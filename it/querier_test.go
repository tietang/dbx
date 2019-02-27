package it

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/segmentio/ksuid"
	"github.com/shopspring/decimal"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestGetAccount(t *testing.T) {

	Convey("写入测试", t, func() {
		Convey("正常", func() {
			a := &Account{
				Balance: decimal.NewFromFloat(100), Status: 1,
				UserId:    ksuid.New().Next().String(),
				Username:  sql.NullString{String: "测试用户1", Valid: true},
				AccountNo: ksuid.New().Next().String(),
			}
			rs, err := db.Insert(a)
			So(err, ShouldBeNil)
			fmt.Printf("AccountNo=%s\n", a.AccountNo)
			So(rs, ShouldNotBeNil)
			q := &Account{AccountNo: a.AccountNo}
			ok, err := db.GetOne(q)
			So(err, ShouldBeNil)
			So(ok, ShouldBeTrue)
			So(q.AccountNo, ShouldEqual, a.AccountNo)
			So(q.Balance.String(), ShouldEqual, a.Balance.String())
			So(q.UserId, ShouldEqual, a.UserId)
			So(q.CreatedAt, ShouldNotBeNil)
			So(q.UpdatedAt, ShouldNotBeNil)
		})
	})
}

func TestInsert(t *testing.T) {
	Convey("写入测试", t, func() {
		Convey("正常", func() {
			g0 := newEnvelopeGoods(1)
			rs, err := db.InsertContext(context.Background(), &g0)
			So(err, ShouldBeNil)
			So(rs, ShouldNotBeNil)
			id, err := rs.LastInsertId()
			So(err, ShouldBeNil)
			So(id > 0, ShouldBeTrue)
			rows, err := rs.RowsAffected()
			So(err, ShouldBeNil)
			So(rows, ShouldEqual, 1)
			g1 := &EnvelopeGoods{EnvelopeNo: g0.EnvelopeNo}

			ok, err := db.GetOne(g1)
			fmt.Println()
			fmt.Printf("%+v\n", g0)
			fmt.Printf("%+v\n", g1)
			So(err, ShouldBeNil)
			So(ok, ShouldBeTrue)
			So(*g1.EnvelopeNo, ShouldEqual, *g0.EnvelopeNo)
			So(g1.EnvelopeType, ShouldEqual, g0.EnvelopeType)
			So(g1.Username, ShouldEqual, g0.Username)
			So(g1.UserId, ShouldEqual, g0.UserId)
			So(g1.Blessing.String, ShouldEqual, g0.Blessing.String)
			So(g1.Amount.String(), ShouldEqual, g0.Amount.String())
			So(g1.AmountOne.String(), ShouldEqual, g0.AmountOne.String())
			So(g1.Quantity, ShouldEqual, g0.Quantity)
			So(g1.RemainAmount.String(), ShouldEqual, g0.RemainAmount.String())
			So(g1.RemainQuantity, ShouldEqual, g0.RemainQuantity)
			So(g1.ExpiredAt.Second(), ShouldEqual, g0.ExpiredAt.Second())
			So(g1.Status, ShouldEqual, g0.Status)
			So(g1.OrderType, ShouldEqual, g0.OrderType)
			So(g1.PayStatus, ShouldEqual, g0.PayStatus)

		})
	})
}

func TestUpdate(t *testing.T) {
	Convey("更新测试", t, func() {
		Convey("正常", func() {
			g0 := newEnvelopeGoods(1)
			rs, err := db.InsertContext(context.Background(), &g0)
			So(err, ShouldBeNil)
			So(rs, ShouldNotBeNil)
			id, err := rs.LastInsertId()
			So(err, ShouldBeNil)
			So(id > 0, ShouldBeTrue)
			rows, err := rs.RowsAffected()
			So(err, ShouldBeNil)
			So(rows, ShouldEqual, 1)
			g1 := newEnvelopeGoods(2)
			g1.EnvelopeNo = g0.EnvelopeNo
			//
			rs, err = db.Update(&g1)
			fmt.Println()
			fmt.Printf("%+v\n", g0)
			fmt.Printf("%+v\n", g1)
			So(err, ShouldBeNil)
			So(err, ShouldBeNil)
			So(rs, ShouldNotBeNil)
			id, err = rs.LastInsertId()
			So(err, ShouldBeNil)
			So(id == 0, ShouldBeTrue)
			rows, err = rs.RowsAffected()
			So(err, ShouldBeNil)
			So(rows, ShouldEqual, 1)
			//
			g2 := &EnvelopeGoods{EnvelopeNo: g0.EnvelopeNo}
			ok, err := db.GetOne(g2)
			fmt.Println()
			fmt.Printf("%+v\n", g2)
			fmt.Printf("%+v\n", g1)
			So(err, ShouldBeNil)
			So(ok, ShouldBeTrue)
			So(*g2.EnvelopeNo, ShouldEqual, *g1.EnvelopeNo)
			So(g2.EnvelopeType, ShouldEqual, g1.EnvelopeType)
			So(g2.Username, ShouldEqual, g1.Username)
			So(g2.UserId, ShouldEqual, g1.UserId)
			So(g2.Blessing.String, ShouldEqual, g1.Blessing.String)
			So(g2.Amount.String(), ShouldEqual, g1.Amount.String())
			So(g2.AmountOne.String(), ShouldEqual, g1.AmountOne.String())
			So(g2.Quantity, ShouldEqual, g1.Quantity)
			So(g2.RemainAmount.String(), ShouldEqual, g1.RemainAmount.String())
			So(g2.RemainQuantity, ShouldEqual, g1.RemainQuantity)
			So(g2.ExpiredAt.Second(), ShouldEqual, g1.ExpiredAt.Second())
			So(g2.Status, ShouldEqual, g1.Status)
			So(g2.OrderType, ShouldEqual, g1.OrderType)
			So(g2.PayStatus, ShouldEqual, g1.PayStatus)

		})
	})
}

var Drop = "DROP TABLE IF EXISTS `red_envelope_goods`;"
var Create = "CREATE TABLE `red_envelope_goods`" +
	"(" +
	"  `id`              bigint(20)     NOT NULL AUTO_INCREMENT COMMENT '自增ID'," +
	"  `envelope_no`     varchar(32)    NOT NULL COMMENT '红包编号,红包唯一标识 '," +
	"  `envelope_type`   tinyint(2)     NOT NULL COMMENT '红包类型：普通红包，碰运气红包,过期红包'," +
	"  `username`        varchar(64)             DEFAULT NULL COMMENT '用户名称'," +
	"  `user_id`         varchar(40)    NOT NULL COMMENT '用户编号, 红包所属用户 '," +
	"  `blessing`        varchar(64)             DEFAULT NULL COMMENT '祝福语'," +
	"  `amount`          decimal(30, 6) NOT NULL DEFAULT '0.000000' COMMENT '红包总金额'," +
	"  `amount_one`      decimal(30, 6) NOT NULL DEFAULT '0.000000' COMMENT '单个红包金额，碰运气红包无效'," +
	"  `quantity`        int(10)        NOT NULL COMMENT '红包总数量 '," +
	"  `remain_amount`   decimal(30, 6) NOT NULL DEFAULT '0.000000' COMMENT '红包剩余金额额'," +
	"  `remain_quantity` int(10)        NOT NULL COMMENT '红包剩余数量 '," +
	"  `expired_at`      datetime(3)    NOT NULL COMMENT '过期时间'," +
	"  `status`          tinyint(2)     NOT NULL COMMENT '红包/订单状态：0 创建、1 发布启用、2过期、3失效'," +
	"  `order_type`      tinyint(2)     NOT NULL COMMENT '订单类型：发布单、退款单 '," +
	"  `pay_status`      tinyint(2)     NOT NULL COMMENT '支付状态：未支付，支付中，已支付，支付失败 '," +
	"  `created_at`      datetime(3)    NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间'," +
	"  `updated_at`      datetime(3)    NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间'," +
	"  PRIMARY KEY (`id`) USING BTREE," +
	"  UNIQUE KEY `envelope_no_idx` (`envelope_no`) USING BTREE," +
	"  KEY `id_user_idx` (`user_id`) USING BTREE" +
	") ENGINE = InnoDB" +
	"  AUTO_INCREMENT = 101" +
	"  DEFAULT CHARSET = utf8" +
	"  ROW_FORMAT = DYNAMIC;"
