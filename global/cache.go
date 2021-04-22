package global

import (
	"github.com/vmihailenco/msgpack"
	"strings"
	"sync"
)

var (
	_cacheInsMap       = make(map[string]*Cache)
	_lockOfCacheInsMap sync.RWMutex
)

type Cache struct {
	//IsDoSql					 bool `yaml:"is_do_sql"`					//是否仅执行sql，如果是,则只传入do_sql
	//doSql					 string `yaml:"do_sql`						//sql内容
	Schema                  string `yaml:"schema"`
	Table                   string `yaml:"table"`
	StorageKey              string `yaml:"storage_key"`                //作为键的字段，一般选择唯一主键
	ColumnLowerCase         bool   `yaml:"column_lower_case"`          // 列名称转为小写
	ColumnUpperCase         bool   `yaml:"column_upper_case"`          // 列名称转为大写
	ColumnUnderscoreToCamel bool   `yaml:"column_underscore_to_camel"` // 列名称下划线转驼峰
	IncludeColumnConfig     string `yaml:"include_columns"`            // 包含的列
	ExcludeColumnConfig     string `yaml:"exclude_columns"`            // 排除掉的列
	ColumnMappingConfigs    string `yaml:"column_mappings"`            // 列名称映射
}

func CacheDeepClone(res *Cache) (*Cache, error) {
	data, err := msgpack.Marshal(res)
	if err != nil {
		return nil, err
	}

	var r Cache
	err = msgpack.Unmarshal(data, &r)
	if err != nil {
		return nil, err
	}

	return &r, nil
}

func CacheKey(schema string, table string) string {
	return strings.ToLower(schema + ":" + table)
}

func CacheStorageKey(cacheKey string, key string) string {
	return strings.ToLower(cacheKey + ":" + key)
}

func AddCacheIns(cacheKey string, r *Cache) {
	_lockOfCacheInsMap.Lock()
	defer _lockOfCacheInsMap.Unlock()

	_cacheInsMap[cacheKey] = r
}

func CacheIns(cacheKey string) (*Cache, bool) {
	_lockOfCacheInsMap.RLock()
	defer _lockOfCacheInsMap.RUnlock()

	r, ok := _cacheInsMap[cacheKey]

	return r, ok
}

func CacheInsExist(cacheKey string) bool {
	_lockOfCacheInsMap.RLock()
	defer _lockOfCacheInsMap.RUnlock()

	_, ok := _cacheInsMap[cacheKey]

	return ok
}

func CacheInsDel(cacheKey string) bool {
	_lockOfCacheInsMap.RLock()
	defer _lockOfCacheInsMap.RUnlock()

	delete(_cacheInsMap, cacheKey)

	return true
}

func CacheInsTotal() int {
	_lockOfCacheInsMap.RLock()
	defer _lockOfCacheInsMap.RUnlock()

	return len(_cacheInsMap)
}

func CacheInsList() []*Cache {
	_lockOfCacheInsMap.RLock()
	defer _lockOfCacheInsMap.RUnlock()

	list := make([]*Cache, 0, len(_cacheInsMap))
	for _, cache := range _cacheInsMap {
		list = append(list, cache)
	}

	return list
}

func CacheKeyList() []string {
	_lockOfCacheInsMap.RLock()
	defer _lockOfCacheInsMap.RUnlock()

	list := make([]string, 0, len(_cacheInsMap))
	for k, _ := range _cacheInsMap {
		list = append(list, k)
	}

	return list
}

func (s *Cache) Initialize() error {

	return nil
}
