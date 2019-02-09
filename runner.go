package dbx

import (
    "database/sql"
    "database/sql/driver"
    "github.com/tietang/dbx/mapping"
)

type Mapper func(model interface{}, rows *sql.Rows) (interface{}, error)
type RowsMapper func(rows *sql.Rows) (interface{}, error)
type RowMapper func(row *sql.Row) (interface{}, error)

type TxRunner struct {
    *Runner
    driver.Tx
}

type Runner struct {
    sqlExecutor
    mapperExecutor
    mapping.EntityMapper
    ILogger
    LoggerSettings
}

func NewTxRunner(tx *sql.Tx) *TxRunner {
    r := &TxRunner{}
    r.Runner = &Runner{}
    r.sqlExecutor = tx
    r.Tx = tx
    return r
}

func NewRunner(se sqlExecutor, em mapping.EntityMapper) *Runner {
    r := &Runner{}
    r.sqlExecutor = se
    r.EntityMapper = em
    return r
}
