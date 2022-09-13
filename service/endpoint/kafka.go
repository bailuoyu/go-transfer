/*
 * Copyright 2020-2021 the original author(https://github.com/wj596)
 *
 * <p>
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 * </p>
 */
package endpoint

import (
	"crypto/sha256"
	"crypto/sha512"
	"github.com/go-mysql-org/go-mysql/canal"
	"github.com/xdg/scram"
	"hash"
	"log"
	"strings"
	"sync"

	"github.com/Shopify/sarama"
	"github.com/go-mysql-org/go-mysql/mysql"
	"github.com/juju/errors"

	"github.com/bailuoyu/go-transfer/global"
	"github.com/bailuoyu/go-transfer/metrics"
	"github.com/bailuoyu/go-transfer/model"
	"github.com/bailuoyu/go-transfer/service/luaengine"
	"github.com/bailuoyu/go-transfer/util/logs"
)

var SHA256 scram.HashGeneratorFcn = func() hash.Hash { return sha256.New() }
var SHA512 scram.HashGeneratorFcn = func() hash.Hash { return sha512.New() }

type XDGSCRAMClient struct {
	*scram.Client
	*scram.ClientConversation
	scram.HashGeneratorFcn
}

func (x *XDGSCRAMClient) Begin(userName, password, authzID string) (err error) {
	x.Client, err = x.HashGeneratorFcn.NewClient(userName, password, authzID)
	if err != nil {
		return err
	}
	x.ClientConversation = x.Client.NewConversation()
	return nil
}

func (x *XDGSCRAMClient) Step(challenge string) (response string, err error) {
	response, err = x.ClientConversation.Step(challenge)
	return
}

func (x *XDGSCRAMClient) Done() bool {
	return x.ClientConversation.Done()
}

type KafkaEndpoint struct {
	client   sarama.Client
	producer sarama.AsyncProducer

	retryLock sync.Mutex
}

func newKafkaEndpoint() *KafkaEndpoint {
	r := &KafkaEndpoint{}
	return r
}

func (s *KafkaEndpoint) Connect() error {
	cfg := sarama.NewConfig()
	cfg.Producer.Partitioner = sarama.NewRandomPartitioner

	if global.Cfg().KafkaSASLUser != "" && global.Cfg().KafkaSASLPassword != "" {
		cfg.Net.SASL.Enable = true
		if global.Cfg().KafkaSASLType == sarama.SASLTypeSCRAMSHA256 {
			cfg.Net.SASL.Mechanism = sarama.SASLTypeSCRAMSHA256
			cfg.Net.SASL.SCRAMClientGeneratorFunc = func() sarama.SCRAMClient { return &XDGSCRAMClient{HashGeneratorFcn: SHA256} }
		} else {
			cfg.Net.SASL.Mechanism = sarama.SASLTypeSCRAMSHA512
			cfg.Net.SASL.SCRAMClientGeneratorFunc = func() sarama.SCRAMClient { return &XDGSCRAMClient{HashGeneratorFcn: SHA512} }
		}
		cfg.Net.SASL.Handshake = true
		cfg.Net.SASL.User = global.Cfg().KafkaSASLUser
		cfg.Net.SASL.Password = global.Cfg().KafkaSASLPassword
	}

	if global.Cfg().KafkaVersion == 1 {
		cfg.Version = sarama.V1_0_0_0
	} else {
		cfg.Version = sarama.V2_0_0_0
	}

	var err error
	var client sarama.Client
	ls := strings.Split(global.Cfg().KafkaAddr, ",")
	client, err = sarama.NewClient(ls, cfg)
	if err != nil {
		return errors.Errorf("unable to create kafka client: %q", err)
	}

	var producer sarama.AsyncProducer
	producer, err = sarama.NewAsyncProducerFromClient(client)
	if err != nil {
		return errors.Errorf("unable to create kafka producer: %q", err)
	}

	s.producer = producer
	s.client = client

	return nil
}

func (s *KafkaEndpoint) Ping() error {
	return s.client.RefreshMetadata()
}

func (s *KafkaEndpoint) Consume(from mysql.Position, rows []*model.RowRequest) error {
	var ms []*sarama.ProducerMessage
	for _, row := range rows {
		rule, _ := global.RuleIns(row.RuleKey)
		if rule.TableColumnSize != len(row.Row) {
			logs.Warnf("%s schema mismatching", row.RuleKey)
			continue
		}

		metrics.UpdateActionNum(row.Action, row.RuleKey)

		if rule.LuaEnable() {
			ls, err := s.buildMessages(row, rule)
			if err != nil {
				log.Println("Lua 脚本执行失败!!! ,详情请参见日志")
				return errors.Errorf("lua 脚本执行失败 : %s ", errors.ErrorStack(err))
			}
			ms = append(ms, ls...)
		} else {
			m, err := s.buildMessage(row, rule)
			if err != nil {
				return errors.Errorf(errors.ErrorStack(err))
			}
			ms = append(ms, m)
		}
	}

	for _, m := range ms {
		s.producer.Input() <- m
		select {
		case err := <-s.producer.Errors():
			return err
		default:
		}
	}

	//logs.Infof("处理完成 %d 条数据", len(rows))
	logs.PrintInfof("处理完成 %d 条数据", len(rows))
	return nil
}

func (s *KafkaEndpoint) Stock(rows []*model.RowRequest) int64 {
	expect := true
	for _, row := range rows {
		rule, _ := global.RuleIns(row.RuleKey)
		if rule.TableColumnSize != len(row.Row) {
			logs.Warnf("%s schema mismatching", row.RuleKey)
			continue
		}

		if rule.LuaEnable() {
			ls, err := s.buildMessages(row, rule)
			if err != nil {
				logs.Errorf(errors.ErrorStack(err))
				expect = false
				break
			}
			for _, m := range ls {
				s.producer.Input() <- m
				select {
				case err := <-s.producer.Errors():
					logs.Error(err.Error())
					expect = false
					break
				default:
				}
			}
			if !expect {
				break
			}
		} else {
			m, err := s.buildMessage(row, rule)
			if err != nil {
				logs.Errorf(errors.ErrorStack(err))
				expect = false
				break
			}
			s.producer.Input() <- m
			select {
			case err := <-s.producer.Errors():
				logs.Error(err.Error())
				expect = false
				break
			default:

			}
		}
	}

	if !expect {
		return 0
	}

	return int64(len(rows))
}

func (s *KafkaEndpoint) buildMessages(row *model.RowRequest, rule *global.Rule) ([]*sarama.ProducerMessage, error) {
	var err error
	var ls []*model.MQRespond
	kvm := rowMap(row, rule, true)
	if row.Action == canal.UpdateAction {
		previous := oldRowMap(row, rule, true)
		ls, err = luaengine.DoMQOps(kvm, previous, row, rule)
	} else {
		ls, err = luaengine.DoMQOps(kvm, nil, row, rule)
	}
	if err != nil {
		return nil, errors.Errorf("lua 脚本执行失败 : %s ", err)
	}

	var ms []*sarama.ProducerMessage
	for _, resp := range ls {
		m := &sarama.ProducerMessage{
			Topic: resp.Topic,
			Value: sarama.ByteEncoder(resp.ByteArray),
		}
		logs.Infof("topic: %s, message: %s", resp.Topic, string(resp.ByteArray))
		ms = append(ms, m)
	}

	return ms, nil
}

func (s *KafkaEndpoint) buildMessage(row *model.RowRequest, rule *global.Rule) (*sarama.ProducerMessage, error) {
	kvm := rowMap(row, rule, false)
	resp := new(model.MQRespond)
	resp.Schema = rule.Schema
	resp.Table = rule.Table
	resp.LogName = row.LogName
	resp.LogPos = row.LogPos
	resp.Action = row.Action
	resp.Timestamp = row.Timestamp
	if rule.ValueEncoder == global.ValEncoderJson {
		resp.Data = kvm
	} else {
		resp.Data = encodeValue(rule, kvm)
	}

	if rule.ReserveRawData && canal.UpdateAction == row.Action {
		resp.Old = oldRowMap(row, rule, false)
	}

	body, err := json.Marshal(resp)
	if err != nil {
		return nil, err
	}
	m := &sarama.ProducerMessage{
		Topic: rule.KafkaTopic,
		Value: sarama.ByteEncoder(body),
	}
	//logs.Infof("topic: %s, message: %s", rule.KafkaTopic, string(body))
	logs.PrintInfof("topic: %s, message: %s", rule.KafkaTopic, body)
	return m, nil
}

func (s *KafkaEndpoint) Close() {
	if s.producer != nil {
		s.producer.Close()
	}
	if s.client != nil {
		s.client.Close()
	}
}
