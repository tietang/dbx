package mapping

import (
	"fmt"
	"github.com/tietang/dbx/reflectx"
	"reflect"
	"strings"
)

type FieldModel struct {
	reflectx.FieldInfo
	ColumnName      string
	TypeFullName    string
	FieldType       reflect.StructField
	FieldName       string
	ParentFieldName string
	ParentFieldType reflect.StructField
	FieldIndex      []int
	Omitempty       bool
	EmbeddedStruct  bool
	IsUnique        bool
	IsPk            bool
}

type EntityInfo struct {
	TableName    string
	TypeFullName string
	StructType   reflect.Type
	//
	FieldModels []*FieldModel
	Columns     map[string]*FieldModel
	//Fields       map[string]*FieldModel
	PkField      *FieldModel
	UniqueFields []*FieldModel
	//
	isDbInit bool
}

func newEntity() *EntityInfo {
	f := &EntityInfo{
		FieldModels:  make([]*FieldModel, 0),
		UniqueFields: make([]*FieldModel, 0),
		Columns:      make(map[string]*FieldModel),
		//Fields:      make(map[string]*FieldModel),
	}
	return f

}

func (e *EntityInfo) Print() {

	for key, value := range e.FieldModels {
		fmt.Printf("%v %#v\n", key, value)
	}
}

func (f *EntityInfo) Append(fd *FieldModel) {
	f.FieldModels = append(f.FieldModels, fd)

	c, ok := f.Columns[fd.ColumnName]
	if ok {
		if len(c.Index) > len(fd.Index) {
			f.Columns[fd.ColumnName] = fd
		}
	} else {
		f.Columns[fd.ColumnName] = fd
	}
	//f.Fields[fd.FieldName] = fd
	if fd.IsUnique {
		f.UniqueFields = append(f.UniqueFields, fd)
	}
	if fd.IsPk {
		f.PkField = fd
	}
}

func (f *EntityInfo) GetFieldModel(columnName string) (*FieldModel, bool) {
	fi, ok := f.Columns[columnName]
	return fi, ok
}

//
//var mapper = reflectx.NewMapper("db")
//var lock sync.Mutex
//
//var tableNames = make(map[reflect.Type]string)
//
//func RegisterTable(model interface{}, tableName string) {
//	val := reflect.ValueOf(model)
//	ind := reflect.Indirect(val)
//	typ := ind.Type()
//	//读取slice类型
//	if ind.Kind() == reflect.Slice {
//		etyp := ind.Type().Elem()
//		typ = etyp
//		if typ.Kind() == reflect.Ptr {
//			typ = typ.Elem()
//		}
//	}
//	tableNames[typ] = tableName
//}
//
//func GetTableName(typ reflect.Type) (string, bool) {
//	b, ok := tableNames[typ]
//	return b, ok
//}
//
//func GetEntity(model interface{}) (*EntityInfo, reflect.Value) {
//	val := reflect.ValueOf(model)
//	ind := reflect.Indirect(val)
//	return GetEntityByValue(ind), ind
//}
//
//func GetEntityByValue(ind reflect.Value) *EntityInfo {
//	typ := ind.Type()
//
//	//读取slice类型
//	if ind.Kind() == reflect.Slice {
//		etyp := ind.Type().Elem()
//		typ = etyp
//		if typ.Kind() == reflect.Ptr {
//			typ = typ.Elem()
//		}
//	}
//	//if typ.Kind() != reflect.Struct {
//	//}
//	//&& ind.Kind() != reflect.Map
//	e, ok := EntityCache.Get(typ.String())
//	if ok && e.isDbInit {
//		return e
//	}
//	lock.Lock()
//	if !ok {
//		e = reflect(typ, "", "")
//	}
//
//	lock.Unlock()
//
//	return e
//
//}
//
//func reflect(typ reflect.Type, prefix, suffix string) *EntityInfo {
//	//ind := reflect.Indirect(val)
//	//typ := ind.Type()
//
//	//if typ.Kind() != reflect.Ptr {
//	//    panic(fmt.Errorf(" cannot use non-ptr entity struct `%s`", typ.String()))
//	//}
//	//// For this case:
//	//// u := &User{}
//	//// registerModel(&u)
//	if typ.Kind() == reflect.Ptr {
//		panic(fmt.Errorf(" only allow ptr entity struct, it looks you use two reference to the struct `%s`", typ))
//	}
//	fmt.Println(tableNames)
//	e := newEntity()
//	e.TypeFullName = typ.String()
//
//	e.StructType = typ
//	tb, ok := GetTableName(typ)
//	if ok {
//		e.TableName = tb
//	} else {
//		e.TableName = SnakeString(typ.Name())
//	}
//
//	if prefix != "" {
//		e.TableName = prefix + "_" + e.TableName
//	}
//	if suffix != "" {
//		e.TableName = e.TableName + "_" + suffix
//	}
//
//	//switch typ.Kind() {
//	//case reflect.Struct:
//	//    structMap := mapper.TypeMap(typ)
//	//    fmt.Println(structMap)
//	//
//	//case reflect.Map:
//	//
//	//default:
//	//
//	//}
//
//	smap := mapper.TypeMap(typ)
//	expandFields(e, smap.Index)
//
//	//getFeilds(e, typ, false)
//
//	EntityCache.Set(e.TypeFullName, e)
//
//	return e
//}

func expandFields(e *EntityInfo, fis []*reflectx.FieldInfo) {
	//for index, fi := range fis {
	//	fmt.Printf("1: %d %+v %+v %+v\n", index, fi.Field.Name, fi.Name, fi)
	//}
	for _, fi := range fis {

		//if fi.Field.Anonymous {
		//	if len(fi.Children) > 0 {
		//		expandFields(e, fi.Children)
		//	}
		//	continue
		//}
		//如果index路径2个以上，并且其父Field不为嵌入式struct，那么该Field就为结构体字段，忽略
		if len(fi.Index) > 1 && !fi.Parent.Embedded && !fi.Field.Anonymous {
			continue
		}
		if strings.Contains(fi.Path, ".") {
			continue
		}
		//如果是嵌入式结构体，忽略结构体本身
		if fi.Embedded || fi.Field.Anonymous {
			continue
		}

		fd := &FieldModel{}

		fd.FieldInfo = *fi
		fd.FieldName = fi.Field.Name
		if fi.Field.Name == "" {
			continue
		}
		//if fi.Name == "" {
		//	continue
		//}

		fd.ColumnName = SnakeString(fd.FieldName)
		fd.FieldType = fi.Field
		fd.TypeFullName = fd.Field.Type.Name()

		tag := fi.Field.Tag
		tval := tag.Get("db")
		if tval != "" {

			tvals := strings.Split(tval, ",")
			for k, v := range tvals {
				if v == "-" {
					continue
				}
				if strings.ToLower(v) == "omitempty" {
					fd.Omitempty = true
				}
				if k == 0 {
					if v != "" {
						fd.ColumnName = v
					}
				} else {
					fi.Options[v] = v
				}
				if strings.ToLower(v) == "pk" || strings.ToLower(v) == "id" { //
					fd.IsPk = true
				}
				if strings.ToLower(v) == "uni" || strings.ToLower(v) == "unique" {
					fd.IsUnique = true
				}
			}
		}
		if fi.Field.Type.Kind() == reflect.Struct {
			fd.EmbeddedStruct = true
		}
		if fi.Field.Anonymous {
			fd.Omitempty = true
		}
		//fmt.Printf("2: %d %+v %+v %+v\n", index, fi.Field.Name, fi.Name, fi)

		fd.FieldIndex = fi.Index
		fd.ParentFieldType = fi.Parent.Field
		e.Append(fd)

	}
}

func getFeilds(entity *EntityInfo, ind reflect.Type, anonymous bool) {
	var (
		//err error
		sf reflect.StructField
	)

	for i := 0; i < ind.NumField(); i++ {
		sf = ind.Field(i)
		//sf = ind.Type().FieldModel(i)
		if sf.Type.Kind() == reflect.Struct || sf.Anonymous {
			getFeilds(entity, sf.Type, true)
		} else {
			fi := GetField(sf, anonymous)
			fi.ParentFieldName = sf.Name
			if fi != nil {
				entity.Append(fi)
			}

		}

	}

}

func GetField(sf reflect.StructField, anonymous bool) *FieldModel {
	fi := &FieldModel{}
	fi.Options = make(map[string]string)
	name := sf.Name
	fi.FieldName = sf.Name
	fi.ColumnName = SnakeString(name)
	fi.FieldType = sf
	fi.TypeFullName = sf.Type.Name()
	tag := sf.Tag
	tval := tag.Get("db")
	if tval != "" {
		tvals := strings.Split(tval, ",")
		for k, v := range tvals {
			if v == "-" {
				continue
			}
			if strings.ToLower(v) == "omitempty" {
				fi.Omitempty = true
			}
			if k == 0 {
				if v != "" {
					fi.ColumnName = v
				}
			} else {
				fi.Options[v] = v
			}
			if strings.ToLower(v) == "pk" || strings.ToLower(v) == "id" { //
				fi.IsPk = true
			}
			if strings.ToLower(v) == "uni" || strings.ToLower(v) == "unique" {
				fi.IsUnique = true
			}
		}
	}
	fi.Embedded = anonymous
	fi.FieldIndex = sf.Index

	//fmt.Println("####", index, tval, tvals)

	return fi
}
