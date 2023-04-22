[![License](https://img.shields.io/badge/license-Apache%202-4EB1BA.svg)](https://www.apache.org/licenses/LICENSE-2.0.html)

# 简介
go-canal基础架构来源于wj596大神的go-mysql-transfer的v1.0.4 release，在此基础上做了一些本地需要的封装，完善了一些异常处理，增加了缓存处理、mysql、clickhouse等数据库的支持。

以下介绍来源于原版介绍：

go-canal是一款MySQL数据库实时增量同步工具。

能够监听MySQL二进制日志(Binlog)的变动，将变更内容形成指定格式的消息，实时发送到接收端。从而在数据库和接收端之间形成一个高性能、低延迟的增量数据同步更新管道。

# 特性  

1、简单，不依赖其它组件，一键部署

2、集成多种接收端，如：Mysql、Clickhouse、Redis、MongoDB、Elasticsearch、RocketMQ、Kafka、RabbitMQ、HTTP API等，无需编写客户端，开箱即用

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
             <th width="40%">go-canal</th>
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

直接下载安装包:  [点击下载](https://gitee.com/morningcloud/go-canal/releases/)

**源码编译**

1、依赖Golang 1.14 及以上版本

2、设置' GO111MODULE=on '

3、拉取源码 ' git clone https://gitee.com/morningcloud/go-canal.git '

4、进入目录，执行 ' go build '编译

# 全量数据初始化

go-canal -stock

# 运行

**开启MySQL的binlog**

```
#Linux在my.cnf文件
#Windows在my.ini文件
log-bin=mysql-bin # 开启 binlog
binlog-format=ROW # 选择 ROW 模式
server_id=1 # 配置 MySQL replaction 需要定义，不要和 go-canal 的 slave_id 重复
```

**命令行运行**

1、修改app.yml

2、Windows直接运行 go-canal.exe

3、Linux执行 nohup go-canal &

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

* [go-mysql-transfer](github.com/wj596/go-mysql-transfer)

* [go-mysql](github.com/siddontang/go-mysql)

* [go-mysql-elasticsearch](https://github.com/siddontang/go-mysql-elasticsearch)

* [go-redis](https://github.com/go-redis/redis)

* [rocketmq-client-go](https://github.com/apache/rocketmq-client-go)




# 更新日志
**更新计划**
* 后续持续更新同步端适配
* 解决运行中的bug

**v1.0.5.beta.2 （2022-04-22）**

* 同步端适配mysql、clickhouse；

**v1.0.5.beta.1 （2021-04-16）**

* 添加指定表初始存入本地bolt缓存库功能，并可根据指定key随时读取；
* 添加binlog运行前自动进行全量同步配置项；
* 添加指定binlog文件以及position位置运行的配置项；
* 修复以通配符形式监听表，新增符合条件的表名不能自动进行监听的问题；
* 修复删除表操作后，服务中断问题；

