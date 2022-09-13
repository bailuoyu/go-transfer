package model

import "sync"

var mqRespondPool = sync.Pool{
	New: func() interface{} {
		return new(MQRespond)
	},
}

var esRespondPool = sync.Pool{
	New: func() interface{} {
		return new(ESRespond)
	},
}

var mongoRespondPool = sync.Pool{
	New: func() interface{} {
		return new(MongoRespond)
	},
}

var redisRespondPool = sync.Pool{
	New: func() interface{} {
		return new(RedisRespond)
	},
}

type MQRespond struct {
	Topic     string      `json:"-"`
	Schema    string      `json:"schema"`
	Table     string      `json:"table"`
	LogName   string      `json:"log_name"`
	LogPos    uint32      `json:"log_pos"`
	Action    string      `json:"action"`
	Timestamp uint32      `json:"timestamp"`
	Data      interface{} `json:"data"`
	Old       interface{} `json:"old,omitempty"`
	ByteArray []byte      `json:"-"`
}

type ESRespond struct {
	Index  string
	Id     string
	Action string
	Date   string
}

type MongoRespond struct {
	RuleKey    string
	Collection string
	Action     string
	Id         interface{}
	Table      map[string]interface{}
}

type RedisRespond struct {
	Action    string
	Structure string
	Key       string
	Field     string
	Score     float64
	OldVal    interface{}
	Val       interface{}
}

func BuildMQRespond() *MQRespond {
	return mqRespondPool.Get().(*MQRespond)
}

func ReleaseMQRespond(t *MQRespond) {
	mqRespondPool.Put(t)
}

func BuildESRespond() *ESRespond {
	return esRespondPool.Get().(*ESRespond)
}

func ReleaseESRespond(t *ESRespond) {
	esRespondPool.Put(t)
}

func BuildMongoRespond() *MongoRespond {
	return mongoRespondPool.Get().(*MongoRespond)
}

func ReleaseMongoRespond(t *MongoRespond) {
	mongoRespondPool.Put(t)
}

func BuildRedisRespond() *RedisRespond {
	return redisRespondPool.Get().(*RedisRespond)
}

func ReleaseRedisRespond(t *RedisRespond) {
	redisRespondPool.Put(t)
}
