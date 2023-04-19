package rdbmsopt

import (
	"go-canel/global"
	"go-canel/model"
	"go-canel/service/endpoint/rdbms/helpers"
)

const (
	Type      = "clickhouse"
	Insert    = `INSERT INTO %s.%s(%s) VALUES(%s);`
	Update    = `ALTER TABLE %s.%s UPDATE %s WHERE %s=?;`
	Delete    = `ALTER TABLE %s.%s DELETE WHERE %s=?;`
	DeleteAll = `ALTER TABLE %s.%s DELETE WHERE 1;`
)

type RdbmsOpt interface {
	GetInsert(*model.RdbmsRespond) helpers.Query
	GetUpdate(*model.RdbmsRespond) helpers.Query
	GetDelete(*model.RdbmsRespond) helpers.Query
}

func NewRdbmsOpt() RdbmsOpt {
	cfg := global.Cfg()
	if cfg.IsMysql() {
		return newMysqlOpt()
	} else if cfg.IsClickhouse() {
		return newClickhouseOpt()
	}
	return nil
}
