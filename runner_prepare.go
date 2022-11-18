package dbx

import (
	"context"
	"database/sql"
	"time"
)

func (r *PrepareTxRunner) Exec(params ...interface{}) (rs sql.Result, err error) {
	return r.ExecContext(context.Background(), params...)
}
func (r *PrepareTxRunner) ExecContext(ctx context.Context, params ...interface{}) (rs sql.Result, err error) {
	if r.LoggingEnabled() {
		defer func(start time.Time) {
			r.Logger().Log(&QueryStatus{
				Query:   r.sql,
				Args:    params,
				Err:     err,
				Start:   start,
				End:     time.Now(),
				Context: ctx,
			})
		}(time.Now())
	}

	return r.sqlPrepareExecutor.ExecContext(ctx, params...)
}

func (r *PrepareTxRunner) Execute(params ...interface{}) (rowsAffected int64, err error) {
	return r.ExecuteContext(context.Background(), params...)
}
func (r *PrepareTxRunner) ExecuteContext(ctx context.Context, params ...interface{}) (rowsAffected int64, err error) {
	rs, err := r.ExecContext(ctx, params...)
	if err != nil {
		return 0, err
	}
	rowsAffected, err = rs.RowsAffected()
	return rowsAffected, err
}

func (r *PrepareTxRunner) Query(params ...interface{}) (*sql.Rows, error) {
	return r.QueryContext(context.Background(), params...)
}

func (r *PrepareTxRunner) QueryContext(ctx context.Context, params ...interface{}) (rows *sql.Rows, err error) {
	if r.LoggingEnabled() {
		defer func(start time.Time) {
			r.Logger().Log(&QueryStatus{
				Query:   r.sql,
				Args:    params,
				Err:     err,
				Start:   start,
				End:     time.Now(),
				Context: ctx,
			})
		}(time.Now())
	}

	return r.sqlPrepareExecutor.QueryContext(ctx, params...)

}

func (r *PrepareTxRunner) QueryRow(params ...interface{}) *sql.Row {
	return r.QueryRowContext(context.Background(), params...)
}

func (r *PrepareTxRunner) QueryRowContext(ctx context.Context, params ...interface{}) *sql.Row {
	var err error
	if r.LoggingEnabled() {
		defer func(start time.Time) {
			r.Logger().Log(&QueryStatus{
				Query:   r.sql,
				Args:    params,
				Err:     err,
				Start:   start,
				End:     time.Now(),
				Context: ctx,
			})
		}(time.Now())
	}
	return r.sqlPrepareExecutor.QueryRowContext(ctx, params...)

}
