package mysqlopt

import (
	"encoding/json"
	"github.com/siddontang/go-mysql/canal"
	lua "github.com/yuin/gopher-lua"
	"go-canel/util/httpclient"
	"go-canel/util/logs"
	"strconv"
	"sync"
)

var (
	_pool *mysqlOptPool
	_ds   *canal.Canal

	_httpClient *httpclient.HttpClient
)

type mysqlOptPool struct {
	lock  sync.Mutex
	saved []*lua.LState
}

func InitMysqlOpt(ds *canal.Canal) {
	_ds = ds
}

func SelectList(key string, sql string) (map[string]interface{}, error) {
	//sql := fmt.Sprintf("select * from %s ",cache_table);
	res := make(map[string]interface{})
	rs, err := _ds.Execute(sql)
	if err != nil {
		return res, err
	}
	defer rs.Close()
	rowNumber := rs.RowNumber()
	if rowNumber > 0 {
		for i := 0; i < rowNumber; i++ {
			_map := make(map[string]interface{})
			for field, index := range rs.FieldNames {
				v, err := rs.GetValue(i, index)

				if err != nil {
					logs.Error(err.Error())
					continue
				}
				_map[field] = interfaceToV(v)
			}
			res[InterfaceToString(_map[key])] = _map
		}
	}
	return res, nil
}

func interfaceToV(v interface{}) interface{} {
	switch v.(type) {
	case float64:
		ft := v.(float64)
		return ft
	case float32:
		ft := v.(float32)
		return ft
	case int:
		ft := v.(int)
		return ft
	case uint:
		ft := v.(uint)
		return ft
	case int8:
		ft := v.(int8)
		return ft
	case uint8:
		ft := v.(uint8)
		return ft
	case int16:
		ft := v.(int16)
		return ft
	case uint16:
		ft := v.(uint16)
		return ft
	case int32:
		ft := v.(int32)
		return ft
	case uint32:
		ft := v.(uint32)
		return ft
	case int64:
		ft := v.(int64)
		return ft
	case uint64:
		ft := v.(uint64)
		return ft
	case string:
		ft := v.(string)
		return ft
	case []byte:
		ft := string(v.([]byte))
		return ft
	case nil:
		return nil
	default:
		jsonValue, _ := json.Marshal(v)
		return string(jsonValue)
	}

}
func InterfaceToString(v interface{}) string {
	switch v.(type) {
	case float64:
		ft := strconv.FormatFloat(v.(float64), 'E', -1, 64)
		return ft
	case float32:
		ft := strconv.FormatFloat(float64(v.(float32)), 'E', -1, 64)
		return ft
	case int:
		ft := v.(int)
		return strconv.FormatInt(int64(ft), 10)
	case uint:
		ft := v.(uint)
		return strconv.FormatInt(int64(ft), 10)
	case int8:
		ft := v.(int8)
		return strconv.FormatInt(int64(ft), 10)
	case uint8:
		ft := v.(uint8)
		return strconv.FormatInt(int64(ft), 10)
	case int16:
		ft := v.(int16)
		return strconv.FormatInt(int64(ft), 10)
	case uint16:
		ft := v.(uint16)
		return strconv.FormatInt(int64(ft), 10)
	case int32:
		ft := v.(int32)
		return strconv.FormatInt(int64(ft), 10)
	case uint32:
		ft := v.(uint32)
		return strconv.FormatInt(int64(ft), 10)
	case int64:
		ft := v.(int64)
		return strconv.FormatInt(ft, 10)
	case uint64:
		ft := v.(uint64)
		return strconv.FormatInt(int64(ft), 10)
	case string:
		ft := v.(string)
		return ft
	case []byte:
		ft := string(v.([]byte))
		return ft
	case nil:
		return ""
	default:
		jsonValue, _ := json.Marshal(v)
		return string(jsonValue)
	}
}
