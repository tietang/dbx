package main

import (
	"context"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/segmentio/ksuid"
	"github.com/shopspring/decimal"
	"github.com/tietang/dbx"
	"time"
)

var db *dbx.Database

func init() {
	//url := "root:123456@tcp(192.168.232.175:3306)/po0?charset=utf8&parseTime=true"
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
	var err error
	db, err = dbx.Open(settings)
	if err != nil {
		panic(err)
	}
	db.RegisterTable(&EnvelopeGoods{}, "red_envelope_goods")

}

func main() {
	g := insert(db)
	FindExampleContext(db, g.EnvelopeNo)
}

func FindExampleContext(db *dbx.Database, envelopeNo *string) {

	var egs []EnvelopeGoods
	err := db.Find(&egs, "select * from red_envelope_goods where id<=? ", 111)
	fmt.Println(err)
	for key, value := range egs {
		fmt.Printf("%+v %+v \n", key, value)
	}
}
func FindExampleContext2(db *dbx.Database, envelopeNo *string) {
	q := EnvelopeGoods{
		EnvelopeNo: envelopeNo,
	}
	var egs []EnvelopeGoods
	err := db.FindExampleContext(context.Background(), q, &egs)
	fmt.Println(err)
	for key, value := range egs {
		fmt.Printf("%+v %+v \n", key, value)
	}

}

func insert(db *dbx.Database) *EnvelopeGoods {
	g := newEnvelopeGoods()
	rs, err := db.InsertContext(context.Background(), &g)
	id, e := rs.LastInsertId()
	fmt.Printf("id=%#v %+v ", id, e)
	r, e := rs.RowsAffected()
	fmt.Printf("rows=%#v %+v ", r, e)
	fmt.Println(err)
	g.Username = "1231"
	g.Blessing = "1312312"
	rs, err = db.UpdateContext(context.Background(), &g)
	id, e = rs.LastInsertId()
	fmt.Printf("id=%#v %+v ", id, e)
	r, e = rs.RowsAffected()
	fmt.Printf("rows=%#v %+v ", r, e)
	fmt.Println(err)
	return &g
}

func newEnvelopeGoods() EnvelopeGoods {
	g := EnvelopeGoods{

		EnvelopeType: 1,
		Username:     "test-name",
		UserId:       ksuid.New().Next().String(),
		Blessing:     "test-blessing",
		Amount:       decimal.NewFromFloat(100.01),
		Quantity:     10,
		Status:       1,
		OrderType:    1,
		PayStatus:    1,
	}
	amountOne := decimal.NewFromFloat(10.01)
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
	Blessing     string           `db:"blessing"`             //祝福语
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
