package dbx

import (
	"context"
	"database/sql"
	"errors"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tietang/dbx/mapping"
)

// fakeExecutor 捕获 Exec 的 SQL 和参数，不连接任何数据库
type fakeExecutor struct {
	execSQL  string
	execArgs []interface{}
	execErr  error
}

func (f *fakeExecutor) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return nil, errors.New("not implemented")
}
func (f *fakeExecutor) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return nil, errors.New("not implemented")
}
func (f *fakeExecutor) QueryRow(query string, args ...interface{}) *sql.Row { return nil }
func (f *fakeExecutor) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return nil
}
func (f *fakeExecutor) Exec(query string, args ...interface{}) (sql.Result, error) {
	f.execSQL = query
	f.execArgs = args
	if f.execErr != nil {
		return nil, f.execErr
	}
	return fakeResult{lastInsertId: 1, rowsAffected: 1}, nil
}
func (f *fakeExecutor) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return f.Exec(query, args...)
}
func (f *fakeExecutor) Prepare(query string) (*sql.Stmt, error) {
	return nil, errors.New("not implemented")
}
func (f *fakeExecutor) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	return nil, errors.New("not implemented")
}

type fakeResult struct {
	lastInsertId int64
	rowsAffected int64
}

func (r fakeResult) LastInsertId() (int64, error) { return r.lastInsertId, nil }
func (r fakeResult) RowsAffected() (int64, error) { return r.rowsAffected, nil }

func newTestRunner(exec *fakeExecutor) *Runner {
	r := &Runner{}
	r.sqlExecutor = exec
	r.EntityMapper = mapping.NewEntityMapper()
	r.LoggerSettings = defaultLoggerSettings
	r.ILogger = &defaultLogger{}
	r.SetLogging(false)
	return r
}

type txUser struct {
	Id      int64  `db:"id,pk"`
	Name    string `db:"name"`
	Email   string `db:"email,unique"`
	Age     int
	Skipped string `db:"-"`
}

// splitColumns 从 SQL 中提取列名集合（排序后返回），规避 map 遍历导致的列顺序不确定
func splitColumns(t *testing.T, query string) []string {
	t.Helper()
	start := strings.Index(query, "(")
	end := strings.Index(query, ")")
	require.GreaterOrEqual(t, end, start, "no columns in sql: %q", query)
	require.GreaterOrEqual(t, start, 0, "no columns in sql: %q", query)
	var cols []string
	for _, c := range strings.Split(query[start+1:end], ",") {
		cols = append(cols, strings.ReplaceAll(strings.TrimSpace(c), "`", ""))
	}
	sort.Strings(cols)
	return cols
}

func TestInsertContext_SQL(t *testing.T) {
	exec := &fakeExecutor{}
	r := newTestRunner(exec)

	u := &txUser{Name: "tom", Email: "tom@x.com", Age: 20}
	_, err := r.InsertContext(context.Background(), u)
	require.NoError(t, err)

	assert.True(t, strings.HasPrefix(exec.execSQL, "insert into `tx_user`("), "sql = %q", exec.execSQL)
	assert.Equal(t, 3, strings.Count(exec.execSQL, "?"), "want 3 placeholders, sql = %q", exec.execSQL)
	assert.Equal(t, []string{"age", "email", "name"}, splitColumns(t, exec.execSQL))
	assert.Len(t, exec.execArgs, 3)
}

func TestInsertContext_SkipZeroWithOmitempty(t *testing.T) {
	exec := &fakeExecutor{}
	r := newTestRunner(exec)

	// 所有可插入列均为零值（pk 之外），应返回 error 而不是 panic
	type rec struct {
		Id   int64  `db:"id,pk"`
		Name string `db:",omitempty"`
	}
	_, err := r.InsertContext(context.Background(), &rec{Name: ""})
	require.Error(t, err)
	assert.Empty(t, exec.execSQL, "should not execute sql")
}

func TestInsertOrIgnoreContext_SQL(t *testing.T) {
	exec := &fakeExecutor{}
	r := newTestRunner(exec)

	u := &txUser{Name: "tom"}
	_, err := r.InsertOrIgnoreContext(context.Background(), u)
	require.NoError(t, err)

	assert.True(t, strings.HasPrefix(exec.execSQL, "insert or IGNORE  into `tx_user`("), "sql = %q", exec.execSQL)
	assert.Equal(t, 3, strings.Count(exec.execSQL, "?"), "sql = %q", exec.execSQL)
	assert.Equal(t, []string{"age", "email", "name"}, splitColumns(t, exec.execSQL))
}

func TestInsertOrReplaceContext_SQL(t *testing.T) {
	exec := &fakeExecutor{}
	r := newTestRunner(exec)

	_, err := r.InsertOrReplaceContext(context.Background(), &txUser{Name: "tom"})
	require.NoError(t, err)

	assert.True(t, strings.HasPrefix(exec.execSQL, "insert or REPLACE  into `tx_user`("), "sql = %q", exec.execSQL)
	assert.Equal(t, 3, strings.Count(exec.execSQL, "?"), "sql = %q", exec.execSQL)
	assert.Equal(t, []string{"age", "email", "name"}, splitColumns(t, exec.execSQL))
}

func TestUpsertContext_UniquesParam(t *testing.T) {
	exec := &fakeExecutor{}
	r := newTestRunner(exec)

	// 实体没有 upsert tag，通过 uniques 参数指定冲突键
	_, err := r.UpsertContext(context.Background(),
		&txUser{Id: 1, Name: "tom", Email: "tom@x.com"}, "email")
	require.NoError(t, err)

	assert.True(t, strings.HasPrefix(exec.execSQL, "insert into `tx_user`("), "sql = %q", exec.execSQL)
	assert.Contains(t, exec.execSQL, "ON CONFLICT(email)")
	assert.Contains(t, exec.execSQL, "DO UPDATE SET")
	assert.Equal(t, []string{"age", "email", "name"}, splitColumns(t, exec.execSQL))
	// insert 参数 + update 参数 = 2 倍
	assert.Len(t, exec.execArgs, 6)
}

func TestUpsertContext_UpsertTag(t *testing.T) {
	exec := &fakeExecutor{}
	r := newTestRunner(exec)

	type kv struct {
		K   string `db:"k,upsertkey"`
		Val string `db:"val"`
	}
	// 实体有 upsert tag：tag 字段也参与 insert/update（保留既有行为），冲突键取自 tag
	_, err := r.UpsertContext(context.Background(), &kv{K: "a", Val: "b"})
	require.NoError(t, err)

	assert.True(t, strings.HasPrefix(exec.execSQL, "insert into `kv`("), "sql = %q", exec.execSQL)
	assert.Contains(t, exec.execSQL, "ON CONFLICT(k)")
	assert.Contains(t, exec.execSQL, "DO UPDATE SET")
	assert.Equal(t, []string{"k", "val"}, splitColumns(t, exec.execSQL))
}

func TestUpsertContext_NoKeys_Error(t *testing.T) {
	exec := &fakeExecutor{}
	r := newTestRunner(exec)

	// 无 upsert tag 也未传 uniques，必须返回 error 而不是 panic
	_, err := r.UpsertContext(context.Background(), &txUser{Name: "tom"})
	assert.Error(t, err)
}

func TestUpsertContext_EmptyColumns_Error(t *testing.T) {
	exec := &fakeExecutor{}
	r := newTestRunner(exec)

	// 所有列都是 pk，names 为空时不能 panic
	type allpk struct {
		Id int64 `db:"id,pk"`
	}
	_, err := r.UpsertContext(context.Background(), &allpk{Id: 1}, "id")
	assert.Error(t, err)
}

func TestExecute(t *testing.T) {
	exec := &fakeExecutor{}
	r := newTestRunner(exec)

	id, affected, err := r.Execute("update t set a=?", 1)
	require.NoError(t, err)
	assert.Equal(t, int64(1), id)
	assert.Equal(t, int64(1), affected)
}

func TestExecute_Err(t *testing.T) {
	exec := &fakeExecutor{execErr: errors.New("boom")}
	r := newTestRunner(exec)

	_, _, err := r.Execute("update t set a=?", 1)
	assert.Error(t, err)
}

func TestUpdateContext_NoUnique_Error(t *testing.T) {
	exec := &fakeExecutor{}
	r := newTestRunner(exec)

	// 无 pk/unique 字段，必须返回 error 且不执行 SQL
	type plain struct {
		Name string `db:"name"`
	}
	_, err := r.UpdateContext(context.Background(), &plain{Name: "tom"})
	require.Error(t, err)
	assert.Empty(t, exec.execSQL, "should not execute sql")
}

func TestUpdateContext_SQL(t *testing.T) {
	exec := &fakeExecutor{}
	r := newTestRunner(exec)

	u := &txUser{Id: 1, Name: "tom"}
	_, err := r.UpdateContext(context.Background(), u)
	require.NoError(t, err)

	// 零值字段（Email/Age）被跳过；pk 非零时会同时出现在 SET 和 WHERE（既有行为）
	assert.True(t, strings.HasPrefix(exec.execSQL, "update `tx_user` set "), "sql = %q", exec.execSQL)
	assert.Contains(t, exec.execSQL, "where")
	assert.Equal(t, 3, strings.Count(exec.execSQL, "?"), "want set+where placeholders, sql = %q", exec.execSQL)
	// SET 子句应包含非零值列 name（pk 非零时也在 SET 中，既有行为）
	assert.Contains(t, exec.execSQL, "`name`=?")
	assert.Len(t, exec.execArgs, 3)
}

func TestGetOneContext_NoUnique_Error(t *testing.T) {
	exec := &fakeExecutor{}
	r := newTestRunner(exec)

	type plain struct {
		Name string `db:"name"`
	}
	_, err := r.GetOneContext(context.Background(), &plain{Name: "tom"})
	assert.Error(t, err)
}

func TestFindExampleContext_NoCondition_Error(t *testing.T) {
	exec := &fakeExecutor{}
	r := newTestRunner(exec)

	// 全部字段为零值，无可用查询条件
	var out []*txUser
	err := r.FindExampleContext(context.Background(), &txUser{}, &out)
	assert.Error(t, err)
}

func TestSelectContext_NotPtr_Error(t *testing.T) {
	exec := &fakeExecutor{}
	r := newTestRunner(exec)

	err := r.SelectContext(context.Background(), func(rows *sql.Rows) (interface{}, error) {
		return nil, nil
	}, txUser{}, "select 1")
	assert.Error(t, err)
}
