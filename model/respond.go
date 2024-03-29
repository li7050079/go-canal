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

var rdbmsRespondPool = sync.Pool{
	New: func() interface{} {
		return new(RdbmsRespond)
	},
}

var redisRespondPool = sync.Pool{
	New: func() interface{} {
		return new(RedisRespond)
	},
}

type MQRespond struct {
	Topic     string      `json:"-"`
	Action    string      `json:"action"`
	Timestamp uint32      `json:"timestamp"`
	Raw       interface{} `json:"raw,omitempty"`
	Date      interface{} `json:"date"`
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

type RdbmsRespond struct {
	RuleKey   string
	Schema    string //库名
	TableName string //同步后的表名
	Action    string
	Id        interface{}            //主键值
	OldId     interface{}            //原主键值
	IdName    string                 //！同步后！的主键名称
	Table     map[string]interface{} //数据
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

func BuildRdbmsRespond() *RdbmsRespond {
	return rdbmsRespondPool.Get().(*RdbmsRespond)
}

func ReleaseRdbmsRespond(t *RdbmsRespond) {
	rdbmsRespondPool.Put(t)
}
