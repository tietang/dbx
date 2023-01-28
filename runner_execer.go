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

func (r *Runner) Exec(sql string, params ...interface{}) (rs sql.Result, err error) {
	return r.ExecContext(context.Background(), sql, params...)
}

func (r *Runner) Execute(sql string, params ...interface{}) (lastInsertId, rowsAffected int64, err error) {
	rs, err := r.Exec(sql, params...)
	if err != nil {
		return 0, 0, err
	}
	lastInsertId, err = rs.LastInsertId()
	if err != nil {
		return 0, 0, err
	}
	rowsAffected, err = rs.RowsAffected()
	return lastInsertId, rowsAffected, err
}

func (r *Runner) ExecContext(ctx context.Context, sql string, params ...interface{}) (rs sql.Result, err error) {
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
	return r.sqlExecutor.ExecContext(ctx, sql, params...)

	//stmt, err := r.Prepare(sql)
	//if err != nil {
	//	return nil, err
	//}
	//defer stmt.Close()
	//rs, err = stmt.ExecContext(ctx, params...)
	//return rs, err
}

func (r *Runner) Insert(model interface{}) (rs sql.Result, err error) {
	return r.InsertContext(context.Background(), model)
}

func (r *Runner) InsertContext(ctx context.Context, model interface{}) (rs sql.Result, err error) {
	entity, ind := r.GetEntity(model)
	names, placeholders := "", ""
	params := make([]interface{}, 0)

	for _, fd := range entity.Columns {

		if fd.Field.Anonymous || fd.Embedded {
			continue
		}
		lastv := ind.FieldByIndex(fd.Index)
		if lastv.Kind() == reflect.Ptr && lastv.IsNil() {
			continue
		}
		if !lastv.IsValid() {
			continue
		}
		if fd.ColumnName == "" || (fd.Omitempty && lastv.IsZero()) {
			//if fd.ColumnName == "" || fd.Omitempty {
			continue
		}
		if fd.IsPk {
			continue
		}

		params = append(params, lastv.Interface())
		names += "`" + fd.ColumnName + "`,"
		placeholders += "?,"
	}
	sql := fmt.Sprintf("insert into `%s`(%s) values(%s)", entity.TableName, names[0:len(names)-1], placeholders[0:len(placeholders)-1])

	return r.ExecContext(ctx, sql, params...)
}
func (r *Runner) InsertOrIgnore(model interface{}) (rs sql.Result, err error) {
	return r.InsertOrContext(context.Background(), model, INSERT_IGNORE)
}

func (r *Runner) InsertOrIgnoreContext(ctx context.Context, model interface{}) (rs sql.Result, err error) {
	return r.InsertOrContext(ctx, model, INSERT_IGNORE)
}

func (r *Runner) InsertOrReplace(model interface{}) (rs sql.Result, err error) {
	return r.InsertOrContext(context.Background(), model, INSERT_REPLACE)
}

func (r *Runner) InsertOrReplaceContext(ctx context.Context, model interface{}) (rs sql.Result, err error) {
	return r.InsertOrContext(ctx, model, INSERT_REPLACE)
}
func (r *Runner) InsertOr(model interface{}, onConflict string) (rs sql.Result, err error) {
	return r.InsertOrContext(context.Background(), model, onConflict)
}

func (r *Runner) InsertOrContext(ctx context.Context, model interface{}, onConflict string) (rs sql.Result, err error) {
	entity, ind := r.GetEntity(model)
	names, placeholders := "", ""
	params := make([]interface{}, 0)

	for _, fd := range entity.Columns {

		if fd.Field.Anonymous || fd.Embedded {
			continue
		}
		lastv := ind.FieldByIndex(fd.Index)
		if lastv.Kind() == reflect.Ptr && lastv.IsNil() {
			continue
		}
		if !lastv.IsValid() {
			continue
		}
		if fd.ColumnName == "" || (fd.Omitempty && lastv.IsZero()) {
			//if fd.ColumnName == "" || fd.Omitempty {
			continue
		}
		if fd.IsPk {
			continue
		}

		params = append(params, lastv.Interface())
		names += "`" + fd.ColumnName + "`,"
		placeholders += "?,"
	}
	sql := fmt.Sprintf("insert or %s  into `%s`(%s) values(%s)", onConflict, entity.TableName, names[0:len(names)-1], placeholders[0:len(placeholders)-1])

	return r.ExecContext(ctx, sql, params...)
}

func (r *Runner) InsertOrGet(model interface{}) (id int64, err error) {
	return r.InsertOrGetContext(context.Background(), model)
}
func (r *Runner) InsertOrGetContext(ctx context.Context, model interface{}) (id int64, err error) {
	entity, ind := r.GetEntity(model)
	names, placeholders := "", ""
	params := make([]interface{}, 0)

	params, whereArgs := make([]interface{}, 0), make([]interface{}, 0)
	wheres := make([]string, 0)
	allNames := ""
	for _, fd := range entity.Columns {
		if fd.ColumnName == "" {
			continue
		}
		if fd.Field.Anonymous || fd.Embedded {
			continue
		}
		allNames += ",`" + fd.ColumnName + "`"

		lastv := ind.FieldByIndex(fd.Index)
		if lastv.Kind() == reflect.Ptr && lastv.IsNil() {
			continue
		}
		if !lastv.IsValid() {
			continue
		}
		if fd.ColumnName == "" || (fd.Omitempty && lastv.IsZero()) {
			//if fd.ColumnName == "" || fd.Omitempty   {
			continue
		}

		if fd.IsPk {
			continue
		}

		if fd.IsUnique {
			wheres = append(wheres, "`"+fd.ColumnName+"`=? ")
			whereArgs = append(whereArgs, lastv.Interface())
		}

		params = append(params, lastv.Interface())
		names += "`" + fd.ColumnName + "`,"
		placeholders += "?,"

	}
	sql := fmt.Sprintf("insert or %s  into `%s`(%s) values(%s)", INSERT_IGNORE, entity.TableName, names[0:len(names)-1], placeholders[0:len(placeholders)-1])

	rs, err := r.ExecContext(ctx, sql, params...)
	if err != nil {
		return 0, err
	}
	lastId, _ := rs.LastInsertId()
	rows, _ := rs.RowsAffected()
	if rows == 0 && err == nil {
		where := strings.Join(wheres, " or ")

		sql = fmt.Sprintf("select %s from %s where %s", allNames[1:], entity.TableName, where)
		ok, err := r.Get(model, sql, whereArgs...)
		if ok {
			pkFieldName := entity.PkField.FieldName
			val := reflect.ValueOf(model)
			field := val.Elem().FieldByName(pkFieldName)
			id = field.Int()
			return id, err
		} else {
			return id, err
		}
	}
	return lastId, err
}

func (r *Runner) Upsert(model interface{}, uniques ...string) (rs sql.Result, err error) {
	return r.UpsertContext(context.Background(), model, uniques...)
}

func (r *Runner) UpsertContext(ctx context.Context, model interface{}, uniques ...string) (rs sql.Result, err error) {
	entity, ind := r.GetEntity(model)
	names, placeholders := "", ""
	params := make([]interface{}, 0)
	//
	namesForUpdate, upsertKeys := "", ""
	for _, fd := range entity.Columns {

		if fd.Field.Anonymous || fd.Embedded {
			continue
		}
		lastv := ind.FieldByIndex(fd.Index)
		if lastv.Kind() == reflect.Ptr && lastv.IsNil() {
			continue
		}
		if !lastv.IsValid() {
			continue
		}
		if fd.ColumnName == "" || (fd.Omitempty && lastv.IsZero()) {
			//if fd.ColumnName == "" || fd.Omitempty {
			continue
		}
		if fd.IsPk {
			continue
		}

		params = append(params, lastv.Interface())
		names += ",`" + fd.ColumnName + "`"
		placeholders += ",?"
		//
		namesForUpdate += ",`" + fd.ColumnName + "`=?"
		if fd.IsUpsertKey {
			//wheres = append(wheres, "`"+fd.ColumnName+"`=? ")
			//whereArgs = append(whereArgs, lastv.Interface())
			upsertKeys += "," + fd.ColumnName
		}
	}
	if len(uniques) == 0 && upsertKeys == "" {
		return nil, errors.New("未指定upsert Unique Keys！")
	}
	//sql := fmt.Sprintf("insert into `%s`(%s) values(%s)", entity.TableName, names[0:len(names)-1], placeholders[0:len(placeholders)-1])
	//where := strings.Join(wheres, " and ")
	//sql = fmt.Sprintf("update `%s` set %s where %s", entity.TableName, names[0:len(names)-1], where)
	sql := fmt.Sprintf(
		"insert into `%s`(%s) values(%s)  ON CONFLICT(%s)  DO UPDATE SET  %s ",
		entity.TableName, names[1:], placeholders[1:],
		upsertKeys[1:], namesForUpdate[1:])
	args := append(params, params...)
	return r.ExecContext(ctx, sql, args...)
}

func (r *Runner) Update(model interface{}) (rs sql.Result, err error) {
	return r.UpdateContext(context.Background(), model)
}

func (r *Runner) UpdateContext(ctx context.Context, model interface{}) (rs sql.Result, err error) {
	entity, ind := r.GetEntity(model)
	var sql string
	//fmt.Printf("%+v  \n", model)
	names := ""
	params, whereArgs := make([]interface{}, 0), make([]interface{}, 0)
	wheres := make([]string, 0)
	for _, fd := range entity.Columns {

		if fd.Field.Anonymous || fd.Embedded {
			continue
		}
		lastv := ind.FieldByIndex(fd.Index)

		if lastv.Kind() == reflect.Ptr && lastv.IsNil() {
			continue
		}
		if !lastv.IsValid() {
			continue
		}
		//if fd.ColumnName == "" || fd.Omitempty {
		if fd.ColumnName == "" || (fd.Omitempty && lastv.IsZero()) {
			continue
		}

		if lastv.Interface() == fd.Zero.Interface() {
			continue
		}

		params = append(params, lastv.Interface())
		names += "`" + fd.ColumnName + "`=?,"
		if fd.IsUnique {
			wheres = append(wheres, "`"+fd.ColumnName+"`=? ")
			whereArgs = append(whereArgs, lastv.Interface())
		}
		if fd.IsPk {
			wheres = append(wheres, "`"+fd.ColumnName+"`=? ")
			whereArgs = append(whereArgs, lastv.Interface())
		}

	}
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
	if len(wheres) == 0 {
		err = errors.New("no.txt unique column config for db tag.")
		return nil, err
	}
	where := strings.Join(wheres, " and ")
	params = append(params, whereArgs...)

	sql = fmt.Sprintf("update `%s` set %s where %s", entity.TableName, names[0:len(names)-1], where)
	return r.ExecContext(ctx, sql, params...)
}
