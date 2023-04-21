package rdbmsopt

import (
	"fmt"
	"go-canel/model"
	"go-canel/service/endpoint/rdbms/helpers"
	"strings"
)

const (
	ClickhouseInsert    = `INSERT INTO %s.%s(%s) VALUES(%s);`
	ClickhouseUpdate    = `ALTER TABLE %s.%s UPDATE %s WHERE %s;`
	ClickhouseDelete    = `ALTER TABLE %s.%s DELETE WHERE %s;`
	ClickhouseDeleteAll = `ALTER TABLE %s.%s DELETE WHERE 1;`
)

type ClickhouseOpt struct {
}

func newClickhouseOpt() *ClickhouseOpt {
	r := &ClickhouseOpt{}
	return r
}

func (model *ClickhouseOpt) GetInsert(resq *model.RdbmsRespond) helpers.Query {
	var params []interface{}
	var fieldNames []string
	var fieldValues []string

	for key, value := range resq.Table {
		fieldNames = append(fieldNames, "`"+key+"`")
		fieldValues = append(fieldValues, "?")

		params = append(params, value)
	}

	query := fmt.Sprintf(ClickhouseInsert, resq.Schema, resq.TableName, strings.Join(fieldNames, ","), strings.Join(fieldValues, ","))

	return helpers.Query{
		Query:  query,
		Params: params,
	}
}

func (model *ClickhouseOpt) GetUpdate(resq *model.RdbmsRespond) helpers.Query {
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

	query := fmt.Sprintf(ClickhouseUpdate, resq.Schema, resq.TableName, strings.Join(fields, ", "), strings.Join(where, " and "))

	return helpers.Query{
		Query:  query,
		Params: params,
	}
}

func (model *ClickhouseOpt) GetDelete(resq *model.RdbmsRespond) helpers.Query {
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
	query = fmt.Sprintf(ClickhouseDelete, resq.Schema, resq.TableName, strings.Join(where, " and "))

	return helpers.Query{
		Query:  query,
		Params: params,
	}
}

func (model *ClickhouseOpt) GetCommitTransaction() helpers.Query {
	return helpers.Query{
		Query:  "",
		Params: []interface{}{},
	}
}

func (model *ClickhouseOpt) GetBeginTransaction() helpers.Query {
	return helpers.Query{
		Query:  "",
		Params: []interface{}{},
	}
}

func (model *ClickhouseOpt) GetRollbackTransaction() helpers.Query {
	return helpers.Query{
		Query:  "",
		Params: []interface{}{},
	}
}
