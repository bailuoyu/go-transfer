[![License](https://img.shields.io/badge/license-Apache%202-4EB1BA.svg)](https://www.apache.org/licenses/LICENSE-2.0.html)

# 简介

go-transfer是一款MySQL数据库实时增量同步工具,即canal，基于开源项目go-mysql-transfer二次开发。

能够监听MySQL二进制日志(Binlog)的变动，将变更内容形成指定格式的消息，实时发送到接收端。从而在数据库和接收端之间形成一个高性能、低延迟的增量数据同步更新管道。

# 主要二次开发内容

1、修复es账号密码无法连接的严重bug

2、修复kafka账号密码无法连接的严重bug

3、解决同步到队列的顺序一致性问题，可以安全用于生产环境及高并发场景

4、添加了用于docker的必要配置参数

5、优化同步参数，使信息更加详细，并修复时间参数格式化bug

6、添加了更加详细的配置demo

# 特性  

1、简单，不依赖其它组件，一键部署

2、集成多种接收端，如：Redis、MongoDB、Elasticsearch、RocketMQ、Kafka、RabbitMQ、HTTP API等，无需编写客户端，开箱即用

3、内置丰富的数据解析、消息生成规则、模板语法

4、支持Lua脚本扩展，可处理复杂逻辑

5、集成Prometheus客户端，支持监控告警

6、集成Web Admin监控页面

7、支持高可用集群部署

8、数据同步失败重试

9、支持全量数据初始化


# 原理

1、将自己伪装为MySQL的Slave监听binlog，获取binlog的变更数据

2、根据规则或者lua脚本解析数据，生成指定格式的消息

3、将生成的消息批量发送给接收端

# 与同类工具比较

<table>
    <thead>
        <tr>
            <th width="20%">特色</th>
            <th width="20%">Canal</th>
            <th width="20%">mysql_stream</th>
             <th width="40%">go-mysql-transfer</th>
        </tr>
    </thead>
    <tbody>
        <tr>
            <td>开发语言</td>
            <td>Java</td>
             <td>Python</td>
             <td>Golang</td>
        </tr>
        <tr>
            <td>高可用</td>
            <td>支持</td>
             <td>支持</td>
             <td>支持</td>
        </tr>
        <tr>
            <td>接收端</td>
            <td>编码定制</td>
             <td>Kafka等(MQ)</td>
             <td>Redis、MongoDB、Elasticsearch、RabbitMQ、Kafka、RocketMQ、HTTP API  <br>后续支持更多</td>
        </tr>
        <tr>
            <td>全量数据初始化</td>
            <td>不支持</td>
             <td>支持</td>
             <td>支持</td>
        </tr>
        <tr>
            <td>数据格式</td>
            <td>编码定制</td>
             <td>Json（固定格式）</td>
             <td>Json（规则配置)<br>模板语法<br>Lua脚本</td>
        </tr>
    </tbody>
</table>

# 安装包

**二进制安装包**

直接下载安装包:  [点击下载](https://github.com/wj596/go-mysql-transfer/releases)

**源码编译**

1、依赖Golang 1.14 及以上版本

2、设置' GO111MODULE=on '

3、拉取源码 ' git clone https://github.com/wj596/go-mysql-transfer.git '

4、进入目录，执行 ' go build '编译

# 全量数据初始化

go-mysql-transfer -stock

# 运行

**开启MySQL的binlog**

```
#Linux在my.cnf文件
#Windows在my.ini文件
log-bin=mysql-bin # 开启 binlog
binlog-format=ROW # 选择 ROW 模式
server_id=1 # 配置 MySQL replaction 需要定义，不要和 go-mysql-transfer 的 slave_id 重复
```

**命令行运行**

1、修改app.yml

2、Windows直接运行 go-mysql-transfer.exe

3、Linux执行 nohup go-mysql-transfer &


# gitee

如果您的github访问不稳定，可以在码云(gitee)上star项目：[go-mysql-transfer 码云(gitee)](https://gitee.com/wj596/go-mysql-transfer) 

# 使用说明

* [高可用集群](https://www.kancloud.cn/wj596/go-mysql-transfer/2116627)
* [同步数据到Redis](https://www.kancloud.cn/wj596/go-mysql-transfer/2064427)
    * [Redis配置](https://www.kancloud.cn/wj596/go-mysql-transfer/2111996)
    * [基于规则同步](https://www.kancloud.cn/wj596/go-mysql-transfer/2111997)
    * [基于Lua脚本同步](https://www.kancloud.cn/wj596/go-mysql-transfer/2111998)
* [同步数据到MongoDB](https://www.kancloud.cn/wj596/go-mysql-transfer/2064428)
    * [MongoDB配置](https://www.kancloud.cn/wj596/go-mysql-transfer/2111999)
    * [基于规则同步](https://www.kancloud.cn/wj596/go-mysql-transfer/2112000)
    * [基于Lua脚本同步](https://www.kancloud.cn/wj596/go-mysql-transfer/2112001)
* [同步数据到RocketMQ](https://www.kancloud.cn/wj596/go-mysql-transfer/2064429)
    * [RocketMQ配置](https://www.kancloud.cn/wj596/go-mysql-transfer/2112002)
    * [基于规则同步](https://www.kancloud.cn/wj596/go-mysql-transfer/2112003)
    * [基于Lua脚本同步](https://www.kancloud.cn/wj596/go-mysql-transfer/2112004)
* [同步数据到Kafka](https://www.kancloud.cn/wj596/go-mysql-transfer/2064430)
    * [Kafka配置](https://www.kancloud.cn/wj596/go-mysql-transfer/2112005)
    * [基于规则同步](https://www.kancloud.cn/wj596/go-mysql-transfer/2112006)
    * [基于Lua脚本同步](https://www.kancloud.cn/wj596/go-mysql-transfer/2112007)
* [同步数据到RabbitMQ](https://www.kancloud.cn/wj596/go-mysql-transfer/2064431)
    * [RabbitMQ配置](https://www.kancloud.cn/wj596/go-mysql-transfer/2112008)
    * [基于规则同步](https://www.kancloud.cn/wj596/go-mysql-transfer/2112009)
    * [基于Lua脚本同步](https://www.kancloud.cn/wj596/go-mysql-transfer/2112010)
* [同步数据到Elasticsearch](https://www.kancloud.cn/wj596/go-mysql-transfer/2064432)
    * [Elasticsearch配置](https://www.kancloud.cn/wj596/go-mysql-transfer/2112011)
    * [基于规则同步](https://www.kancloud.cn/wj596/go-mysql-transfer/2112012)
    * [基于Lua脚本同步](https://www.kancloud.cn/wj596/go-mysql-transfer/2112013)
* [全量数据导入](https://www.kancloud.cn/wj596/go-mysql-transfer/2116628)
* [Lua脚本](https://www.kancloud.cn/wj596/go-mysql-transfer/2064433)
    * [基础模块](https://www.kancloud.cn/wj596/go-mysql-transfer/2112014)
    * [Json模块](https://www.kancloud.cn/wj596/go-mysql-transfer/2112015)
    * [HttpClient模块](https://www.kancloud.cn/wj596/go-mysql-transfer/2112016)
    * [DBClient模块](https://www.kancloud.cn/wj596/go-mysql-transfer/2112017)
* [监控](https://www.kancloud.cn/wj596/go-mysql-transfer/2064434)
* [性能测试](https://www.kancloud.cn/wj596/go-mysql-transfer/2116629)
* [常见问题](https://www.kancloud.cn/wj596/go-mysql-transfer/2064435)


# 感谢

* [go-mysql](github.com/siddontang/go-mysql)

* [go-mysql-elasticsearch](https://github.com/siddontang/go-mysql-elasticsearch)

* [go-redis](https://github.com/go-redis/redis)

* [rocketmq-client-go](https://github.com/apache/rocketmq-client-go)


