package it

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/segmentio/ksuid"
	"github.com/shopspring/decimal"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/tietang/dbx"
	"strconv"
	"testing"
	"time"
)

var db *dbx.Database

func init() {
	settings := dbx.Settings{
		DriverName: "mysql",
		User:       "root",
		Password:   "123456",
		Host:       "192.168.232.175:3306",
		//Host:            "172.16.1.248:3306",
		Database:        "po0",
		MaxOpenConns:    10,
		MaxIdleConns:    2,
		ConnMaxLifetime: time.Minute * 30,
		LoggingEnabled:  true,
		Options: map[string]string{
			"charset":   "utf8",
			"parseTime": "true",
		},
	}
	sqlDb, err := dbx.Open(settings)
	if err != nil {
		panic(err)
	}
	db = sqlDb
	db.SetLogging(true)
	db.RegisterTable(&EnvelopeGoods{}, "red_envelope_goods")
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

			err = db.GetOne(g1)
			fmt.Println()
			fmt.Printf("%+v\n", g0)
			fmt.Printf("%+v\n", g1)
			So(err, ShouldBeNil)
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
			err = db.GetOne(g2)
			fmt.Println()
			fmt.Printf("%+v\n", g2)
			fmt.Printf("%+v\n", g1)
			So(err, ShouldBeNil)
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

func newEnvelopeGoods(seed int) EnvelopeGoods {
	deedStr := strconv.Itoa(seed)
	g := EnvelopeGoods{
		EnvelopeType: 1,
		Username:     "test-name" + deedStr,
		UserId:       ksuid.New().Next().String(),
		Blessing:     sql.NullString{String: "test-blessing" + deedStr, Valid: true},
		Amount:       decimal.NewFromFloat(100.01 + float64(seed)),
		Quantity:     10 + seed,
		Status:       1 + seed,
		OrderType:    1 + seed,
		PayStatus:    1 + seed,
	}
	amountOne := decimal.NewFromFloat(10.01 + float64(seed))
	g.AmountOne = &amountOne
	EnvelopeNo := ksuid.New().Next().String()
	g.EnvelopeNo = &EnvelopeNo
	t := time.Now().Add(time.Hour * 24)
	g.ExpiredAt = &t

	g.RemainAmount = g.Amount
	g.RemainQuantity = g.Quantity
	return g
}

type EnvelopeGoods struct {
	Inventory
	EnvelopeNo   *string          `db:"envelope_no,uni"`      //红包编号,红包唯一标识
	EnvelopeType int              `db:"envelope_type"`        //红包类型：普通红包，碰运气红包
	Username     string           `db:"username"`             //用户名称
	UserId       string           `db:"user_id"`              //用户编号, 红包所属用户
	Blessing     sql.NullString   `db:"blessing"`             //祝福语
	Amount       decimal.Decimal  `db:"amount"`               //红包总金额
	AmountOne    *decimal.Decimal `db:"amount_one"`           //单个红包金额，碰运气红包无效
	Quantity     int              `db:"quantity"`             //红包总数量
	ExpiredAt    *time.Time       `db:"expired_at"`           //过期时间
	Status       int              `db:"status"`               //红包状态：0红包初始化，1启用，2失效
	OrderType    int              `db:"order_type"`           //订单类型：发布单、退款单
	PayStatus    int              `db:"pay_status"`           //支付状态：未支付，支付中，已支付，支付失败
	CreatedAt    *time.Time       `db:"created_at,omitempty"` //创建时间
	UpdatedAt    time.Time        `db:"updated_at,omitempty"` //更新时间
}

type Inventory struct {
	EnvelopeNo     string          `db:"envelope_no"`     //红包编号,红包唯一标识
	RemainAmount   decimal.Decimal `db:"remain_amount"`   //红包剩余金额额
	RemainQuantity int             `db:"remain_quantity"` //红包剩余数量
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
