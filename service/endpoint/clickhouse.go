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
	"go-canel/model"
	"go-canel/service/luaengine"
	"go-canel/util/logs"
	"reflect"
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
	//for _, row := range rows {
	//	rule, _ := global.RuleIns(row.RuleKey)
	//	if rule.TableColumnSize != len(row.Row) {
	//		logs.Warnf("%s schema mismatching", row.RuleKey)
	//		continue
	//	}
	//
	//	metrics.UpdateActionNum(row.Action, row.RuleKey)
	//
	//	if rule.LuaEnable() {
	//		kvm := rowMap(row, rule, true)
	//		ls, err := luaengine.DoMongoOps(kvm, row.Action, rule)
	//		if err != nil {
	//			return errors.Errorf("lua 脚本执行失败 : %s ", errors.ErrorStack(err))
	//		}
	//		for _, resp := range ls {
	//			var model mongo.WriteModel
	//
	//			switch resp.Action {
	//			case canal.InsertAction:
	//				model = mongo.NewInsertOneModel().SetDocument(resp.Table)
	//			case canal.UpdateAction:
	//				model = mongo.NewUpdateOneModel().SetFilter(bson.M{"_id": resp.Id}).SetUpdate(bson.M{"$set": resp.Table})
	//			case global.UpsertAction:
	//				model = mongo.NewUpdateOneModel().SetFilter(bson.M{"_id": resp.Id}).SetUpsert(true).SetUpdate(bson.M{"$set": resp.Table})
	//			case canal.DeleteAction:
	//				model = mongo.NewDeleteOneModel().SetFilter(bson.M{"_id": resp.Id})
	//			}
	//
	//			key := s.collectionKey(rule.MongodbDatabase, resp.Collection)
	//			array, ok := models[key]
	//			if !ok {
	//				array = make([]mongo.WriteModel, 0)
	//			}
	//
	//			logs.Infof("action:%s, collection:%s, id:%v, data:%v", resp.Action, resp.Collection, resp.Id, resp.Table)
	//
	//			array = append(array, model)
	//			models[key] = array
	//		}
	//	} else {
	//		kvm := rowMap(row, rule, false)
	//		id := primaryKey(row, rule)
	//		kvm["_id"] = id
	//		var model mongo.WriteModel
	//		switch row.Action {
	//		case canal.InsertAction:
	//			model = mongo.NewInsertOneModel().SetDocument(kvm)
	//		case canal.UpdateAction:
	//			model = mongo.NewUpdateOneModel().SetFilter(bson.M{"_id": id}).SetUpdate(bson.M{"$set": kvm})
	//		case canal.DeleteAction:
	//			model = mongo.NewDeleteOneModel().SetFilter(bson.M{"_id": id})
	//		}
	//
	//		ccKey := s.collectionKey(rule.MongodbDatabase, rule.MongodbCollection)
	//		array, ok := models[ccKey]
	//		if !ok {
	//			array = make([]mongo.WriteModel, 0)
	//		}
	//
	//		logs.Infof("action:%s, collection:%s, id:%v, data:%v", row.Action, rule.MongodbCollection, id, kvm)
	//
	//		array = append(array, model)
	//		models[ccKey] = array
	//	}
	//}
	//
	//var slowly bool
	//for key, model := range models {
	//	collection := s.collection(key)
	//	_, err := collection.BulkWrite(context.Background(), model)
	//	if err != nil {
	//		if s.isDuplicateKeyError(err.Error()) {
	//			slowly = true
	//		} else {
	//			return err
	//		}
	//		logs.Error(errors.ErrorStack(err))
	//		break
	//	}
	//}
	//if slowly {
	//	_, err := s.doConsumeSlowly(rows)
	//	if err != nil {
	//		return err
	//	}
	//}
	//
	//logs.Infof("处理完成 %d 条数据", len(rows))
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

func (s *ClickhouseEndpoint) doConsumeSlowly(rows []*model.RowRequest) (int64, error) {
	var sum int64
	for _, row := range rows {
		rule, _ := global.RuleIns(row.RuleKey)
		if rule.TableColumnSize != len(row.Row) {
			logs.Warnf("%s schema mismatching", row.RuleKey)
			continue
		}

		if rule.LuaEnable() {
			kvm := rowMap(row, rule, true)
			ls, err := luaengine.DoMongoOps(kvm, row.Action, rule)
			if err != nil {
				logs.Errorf("lua 脚本执行失败 : %s ", errors.ErrorStack(err))
				return sum, err
			}
			for _, resp := range ls {
				//collection := s.collection(s.collectionKey(rule.MongodbDatabase, resp.Collection))
				switch resp.Action {
				case canal.InsertAction:
					//_, err := collection.InsertOne(context.Background(), resp.Table)
					if err != nil {
						//if s.isDuplicateKeyError(err.Error()) {
						//	logs.Warnf("duplicate key [ %v ]", stringutil.ToJsonString(resp.Table))
						//} else {
						//	return sum, err
						//}
					}
				case canal.UpdateAction:
					//_, err := collection.UpdateOne(context.Background(), bson.M{"_id": resp.Id}, bson.M{"$set": resp.Table})
					//if err != nil {
					//	return sum, err
					//}
				case canal.DeleteAction:
					//_, err := collection.DeleteOne(context.Background(), bson.M{"_id": resp.Id})
					//if err != nil {
					//	return sum, err
					//}
				}
				//logs.Infof("action:%s, collection:%s, id:%v, data:%v",
				//row.Action, collection.Name(), resp.Id, resp.Table)
			}
		} else {
			kvm := rowMap(row, rule, false)
			id := primaryKey(row, rule)
			kvm["_id"] = id

			//collection := s.collection(s.collectionKey(rule.MongodbDatabase, rule.MongodbCollection))

			switch row.Action {
			case canal.InsertAction:
				//_, err := collection.InsertOne(context.Background(), kvm)
				//if err != nil {
				//	if s.isDuplicateKeyError(err.Error()) {
				//		logs.Warnf("duplicate key [ %v ]", stringutil.ToJsonString(kvm))
				//	} else {
				//		return sum, err
				//	}
				//}
			case canal.UpdateAction:
				//_, err := collection.UpdateOne(context.Background(), bson.M{"_id": id}, bson.M{"$set": kvm})
				//if err != nil {
				//	return sum, err
				//}
			case canal.DeleteAction:
				//_, err := collection.DeleteOne(context.Background(), bson.M{"_id": id})
				//if err != nil {
				//	return sum, err
				//}
			}

			//logs.Infof("action:%s, collection:%s, id:%v, data:%v", row.Action, collection.Name(), id, kvm)
		}
		sum++
	}
	return sum, nil
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
