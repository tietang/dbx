package it

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/segmentio/ksuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSlice(t *testing.T) {
	//redEnvelopeGoodss := make([]EnvelopeGoods, 0)
	var redEnvelopeGoodss []EnvelopeGoods
	dstv := reflect.ValueOf(redEnvelopeGoodss)
	fmt.Println(dstv.Type(), dstv.Kind(), dstv.Interface())
	ind := reflect.Indirect(dstv)
	fmt.Println(ind.Type(), ind.Kind())
	typ := ind.Type()
	if ind.Kind() == reflect.Slice {
		etyp := ind.Type().Elem()
		typ = etyp
		if typ.Kind() == reflect.Ptr {
			typ = typ.Elem()
		}
	}
	fmt.Println(typ)
}

func TestFindPtr(t *testing.T) {
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

	var redEnvelopeGoodss []EnvelopeGoods
	query := "select * from red_envelope_goods limit ?,?"
	err = db.Find(&redEnvelopeGoodss, query, 0, 100)

	require.NoError(t, err)
	assert.True(t, len(redEnvelopeGoodss) > 0)
	fmt.Println(redEnvelopeGoodss)
}

func TestFindPtrSlice(t *testing.T) {
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

	var redEnvelopeGoodss []*EnvelopeGoods
	query := "select * from red_envelope_goods limit ?,?"
	err = db.Find(&redEnvelopeGoodss, query, 0, 100)

	require.NoError(t, err)
	assert.True(t, len(redEnvelopeGoodss) > 0)
	fmt.Println(redEnvelopeGoodss)
}

func TestFind(t *testing.T) {
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

	var redEnvelopeGoodss []EnvelopeGoods
	query := "select * from red_envelope_goods limit ?,?"
	err = db.Find(redEnvelopeGoodss, query, 0, 100)

	assert.Error(t, err)
}

func TestGetAccount(t *testing.T) {
	a := &Account{
		Balance: decimal.NewFromFloat(100), Status: 1,
		UserId:    ksuid.New().Next().String(),
		Username:  sql.NullString{String: "测试用户1", Valid: true},
		AccountNo: ksuid.New().Next().String(),
	}
	rs, err := db.Insert(a)
	require.NoError(t, err)
	fmt.Printf("AccountNo=%s\n", a.AccountNo)
	require.NotNil(t, rs)
	q := &Account{AccountNo: a.AccountNo}
	ok, err := db.GetOne(q)
	require.NoError(t, err)
	assert.True(t, ok)
	assert.Equal(t, a.AccountNo, q.AccountNo)
	assert.Equal(t, a.Balance.String(), q.Balance.String())
	assert.Equal(t, a.UserId, q.UserId)
	assert.NotNil(t, q.CreatedAt)
	assert.NotNil(t, q.UpdatedAt)
}

func TestGetAccountById(t *testing.T) {
	fmt.Println(Create)
	a := &Account{
		Balance: decimal.NewFromFloat(100), Status: 1,
		UserId:    ksuid.New().Next().String(),
		Username:  sql.NullString{String: "测试用户1", Valid: true},
		AccountNo: ksuid.New().Next().String(),
	}
	rs, err := db.Insert(a)
	require.NoError(t, err)
	fmt.Printf("AccountNo=%s\n", a.AccountNo)
	require.NotNil(t, rs)
	id, err := rs.LastInsertId()
	require.NoError(t, err)
	assert.GreaterOrEqual(t, id, int64(0))
	q := &Account{Id: id}
	ok, err := db.GetOne(q)
	require.NoError(t, err)
	assert.True(t, ok)
	assert.Equal(t, id, q.Id)
	assert.Equal(t, a.AccountNo, q.AccountNo)
	assert.Equal(t, a.Balance.String(), q.Balance.String())
	assert.Equal(t, a.UserId, q.UserId)
	assert.NotNil(t, q.CreatedAt)
	assert.NotNil(t, q.UpdatedAt)
}

func TestInsert(t *testing.T) {
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
	g1 := &EnvelopeGoods{EnvelopeNo: g0.EnvelopeNo}

	ok, err := db.GetOne(g1)
	fmt.Println()
	fmt.Printf("%+v\n", g0)
	fmt.Printf("%+v\n", g1)
	require.NoError(t, err)
	assert.True(t, ok)
	assert.Equal(t, *g0.EnvelopeNo, *g1.EnvelopeNo)
	assert.Equal(t, g0.EnvelopeType, g1.EnvelopeType)
	assert.Equal(t, g0.Username, g1.Username)
	assert.Equal(t, g0.UserId, g1.UserId)
	assert.Equal(t, g0.Blessing.String, g1.Blessing.String)
	assert.Equal(t, g0.Amount.String(), g1.Amount.String())
	assert.Equal(t, g0.AmountOne.String(), g1.AmountOne.String())
	assert.Equal(t, g0.Quantity, g1.Quantity)
	assert.Equal(t, g0.RemainAmount.String(), g1.RemainAmount.String())
	assert.Equal(t, g0.RemainQuantity, g1.RemainQuantity)
	assert.Equal(t, g0.ExpiredAt.Second(), g1.ExpiredAt.Second())
	assert.Equal(t, g0.Status, g1.Status)
	assert.Equal(t, g0.OrderType, g1.OrderType)
	assert.Equal(t, g0.PayStatus, g1.PayStatus)
}

func TestUpdate(t *testing.T) {
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
	g1 := newEnvelopeGoods(2)
	g1.EnvelopeNo = g0.EnvelopeNo
	//
	rs, err = db.Update(&g1)
	fmt.Println()
	fmt.Printf("%+v\n", g0)
	fmt.Printf("%+v\n", g1)
	require.NoError(t, err)
	require.NotNil(t, rs)
	id, err = rs.LastInsertId()
	require.NoError(t, err)
	assert.Equal(t, int64(0), id)
	rows, err = rs.RowsAffected()
	require.NoError(t, err)
	assert.Equal(t, int64(1), rows)
	//
	g2 := &EnvelopeGoods{EnvelopeNo: g0.EnvelopeNo}
	ok, err := db.GetOne(g2)
	fmt.Println()
	fmt.Printf("%+v\n", g2)
	fmt.Printf("%+v\n", g1)
	require.NoError(t, err)
	assert.True(t, ok)
	assert.Equal(t, *g1.EnvelopeNo, *g2.EnvelopeNo)
	assert.Equal(t, g1.EnvelopeType, g2.EnvelopeType)
	assert.Equal(t, g1.Username, g2.Username)
	assert.Equal(t, g1.UserId, g2.UserId)
	assert.Equal(t, g1.Blessing.String, g2.Blessing.String)
	assert.Equal(t, g1.Amount.String(), g2.Amount.String())
	assert.Equal(t, g1.AmountOne.String(), g2.AmountOne.String())
	assert.Equal(t, g1.Quantity, g2.Quantity)
	assert.Equal(t, g1.RemainAmount.String(), g2.RemainAmount.String())
	assert.Equal(t, g1.RemainQuantity, g2.RemainQuantity)
	assert.Equal(t, g1.ExpiredAt.Second(), g2.ExpiredAt.Second())
	assert.Equal(t, g1.Status, g2.Status)
	assert.Equal(t, g1.OrderType, g2.OrderType)
	assert.Equal(t, g1.PayStatus, g2.PayStatus)
}

var Drop = "DROP TABLE IF EXISTS `red_envelope_goods`;"
var Create = "CREATE TABLE `red_envelope_goods`" +
	"(" +
	"  `id`              bigint(20)     NOT NULL AUTO_INCREMENT COMMENT '自增ID'," +
	"  `envelope_no`     varchar(32)    NOT NULL COMMENT '红包编号,红包唯一标识 '," +
	"  `envelope_type`   tinyint(d)     NOT NULL COMMENT '红包类型：普通红包，碰运气红包,过期红包'," +
	"  `username`        varchar(64)             DEFAULT NULL COMMENT '用户名称'," +
	"  `user_id`         varchar(40)    NOT NULL COMMENT '用户编号, 红包所属用户 '," +
	"  `blessing`        varchar(64)             DEFAULT NULL COMMENT '祝福语'," +
	"  `amount`          decimal(30, 6) NOT NULL DEFAULT '0.000000' COMMENT '红包总金额'," +
	"  `amount_one`      decimal(30, 6) NOT NULL DEFAULT '0.000000' COMMENT '单个红包金额，碰运气红包无效'," +
	"  `quantity`        int(10)        NOT NULL COMMENT '红包总数量 '," +
	"  `remain_amount`   decimal(30, 6) NOT NULL DEFAULT '0.000000' COMMENT '红包剩余金额额'," +
	"  `remain_quantity` int(10)        NOT NULL COMMENT '红包剩余数量 '," +
	"  `expired_at`      datetime(3)    NOT NULL COMMENT '过期时间'," +
	"  `status`          tinyint(d)     NOT NULL COMMENT '红包/订单状态：0 创建、1 发布启用、2过期、3失效'," +
	"  `order_type`      tinyint(d)     NOT NULL COMMENT '订单类型：发布单、退款单 '," +
	"  `pay_status`      tinyint(d)     NOT NULL COMMENT '支付状态：未支付，支付中，已支付，支付失败 '," +
	"  `created_at`      datetime(3)    NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间'," +
	"  `updated_at`      datetime(3)    NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间'," +
	"  PRIMARY KEY (`id`) USING BTREE," +
	"  UNIQUE KEY `envelope_no_idx` (`envelope_no`) USING BTREE," +
	"  KEY `id_user_idx` (`user_id`) USING BTREE" +
	") ENGINE = InnoDB" +
	"  AUTO_INCREMENT = 101" +
	"  DEFAULT CHARSET = utf8" +
	"  ROW_FORMAT = DYNAMIC;"
