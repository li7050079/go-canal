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
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/juju/errors"
	_ "github.com/kshvakov/clickhouse"
	"github.com/siddontang/go-log/log"
	"github.com/siddontang/go-mysql/canal"
	"github.com/siddontang/go-mysql/mysql"
	"go-canel/global"
	"go-canel/metrics"
	"go-canel/model"
	"go-canel/service/endpoint/rdbms/constants"
	"go-canel/service/endpoint/rdbms/helpers"
	"go-canel/service/endpoint/rdbms/rdbmsopt"
	"go-canel/service/luaengine"
	"go-canel/util/logs"
	"reflect"
)

const (
	Type      = "clickhouse"
	Insert    = `INSERT INTO %s.%s(%s) VALUES(%s);`
	Update    = `ALTER TABLE %s.%s UPDATE %s WHERE %s=?;`
	Delete    = `ALTER TABLE %s.%s DELETE WHERE %s=?;`
	DeleteAll = `ALTER TABLE %s.%s DELETE WHERE 1;`
)

type ClickhouseEndpoint struct {
	conn *sqlx.DB
}

const DSN = "tcp://%s?username=%s&password=%s&database=%s&read_timeout=10&write_timeout=20"

func newClickhouseEndpoint() *ClickhouseEndpoint {
	r := &ClickhouseEndpoint{}
	myconn, _ := sqlx.Open("clickhouse", buildDSN())
	r.conn = myconn
	return r
}
func buildDSN() string {
	cred := global.Cfg()
	return fmt.Sprintf(DSN, cred.ClickhouseAddr, cred.ClickhouseUsername, cred.ClickhousePassword, cred.ClickhouseDatabase)
}

func (s *ClickhouseEndpoint) Connect() error {
	s.Ping()
	return nil
}

func (s *ClickhouseEndpoint) Ping() error {
	return s.conn.Ping()
}

func (s *ClickhouseEndpoint) Consume(from mysql.Position, rows []*model.RowRequest) error {
	//models := make(map[cKey][]mongo.WriteModel, 0)
	for _, row := range rows {
		rule, _ := global.RuleIns(row.RuleKey)
		if rule.TableColumnSize != len(row.Row) {
			logs.Warnf("%s schema mismatching", row.RuleKey)
			continue
		}

		metrics.UpdateActionNum(row.Action, row.RuleKey)

		if rule.LuaEnable() {
			kvm := rowMap(row, rule, true)
			ls, err := luaengine.DoRdbmsOps(kvm, row.Action, rule)
			if err != nil {
				return errors.Errorf("lua 脚本执行失败 : %s ", errors.ErrorStack(err))
			}
			for _, resp := range ls {
				rdbmsOpt := rdbmsopt.NewRdbmsOpt()
				var query helpers.Query
				switch resp.Action {
				case canal.InsertAction:
					query = rdbmsOpt.GetInsert(resp)
				case canal.UpdateAction:
					query = rdbmsOpt.GetUpdate(resp)
				case global.UpsertAction:
					query = rdbmsOpt.GetUpdate(resp)
				case canal.DeleteAction:
					query = rdbmsOpt.GetDelete(resp)
				}
				s.Exec(query)
				logs.Infof("action:%s, collection:%s, id:%v, data:%v", resp.Action, resp.TableName, resp.Id, resp.Table)

			}
		} else {
			kvm := rowMap(row, rule, false)
			id := primaryKey(row, rule)
			kvm["_id"] = id
			rdbmsOpt := rdbmsopt.NewRdbmsOpt()
			var query helpers.Query
			resp := new(model.RdbmsRespond)
			resp.Schema = rule.Schema
			resp.TableName = rule.Table
			resp.Id = id
			resp.Action = row.Action
			resp.Table = kvm
			index := rule.TableInfo.PKColumns[0]
			resp.RuleKey = rule.TableInfo.Columns[index].Name
			switch row.Action {
			case canal.InsertAction:
				query = rdbmsOpt.GetInsert(resp)
			case canal.UpdateAction:
				query = rdbmsOpt.GetUpdate(resp)
			case canal.DeleteAction:
				query = rdbmsOpt.GetDelete(resp)
			}
			s.Exec(query)
			logs.Infof("action:%s, collection:%s, id:%v, data:%v", row.Action, rule.Table, id, kvm)

		}
	}

	logs.Infof("处理完成 %d 条数据", len(rows))
	return nil
}

func (s *ClickhouseEndpoint) Stock(rows []*model.RowRequest) int64 {
	//expect := true
	//models := make(map[cKey][]mongo.WriteModel, 0)
	//for _, row := range rows {
	//	rule, _ := global.RuleIns(row.RuleKey)
	//	if rule.TableColumnSize != len(row.Row) {
	//		logs.Warnf("%s schema mismatching", row.RuleKey)
	//		continue
	//	}
	//
	//	if rule.LuaEnable() {
	//		kvm := rowMap(row, rule, true)
	//		ls, err := luaengine.DoMongoOps(kvm, row.Action, rule)
	//		if err != nil {
	//			log.Println("Lua 脚本执行失败!!! ,详情请参见日志")
	//			logs.Errorf("lua 脚本执行失败 : %s ", errors.ErrorStack(err))
	//			expect = false
	//			break
	//		}
	//
	//		for _, resp := range ls {
	//			ccKey := s.collectionKey(rule.MongodbDatabase, resp.Collection)
	//			model := mongo.NewInsertOneModel().SetDocument(resp.Table)
	//			array, ok := models[ccKey]
	//			if !ok {
	//				array = make([]mongo.WriteModel, 0)
	//			}
	//			array = append(array, model)
	//			models[ccKey] = array
	//		}
	//	} else {
	//		kvm := rowMap(row, rule, false)
	//		id := primaryKey(row, rule)
	//		kvm["_id"] = id
	//
	//		ccKey := s.collectionKey(rule.MongodbDatabase, rule.MongodbCollection)
	//		model := mongo.NewInsertOneModel().SetDocument(kvm)
	//		array, ok := models[ccKey]
	//		if !ok {
	//			array = make([]mongo.WriteModel, 0)
	//		}
	//		array = append(array, model)
	//		models[ccKey] = array
	//	}
	//}
	//
	//if !expect {
	//	return 0
	//}
	//
	//var slowly bool
	var sum int64
	//for key, vs := range models {
	//	collection := s.collection(key)
	//	rr, err := collection.BulkWrite(context.Background(), vs)
	//	if err != nil {
	//		if s.isDuplicateKeyError(err.Error()) {
	//			slowly = true
	//		}
	//		logs.Error(errors.ErrorStack(err))
	//		break
	//	}
	//	sum += rr.InsertedCount
	//}
	//
	//if slowly {
	//	logs.Info("do consume slowly ... ... ")
	//	slowlySum, err := s.doConsumeSlowly(rows)
	//	if err != nil {
	//		logs.Warnf(err.Error())
	//	}
	//	return slowlySum
	//}

	return sum
}

func (s *ClickhouseEndpoint) Exec(params helpers.Query) bool {
	if params.Query == "" {
		return true
	}
	tx, _ := s.conn.Begin()
	_, err := tx.Exec(fmt.Sprintf("%v", params.Query), MakeSlice(params.Params)...)

	if err != nil {
		log.Warnf(constants.ErrorExecQuery, "clickhouse", err)
		return false
	}

	defer func() {
		err = tx.Commit()
	}()

	return true
}

func MakeSlice(input interface{}) []interface{} {
	s := reflect.ValueOf(input)
	if s.Kind() != reflect.Slice {
		log.Warnf("sss")
	}

	ret := make([]interface{}, s.Len())

	for i := 0; i < s.Len(); i++ {
		ret[i] = s.Index(i).Interface()
	}

	return ret
}

func (s *ClickhouseEndpoint) Close() {
	if s.conn != nil {
		s.conn.Close()
	}
}
