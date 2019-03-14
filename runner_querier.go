package dbx

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"
)

//type Mapper interface {
//    Scans(rows *sql.Rows) (obj interface{}, rowMappers []interface{})
//}

//type Mapper func(model interface{}, rows *sql.Rows) (interface{}, error)
//type RowsMapper func(rows *sql.Rows) (interface{}, error)
//type RowMapper func(row *sql.Row) (interface{}, error)

func (r *Runner) Find(resultSlicePtr interface{}, query string, params ...interface{}) (err error) {
	return r.Select(func(rows *sql.Rows) (i interface{}, e error) {
		return r.rowsMapping(resultSlicePtr, rows)
	}, resultSlicePtr, query, params...)
}
func (r *Runner) FindContext(ctx context.Context, resultSlicePtr interface{}, query string, params ...interface{}) (err error) {
	return r.SelectContext(ctx, func(rows *sql.Rows) (i interface{}, e error) {
		return r.rowsMapping(resultSlicePtr, rows)
	}, resultSlicePtr, query, params...)
}

func (r *Runner) FindTest(resultSlicePtr interface{}, query string, params ...interface{}) (err error) {
	return r.Select(func(rows *sql.Rows) (i interface{}, e error) {
		return r.rowsMapping(resultSlicePtr, rows)
	}, resultSlicePtr, query, params...)
}
func (r *Runner) FindExample(querier interface{}, resultSlicePtr interface{}) (err error) {
	return r.FindExampleContext(context.Background(), querier, resultSlicePtr)
}

func (r *Runner) FindExampleContext(ctx context.Context, querier interface{}, resultSlicePtr interface{}) (err error) {

	entity, ind := r.GetEntity(querier)
	names := ""
	params := make([]interface{}, 0)
	wheres := make([]string, 0)
	for _, fd := range entity.Columns {
		if fd.Field.Anonymous || fd.Embedded {
			continue
		}
		lastv := ind.FieldByIndex(fd.Index)
		names += "`" + fd.ColumnName + "`,"
		if lastv.Kind() == reflect.Ptr && lastv.IsNil() {
			continue
		}
		if !lastv.IsValid() {
			continue
		}
		if fd.ColumnName == "" || lastv.Interface() == fd.Zero.Interface() {
			continue
		}

		wheres = append(wheres, "`"+fd.ColumnName+"`=?")
		params = append(params, lastv.Interface())

	}
	if len(wheres) == 0 {
		err = errors.New("无可用查询条件")
		return
	}

	where := strings.Join(wheres, " and ")
	query := fmt.Sprintf("select %s from `%s` where %s", names[0:len(names)-1], entity.TableName, where)

	if r.LoggingEnabled() {
		defer func(start time.Time) {
			r.Logger().Log(&QueryStatus{
				Query:   query,
				Args:    params,
				Err:     err,
				Start:   start,
				End:     time.Now(),
				Context: ctx,
			})
		}(time.Now())
	}

	return r.SelectContext(ctx, func(rows *sql.Rows) (i interface{}, e error) {
		return r.rowsMapping(resultSlicePtr, rows)
	}, resultSlicePtr, query, params...)
}

func (r *Runner) Select(mapper RowsMapper, resultSlicePtr interface{}, sql string, params ...interface{}) (err error) {
	return r.SelectContext(context.Background(), mapper, resultSlicePtr, sql, params...)
}

func (r *Runner) SelectContext(ctx context.Context, mapper RowsMapper, resultSlicePtr interface{}, sql string, params ...interface{}) (err error) {
	if r.LoggingEnabled() {
		defer func(start time.Time) {
			r.Logger().Log(&QueryStatus{
				Query:   sql,
				Args:    params,
				Err:     err,
				Start:   start,
				End:     time.Now(),
				Context: ctx,
			})
		}(time.Now())
	}

	dstv := reflect.ValueOf(resultSlicePtr)
	if dstv.Kind() != reflect.Ptr {
		return errors.New("needs a pointer to a slice ")
	}

	rows, err := r.sqlExecutor.QueryContext(ctx, sql, params...)
	if err != nil {
		return err
	}

	slicev := dstv.Elem()
	itemT := slicev.Type().Elem()

	for j := 0; rows.Next(); j++ {
		item, err := mapper(rows)
		if err != nil {
			return err
		}
		//fmt.Printf("--- %#v \n", item)
		val := item.(reflect.Value)
		if itemT.Kind() == reflect.Ptr {
			slicev = reflect.Append(slicev, val)
		} else {
			slicev = reflect.Append(slicev, reflect.Indirect(val))
		}
	}
	if err = rows.Err(); err != nil {
		return err
	}
	//ind.Set(slicev)
	dstv.Elem().Set(slicev)

	return err
}

func (r *Runner) Get(out interface{}, sql string, params ...interface{}) (ok bool, err error) {
	return r.GetContext(context.Background(), out, sql, params...)
}

func (r *Runner) GetContext(ctx context.Context, out interface{}, sql string, params ...interface{}) (ok bool, err error) {
	message := ""

	if r.LoggingEnabled() {
		defer func(start time.Time) {
			r.Logger().Log(&QueryStatus{
				Message: &message,
				Query:   sql,
				Args:    params,
				Err:     err,
				Start:   start,
				End:     time.Now(),
				Context: ctx,
			})
		}(time.Now())
	}
	//stmt, err := r.Prepare(sql)
	//if err != nil {
	//	return false, err
	//}
	//rows, err := stmt.QueryContext(ctx, params...)

	rows, err := r.sqlExecutor.QueryContext(ctx, sql, params...)
	if err != nil {
		return false, err
	}

	defer rows.Close()
	if rows.Next() {
		out, err = r.rowsMapping(out, rows)
		if err != nil {

			return false, err
		}
	} else {
		return false, err
	}

	if err = rows.Err(); err != nil {
		return false, err
	}

	if rows.Next() {
		message = "warn: has more data."
		if err = rows.Close(); err != nil {
			r.Log(&QueryStatus{
				Err:     err,
				Message: &message,
			})
		}
	}
	return true, err
}

func (r *Runner) GetOne(out interface{}) (ok bool, err error) {
	return r.GetOneContext(context.Background(), out)
}

func (r *Runner) GetOneContext(ctx context.Context, out interface{}) (ok bool, err error) {
	entity, ind := r.GetEntity(out)
	//fmt.Printf("%+v  \n", model)
	names := ""
	whereArgs := make([]interface{}, 0)
	wheres := make([]string, 0)
	for _, fd := range entity.Columns {
		if fd.Field.Anonymous || fd.Embedded {
			continue
		}

		lastv := ind.FieldByIndex(fd.Index)
		names += "`" + fd.ColumnName + "`,"

		if lastv.Kind() == reflect.Ptr && lastv.IsNil() {
			continue
		}
		if !lastv.IsValid() {
			continue
		}
		if fd.ColumnName == "" || fd.Omitempty {
			continue
		}

		if lastv.Interface() == fd.Zero.Interface() {
			continue
		}
		if fd.IsUnique {
			wheres = append(wheres, "`"+fd.ColumnName+"`=? ")
			whereArgs = append(whereArgs, lastv.Interface())
		}
		if fd.IsPk {
			wheres = append(wheres, "`"+fd.ColumnName+"`=? ")
			whereArgs = append(whereArgs, lastv.Interface())
		}

	}
	where := strings.Join(wheres, " and ")
	query := fmt.Sprintf("select %s from `%s` where %s", names[0:len(names)-1], entity.TableName, where)

	if len(wheres) == 0 {
		err := errors.New("no unique column for db tag. example: `db:\"order_id,unique\"` : " + ind.Kind().String())
		if r.LoggingEnabled() {
			defer func(start time.Time) {
				r.Logger().Log(&QueryStatus{
					Query:   query,
					Args:    whereArgs,
					Err:     err,
					Start:   start,
					End:     time.Now(),
					Context: ctx,
				})
			}(time.Now())
		}
		return false, err
	}

	return r.GetContext(ctx, out, query, whereArgs...)
}

func (r *Runner) Query(sql string, params ...interface{}) (*sql.Rows, error) {
	return r.QueryContext(context.Background(), sql, params...)
	//stmt, err := r.Prepare(sql)
	//if err != nil {
	//	return nil, err
	//}
	//defer stmt.Close()
	//return stmt.Query(params...)
}

func (r *Runner) QueryContext(ctx context.Context, sql string, params ...interface{}) (rows *sql.Rows, err error) {
	if r.LoggingEnabled() {
		defer func(start time.Time) {
			r.Logger().Log(&QueryStatus{
				Query:   sql,
				Args:    params,
				Err:     err,
				Start:   start,
				End:     time.Now(),
				Context: ctx,
			})
		}(time.Now())
	}

	return r.sqlExecutor.QueryContext(ctx, sql, params...)
	//stmt, err := r.Prepare(sql)
	//if err != nil {
	//	return nil, err
	//}
	//defer stmt.Close()
	//return stmt.QueryContext(ctx, params...)
}

func (r *Runner) QueryRow(sql string, params ...interface{}) *sql.Row {
	return r.QueryRowContext(context.Background(), sql, params...)
}

func (r *Runner) QueryRowContext(ctx context.Context, sql string, params ...interface{}) *sql.Row {
	var err error
	if r.LoggingEnabled() {
		defer func(start time.Time) {
			r.Logger().Log(&QueryStatus{
				Query:   sql,
				Args:    params,
				Err:     err,
				Start:   start,
				End:     time.Now(),
				Context: ctx,
			})
		}(time.Now())
	}
	return r.sqlExecutor.QueryRowContext(ctx, sql, params...)

	//if err != nil {
	//	return nil
	//}
	//defer stmt.Close()
	//row := stmt.QueryRowContext(ctx, params...)
	//return row
}
