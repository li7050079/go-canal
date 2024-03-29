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
package luaengine

import (
	"github.com/juju/errors"
	"github.com/siddontang/go-mysql/canal"
	"github.com/yuin/gopher-lua"

	"go-canel/global"
	"go-canel/model"
	"go-canel/util/stringutil"
)

func rdbmsModule(L *lua.LState) int {
	t := L.NewTable()
	L.SetFuncs(t, _rdbmsModuleApi)
	L.Push(t)
	return 1
}

var _rdbmsModuleApi = map[string]lua.LGFunction{
	"rawRow":    rawRow,
	"rawAction": rawAction,

	"INSERT": rdbmsInsert,
	"UPDATE": rdbmsUpdate,
	"DELETE": rdbmsDelete,
	"UPSERT": rdbmsUpsert,
}

func rdbmsInsert(L *lua.LState) int {
	tableName := L.CheckAny(1)
	table := L.CheckAny(2)

	data := L.NewTable()
	L.SetTable(data, lua.LString("tableName"), tableName)
	L.SetTable(data, lua.LString("action"), lua.LString(canal.InsertAction))
	L.SetTable(data, lua.LString("table"), table)

	ret := L.GetGlobal(_globalRET)
	L.SetTable(ret, lua.LString(stringutil.UUID()), data)
	return 0
}

func rdbmsUpdate(L *lua.LState) int {
	tableName := L.CheckAny(1)
	id := L.CheckAny(2)
	table := L.CheckAny(3)

	data := L.NewTable()
	L.SetTable(data, lua.LString("tableName"), tableName)
	L.SetTable(data, lua.LString("action"), lua.LString(canal.UpdateAction))
	L.SetTable(data, lua.LString("id"), id)
	L.SetTable(data, lua.LString("table"), table)

	ret := L.GetGlobal(_globalRET)
	L.SetTable(ret, lua.LString(stringutil.UUID()), data)
	return 0
}

func rdbmsUpsert(L *lua.LState) int {
	tableName := L.CheckAny(1)
	id := L.CheckAny(2)
	table := L.CheckAny(3)

	data := L.NewTable()
	L.SetTable(data, lua.LString("tableName"), tableName)
	L.SetTable(data, lua.LString("action"), lua.LString(global.UpsertAction))
	L.SetTable(data, lua.LString("id"), id)
	L.SetTable(data, lua.LString("table"), table)

	ret := L.GetGlobal(_globalRET)
	L.SetTable(ret, lua.LString(stringutil.UUID()), data)
	return 0
}

func rdbmsDelete(L *lua.LState) int {
	tableName := L.CheckAny(1)
	id := L.CheckAny(2)

	data := L.NewTable()
	L.SetTable(data, lua.LString("tableName"), tableName)
	L.SetTable(data, lua.LString("action"), lua.LString(canal.DeleteAction))
	L.SetTable(data, lua.LString("id"), id)

	ret := L.GetGlobal(_globalRET)
	L.SetTable(ret, lua.LString(stringutil.UUID()), data)
	return 0
}

func DoRdbmsOps(input map[string]interface{}, action string, rule *global.Rule) ([]*model.RdbmsRespond, error) {
	L := _pool.Get()
	defer _pool.Put(L)

	row := L.NewTable()
	paddingTable(L, row, input)
	ret := L.NewTable()
	L.SetGlobal(_globalRET, ret)
	L.SetGlobal(_globalROW, row)
	L.SetGlobal(_globalACT, lua.LString(action))

	funcFromProto := L.NewFunctionFromProto(rule.LuaProto)
	L.Push(funcFromProto)
	err := L.PCall(0, lua.MultRet, nil)
	if err != nil {
		return nil, err
	}

	asserted := true
	responds := make([]*model.RdbmsRespond, 0, ret.Len())
	ret.ForEach(func(k lua.LValue, v lua.LValue) {
		resp := new(model.RdbmsRespond)
		resp.TableName = lvToString(L.GetTable(v, lua.LString("tableName")))
		resp.Action = lvToString(L.GetTable(v, lua.LString("action")))
		resp.Id = lvToInterface(L.GetTable(v, lua.LString("id")), true)
		lvTable := L.GetTable(v, lua.LString("table"))

		var table map[string]interface{}
		if action != canal.DeleteAction {
			table, asserted = lvToMap(lvTable)
			if !asserted {
				return
			}
			resp.Table = table
		}

		responds = append(responds, resp)
	})

	if !asserted {
		return nil, errors.New("The parameter must be of table type")
	}

	return responds, nil
}
