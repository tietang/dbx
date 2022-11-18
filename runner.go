package dbx

import (
	"database/sql"
	"database/sql/driver"
	"github.com/tietang/dbx/mapping"
)

type Mapper func(model interface{}, rows *sql.Rows) (interface{}, error)
type RowsMapper func(rows *sql.Rows) (interface{}, error)
type RowMapper func(row *sql.Row) (interface{}, error)

var (
	_ mapperExecutor        = new(Runner)
	_ sqlExecutor           = new(Runner)
	_ mapperPrepareExecutor = new(PrepareTxRunner)
	_ sqlPrepareExecutor    = new(PrepareTxRunner)
	_ mapperExecutor        = new(TxRunner)
	_ sqlExecutor           = new(TxRunner)
	_ mapping.EntityMapper  = new(Runner)
	_ mapping.EntityMapper  = new(TxRunner)
)

type Runner struct {
	sqlExecutor
	mapperExecutor
	mapping.EntityMapper
	ILogger
	LoggerSettings
}

func NewRunner(se sqlExecutor, em mapping.EntityMapper) *Runner {
	r := &Runner{}
	r.sqlExecutor = se
	r.EntityMapper = em
	return r
}

type TxRunner struct {
	*Runner
	driver.Tx
}

func NewTxRunner(tx *sql.Tx) *TxRunner {
	r := &TxRunner{}
	r.Runner = &Runner{}
	r.sqlExecutor = tx
	r.Tx = tx
	return r
}

type PrepareTxRunner struct {
	sqlPrepareExecutor
	mapperPrepareExecutor
	ILogger
	LoggerSettings
	sql string
}

func NewPrepareTxRunner(sql string, stmt *sql.Stmt) *PrepareTxRunner {
	r := &PrepareTxRunner{}
	r.sqlPrepareExecutor = stmt
	r.sql = sql
	return r
}
