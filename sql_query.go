package dbx

import (
	"context"
	"database/sql"
)

type sqlExecutor interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	Exec(query string, args ...interface{}) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	Prepare(query string) (*sql.Stmt, error)
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
}

type mapperExecutor interface {
	Execute(sql string, params ...interface{}) (lastInsertId, rowsAffected int64, err error)
	Insert(model interface{}) (rs sql.Result, err error)
	InsertContext(ctx context.Context, model interface{}) (rs sql.Result, err error)
	Update(model interface{}) (rs sql.Result, err error)
	UpdateContext(ctx context.Context, model interface{}) (rs sql.Result, err error)
	Find(resultSlice interface{}, query string, params ...interface{}) (err error)
	FindContext(ctx context.Context, resultSlice interface{}, query string, params ...interface{}) (err error)
	FindExample(querier interface{}, resultSlice interface{}) (err error)
	FindExampleContext(ctx context.Context, querier interface{}, resultSlice interface{}) (err error)
	Select(mapper RowsMapper, resultSlice interface{}, sql string, params ...interface{}) (err error)
	SelectContext(ctx context.Context, mapper RowsMapper, resultSlice interface{}, sql string, params ...interface{}) (err error)
	Get(out interface{}, sql string, params ...interface{}) (ok bool, err error)
	GetContext(ctx context.Context, out interface{}, sql string, params ...interface{}) (ok bool, err error)
	GetOne(out interface{}) (ok bool, err error)
	GetOneContext(ctx context.Context, out interface{}) (ok bool, err error)
}

type Scaner interface {
	Scan(dest ...interface{}) error
}
