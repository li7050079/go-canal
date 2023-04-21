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
	"github.com/juju/errors"
	"github.com/siddontang/go-mysql/canal"
	"github.com/siddontang/go-mysql/client"
	"github.com/siddontang/go-mysql/mysql"
	"github.com/siddontang/go/log"
	"go-canel/global"
	"go-canel/metrics"
	"go-canel/model"
	"go-canel/service/endpoint/rdbms/constants"
	"go-canel/service/endpoint/rdbms/helpers"
	"go-canel/service/endpoint/rdbms/rdbmsopt"
	"go-canel/service/luaengine"
	"go-canel/util/logs"
)

type MysqlEndpoint struct {
	conn *client.Conn
}

func newMysqlEndpoint() *MysqlEndpoint {

	r := &MysqlEndpoint{}
	cfg := global.Cfg()
	myconn, _ := client.Connect(cfg.MysqlAddr, cfg.MysqlUsername, cfg.MysqlPassword, cfg.MysqlDatabase)
	r.conn = myconn
	return r
}

func (s *MysqlEndpoint) Connect() error {
	s.Ping()
	return nil
}

func (s *MysqlEndpoint) Ping() error {
	return s.conn.Ping()
}

func (s *MysqlEndpoint) Consume(from mysql.Position, rows []*model.RowRequest) error {
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
				resp.Schema = rule.Schema
				resp.IdName = primaryKeyName(rule)
				resp.OldId = primaryOldKey(row, rule)
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
			rdbmsOpt := rdbmsopt.NewRdbmsOpt()
			var query helpers.Query
			resp := new(model.RdbmsRespond)
			resp.Schema = rule.Schema
			resp.TableName = rule.Table
			resp.Id = id
			resp.Action = row.Action
			resp.Table = kvm
			resp.IdName = primaryKeyName(rule)
			resp.OldId = primaryOldKey(row, rule)
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

func (s *MysqlEndpoint) Stock(rows []*model.RowRequest) int64 {
	expect := true
	rdbmsOpt := rdbmsopt.NewRdbmsOpt()
	for _, row := range rows {
		rule, _ := global.RuleIns(row.RuleKey)
		if rule.TableColumnSize != len(row.Row) {
			logs.Warnf("%s schema mismatching", row.RuleKey)
			continue
		}

		if rule.LuaEnable() {
			kvm := rowMap(row, rule, true)
			ls, err := luaengine.DoRdbmsOps(kvm, row.Action, rule)
			if err != nil {
				log.Errorf("Lua 脚本执行失败!!! ,详情请参见日志")
				logs.Errorf("lua 脚本执行失败 : %s ", errors.ErrorStack(err))
				expect = false
				break
			}

			for _, resp := range ls {
				resp.Schema = rule.Schema
				query := rdbmsOpt.GetInsert(resp)
				s.Exec(query)
			}
		} else {
			kvm := rowMap(row, rule, false)
			id := primaryKey(row, rule)
			resp := new(model.RdbmsRespond)
			resp.Schema = rule.Schema
			resp.TableName = rule.Table
			resp.Id = id
			resp.IdName = primaryKeyName(rule)
			resp.Action = row.Action
			resp.Table = kvm
			query := rdbmsOpt.GetInsert(resp)
			s.Exec(query)

		}
	}

	if !expect {
		return 0
	}
	return int64(len(rows))
}

func (s *MysqlEndpoint) Exec(params helpers.Query) bool {
	if params.Query == "" {
		return true
	}
	s.conn.Begin()
	_, err := s.conn.Execute(fmt.Sprintf("%v", params.Query), MakeSlice(params.Params)...)

	if err != nil {
		log.Warnf(constants.ErrorExecQuery, "mysql", err)
		return false
	}

	defer func() {
		err = s.conn.Commit()
	}()

	return true
}

func (s *MysqlEndpoint) Close() {
	if s.conn != nil {
		s.conn.Close()
	}
}
