package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"go-mysql-transfer/global"
	"go-mysql-transfer/service/mysqlopt"
	"go.etcd.io/bbolt"
)

type bolkCacheStorage struct {
	Name  string
	Value string
}

var (
	_cacheBucket = []byte("cache")
)

func InitCache() error {
	if len(global.Cfg().CacheTable) == 0 {
		return nil
	}
	if _bolt == nil {
		if err := initBolt(); err != nil {
			return err
		}
	}
	for _, rs := range global.Cfg().CacheTable {
		if rs.Schema != "" && rs.Table != "" && rs.StorageKey != "" {
			cacheKey := global.CacheKey(rs.Schema, rs.Table)
			global.AddCacheIns(cacheKey, rs)
			err := RefreshCacheStorage(cacheKey)
			if err != nil {
				errors.New("库名或表名不能为空")
				continue
			}
		} else {
			return errors.New("库名或表名不能为空")
		}
	}
	return nil
}

func RefreshCacheStorage(cacheKey string) error {
	rs, ok := global.CacheIns(cacheKey)
	if ok {
		sql := fmt.Sprintf("select * from %s", rs.Schema+"."+rs.Table)
		res, _ := mysqlopt.SelectList(rs.StorageKey, sql)
		_bolt.Update(func(tx *bbolt.Tx) error {
			tx.DeleteBucket([]byte(cacheKey))
			tx.CreateBucketIfNotExists([]byte(cacheKey))
			bt := tx.Bucket([]byte(cacheKey))
			for k := range res {
				bboltKey := mysqlopt.InterfaceToString(k)
				_json, _ := json.Marshal(res[k])
				bt.Put([]byte(bboltKey), _json)
			}
			return nil
		})

	}
	return nil
}

func Get(schema string, table string, key string) map[string]interface{} {
	cacheKey := global.CacheKey(schema, table)
	var result = make(map[string]interface{})
	err := _bolt.View(func(tx *bbolt.Tx) error {
		bx := tx.Bucket([]byte(cacheKey))
		res := bx.Get([]byte(key))
		err := json.Unmarshal(res, &result)
		return err
	})
	if err != nil {

	}
	return result
}
