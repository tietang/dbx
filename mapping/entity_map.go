package mapping

import (
    "fmt"
    "github.com/tietang/dbx/reflectx"
    "reflect"
    "sync"
)

type EntityMapper interface {
    RegisterTable(model interface{}, tableName string)
    GetEntity(model interface{}) (*EntityInfo, reflect.Value)
    GetTableName(typ reflect.Type) (tableName string, ok bool)
    GetMapper() *reflectx.Mapper
}

type entityMapper struct {
    mapper      *reflectx.Mapper
    lock        sync.Mutex
    tableNames  map[reflect.Type]string
    entityCache EntityCache
}

func NewEntityMapper() EntityMapper {
    e := &entityMapper{
        tableNames: make(map[reflect.Type]string),
        mapper:     reflectx.NewMapper("db"),
    }
    return e
}

func (e *entityMapper) GetMapper() *reflectx.Mapper {
    return e.mapper
}

func (e *entityMapper) RegisterTable(model interface{}, tableName string) {
    val := reflect.ValueOf(model)
    ind := reflect.Indirect(val)
    typ := ind.Type()
    //读取slice类型
    if ind.Kind() == reflect.Slice {
        etyp := ind.Type().Elem()
        typ = etyp
        if typ.Kind() == reflect.Ptr {
            typ = typ.Elem()
        }
    }
    e.tableNames[typ] = tableName
    e.GetEntity(model)
}

func (e *entityMapper) GetTableName(typ reflect.Type) (string, bool) {
    b, ok := e.tableNames[typ]
    return b, ok
}

func (e *entityMapper) GetEntity(model interface{}) (*EntityInfo, reflect.Value) {
    val := reflect.ValueOf(model)
    ind := reflect.Indirect(val)
    return e.getEntityByValue(ind), ind
}

func (e *entityMapper) getEntityByValue(ind reflect.Value) *EntityInfo {
    typ := ind.Type()

    //读取slice类型
    if ind.Kind() == reflect.Slice {
        etyp := ind.Type().Elem()
        typ = etyp
        if typ.Kind() == reflect.Ptr {
            typ = typ.Elem()
        }
    }
    ei, ok := e.entityCache.Get(typ.String())
    if ok && ei.isDbInit {
        return ei
    }
    e.lock.Lock()
    if !ok {
        ei = e.reflect(typ, "", "")
    }

    e.lock.Unlock()

    return ei

}

func (e *entityMapper) reflect(typ reflect.Type, prefix, suffix string) *EntityInfo {
    if typ.Kind() == reflect.Ptr {
        panic(fmt.Errorf(" only allow ptr entity struct, it looks you use two reference to the struct `%s`", typ))
    }
    //fmt.Println(e.tableNames)
    ei := newEntity()
    ei.TypeFullName = typ.String()

    ei.StructType = typ
    tb, ok := e.GetTableName(typ)
    if ok {
        ei.TableName = tb
    } else {
        ei.TableName = SnakeString(typ.Name())
    }

    if prefix != "" {
        ei.TableName = prefix + "_" + ei.TableName
    }
    if suffix != "" {
        ei.TableName = ei.TableName + "_" + suffix
    }

    smap := e.mapper.TypeMap(typ)
    expandFields(ei, smap.Index)

    //getFeilds(e, typ, false)

    e.entityCache.Set(ei.TypeFullName, ei)

    return ei
}
