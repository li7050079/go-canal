package rdbmsopt

import (
	"fmt"
	"go-canel/model"
	"go-canel/service/endpoint/rdbms/helpers"
	"strings"
)

const (
	MysqlInsert    = `INSERT INTO %s.%s(%s) VALUES(%s);`
	MysqlUpdate    = `UPDATE %s.%s SET %s WHERE %s;`
	MysqlDelete    = `DELETE FROM %s.%s WHERE %s;`
	MysqlDeleteAll = `DELETE FROM %s.%s WHERE 1=1;`
)

type MysqlOpt struct {
}

func newMysqlOpt() *MysqlOpt {
	r := &MysqlOpt{}
	return r
}

func (model *MysqlOpt) GetInsert(resq *model.RdbmsRespond) helpers.Query {
	var params []interface{}
	var fieldNames []string
	var fieldValues []string

	for key, value := range resq.Table {
		fieldNames = append(fieldNames, "`"+key+"`")
		fieldValues = append(fieldValues, "?")

		params = append(params, value)
	}

	query := fmt.Sprintf(MysqlInsert, resq.Schema, resq.TableName, strings.Join(fieldNames, ","), strings.Join(fieldValues, ","))

	return helpers.Query{
		Query:  query,
		Params: params,
	}
}

func (model *MysqlOpt) GetUpdate(resq *model.RdbmsRespond) helpers.Query {
	var params []interface{}
	var fields []string
	var where []string
	for key, value := range resq.Table {
		fields = append(fields, "`"+key+"`"+"=?")
		params = append(params, value)
	}
	// add key to params
	if oldId, ok := resq.OldId.(string); ok == true {
		keys := strings.Split(resq.IdName, "@@")
		values := strings.Split(oldId, "@@")
		for index := range keys {
			where = append(where, "`"+keys[index]+"`=?")
			params = append(params, values[index])
		}
	} else {
		where = append(where, "`"+resq.IdName+"`=?")
		params = append(params, resq.OldId)
	}

	query := fmt.Sprintf(MysqlUpdate, resq.Schema, resq.TableName, strings.Join(fields, ", "), strings.Join(where, " and "))

	return helpers.Query{
		Query:  query,
		Params: params,
	}
}

func (model *MysqlOpt) GetDelete(resq *model.RdbmsRespond) helpers.Query {
	var params []interface{}
	var query string
	// add key to params
	var where []string
	if id, ok := resq.Id.(string); ok == true {
		keys := strings.Split(resq.IdName, "@@")
		values := strings.Split(id, "@@")
		for index := range keys {
			where = append(where, "`"+keys[index]+"`=?")
			params = append(params, values[index])
		}
	} else {
		where = append(where, "`"+resq.IdName+"`=?")
		params = append(params, resq.Id)
	}
	query = fmt.Sprintf(MysqlDelete, resq.Schema, resq.TableName, strings.Join(where, " and "))

	return helpers.Query{
		Query:  query,
		Params: params,
	}
}

func (model *MysqlOpt) GetCommitTransaction() helpers.Query {
	return helpers.Query{
		Query:  "",
		Params: []interface{}{},
	}
}

func (model *MysqlOpt) GetBeginTransaction() helpers.Query {
	return helpers.Query{
		Query:  "",
		Params: []interface{}{},
	}
}

func (model *MysqlOpt) GetRollbackTransaction() helpers.Query {
	return helpers.Query{
		Query:  "",
		Params: []interface{}{},
	}
}
