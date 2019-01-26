package dbx

import (
	"database/sql"
	"reflect"
)

func (r *Runner) rowsMapping(model interface{}, rows *sql.Rows) (interface{}, error) {
	e, ind := r.GetEntity(model)
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	scans := make([]interface{}, len(columns))

	obj := ind
	if ind.Kind() != reflect.Struct {
		obj = reflect.New(e.StructType)
	}

	for i, c := range columns {

		fi, ok := e.GetFieldModel(c)
		if !ok {
			scans[i] = new(interface{})
			continue
		}
		//v := obj.FieldByIndex(fi.Index)
		v := reflect.Indirect(obj).FieldByIndex(fi.Index)
		//
		//if v.Kind() == reflect.Ptr && v.IsNil() {
		//    //fmt.Println("--ptr:", v.Type(), fi.FieldName)
		//    alloc := reflect.New(reflectx.Deref(v.Type()))
		//    v.Set(alloc)
		//}
		//if v.Kind() == reflect.Map && v.IsNil() {
		//    v.Set(reflect.MakeMap(v.Type()))
		//    //fmt.Println("--map:", v.Type(), fi.FieldName)
		//}
		//if v.Kind() == reflect.Struct {
		//    //fmt.Println("--Struct:", v.Type(), fi.FieldName)
		//    typ := reflectx.Deref(v.Type())
		//    val := reflect.New(typ)
		//    v.Set(val.Elem())
		//}

		scans[i] = v.Addr().Interface()
		//fmt.Printf("%s %#v %#v \n ", v.Type().String(), scans[i], reflect.ValueOf(scans[i]).Kind() == reflect.Ptr)
	}
	//if r.logEnabled{
	//    log.Debug("mapping columns: ",columns)
	//    log.Debugf("mapping columns: %#v",scans)
	//}
	//fmt.Printf("%#v\n", scans)
	//
	if err = rows.Scan(scans...); err != nil {
		return obj, err
	}
	return obj, nil
}
