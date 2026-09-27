package mapping

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type upsertEntity struct {
	Id  int64  `db:"id,pk"`
	Key string `db:"k,upsertKey"` // 回归：以前 "upsertKey" 因大小写 bug 解析不到
	Val string `db:"val,unique"`
}

type omitemptyEntity struct {
	Id   int64  `db:"id,pk"`
	Name string `db:"user_name,omitempty"`
}

func TestSnakeString(t *testing.T) {
	cases := map[string]string{
		"TxUser":     "tx_user",
		"HTTPCode":   "h_t_t_p_code", // XxYY -> xx_y_y，符合注释描述的规则
		"UserName":   "user_name",
		"user_name":  "user_name",
		"a":          "a",
		"":           "",
		"Already_sn": "already_sn",
	}
	for in, want := range cases {
		assert.Equal(t, want, SnakeString(in), "SnakeString(%q)", in)
	}
}

func TestGetEntity_UpsertKeyTag(t *testing.T) {
	m := NewEntityMapper()
	e, _ := m.GetEntity(&upsertEntity{})

	fd, ok := e.GetFieldModel("k")
	require.True(t, ok, "column k not found")
	assert.True(t, fd.IsUpsertKey, "upsertKey tag should set IsUpsertKey=true")
}

func TestGetEntity_UpsertLowercase(t *testing.T) {
	m := NewEntityMapper()
	e, _ := m.GetEntity(&struct {
		K string `db:"k,upsert"`
	}{})
	fd, ok := e.GetFieldModel("k")
	require.True(t, ok)
	assert.True(t, fd.IsUpsertKey, "lowercase upsert tag should also set IsUpsertKey")
}

func TestGetEntity_PkAndUnique(t *testing.T) {
	m := NewEntityMapper()
	e, _ := m.GetEntity(&upsertEntity{})

	assert.Equal(t, "upsert_entity", e.TableName)
	pk, ok := e.GetFieldModel("id")
	require.True(t, ok)
	assert.True(t, pk.IsPk, "pk field not parsed")
	require.NotNil(t, e.PkField)
	assert.Equal(t, "id", e.PkField.ColumnName)
	require.Len(t, e.UniqueFields, 1)
	assert.Equal(t, "val", e.UniqueFields[0].ColumnName)
}

func TestGetEntity_OmitemptyAndCustomColumn(t *testing.T) {
	m := NewEntityMapper()
	e, _ := m.GetEntity(&omitemptyEntity{})

	fd, ok := e.GetFieldModel("user_name")
	require.True(t, ok, "custom column name not applied")
	assert.True(t, fd.Omitempty, "omitempty tag not parsed")
}

func TestGetEntity_SkipDbDash(t *testing.T) {
	m := NewEntityMapper()
	e, _ := m.GetEntity(&struct {
		Skip    string `db:"-"`
		Keep    string `db:"keep"`
		Unnamed string // 无 tag，按字段名生成列名
	}{})
	_, ok := e.GetFieldModel("skip")
	assert.False(t, ok, `db:"-" should be skipped`)
	_, ok = e.GetFieldModel("keep")
	assert.True(t, ok, "normal column should exist")
	_, ok = e.GetFieldModel("unnamed")
	assert.True(t, ok, "untagged field should map to snake column")
}

func TestGetEntity_RegisterTable(t *testing.T) {
	m := NewEntityMapper()
	m.RegisterTable(&upsertEntity{}, "custom_table")

	typ := reflect.TypeOf(upsertEntity{})
	name, ok := m.GetTableName(typ)
	require.True(t, ok)
	assert.Equal(t, "custom_table", name)

	e, _ := m.GetEntity(&upsertEntity{})
	assert.Equal(t, "custom_table", e.TableName)
}

func TestGetEntity_Cache(t *testing.T) {
	m := NewEntityMapper()
	e1, _ := m.GetEntity(&upsertEntity{})
	e2, _ := m.GetEntity(&upsertEntity{})
	assert.Same(t, e1, e2, "GetEntity should return cached EntityInfo for same type")
}
