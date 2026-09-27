package it

import (
	"context"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCount(t *testing.T) {
	g0 := newEnvelopeGoods(1)
	rs, err := db.InsertContext(context.Background(), &g0)
	require.NoError(t, err)
	require.NotNil(t, rs)
	id, err := rs.LastInsertId()
	require.NoError(t, err)
	assert.True(t, id > 0)
	rows, err := rs.RowsAffected()
	require.NoError(t, err)
	assert.Equal(t, int64(1), rows)

	query := "select count(1) from red_envelope_goods where id>=?"
	var x int64
	row := db.QueryRow(query, id)
	err = row.Scan(&x)

	require.NoError(t, err)
	assert.True(t, x >= 1)
	fmt.Println(x)
}

func TestCountInt32(t *testing.T) {
	g0 := newEnvelopeGoods(1)
	rs, err := db.InsertContext(context.Background(), &g0)
	require.NoError(t, err)
	require.NotNil(t, rs)
	id, err := rs.LastInsertId()
	require.NoError(t, err)
	assert.True(t, id > 0)
	rows, err := rs.RowsAffected()
	require.NoError(t, err)
	assert.Equal(t, int64(1), rows)

	query := "select count(1) from red_envelope_goods where id>=?"
	x, err := db.GetInt32(query, id)
	require.NoError(t, err)
	assert.True(t, x >= 1)
	fmt.Println(x)
}

func TestCountInt64(t *testing.T) {
	g0 := newEnvelopeGoods(1)
	rs, err := db.InsertContext(context.Background(), &g0)
	require.NoError(t, err)
	require.NotNil(t, rs)
	id, err := rs.LastInsertId()
	require.NoError(t, err)
	assert.True(t, id > 0)
	rows, err := rs.RowsAffected()
	require.NoError(t, err)
	assert.Equal(t, int64(1), rows)

	query := "select count(1) from red_envelope_goods where id>= ?"
	x, err := db.GetInt64(query, id)
	require.NoError(t, err)
	assert.True(t, x >= 1)
	fmt.Println(x)
}
