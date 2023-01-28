package it

import (
	"database/sql"
	"github.com/segmentio/ksuid"
	"github.com/shopspring/decimal"
	"strconv"
	"time"
)

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
	EnvelopeNo   *string          `db:"envelope_no,uni"` //红包编号,红包唯一标识
	EnvelopeType int              //`db:"envelope_type"`        //红包类型：普通红包，碰运气红包
	Username     string           `db:"username"` //用户名称
	UserId       string           //`db:"user_id"`              //用户编号, 红包所属用户
	Blessing     sql.NullString   `db:"blessing"` //祝福语
	Amount       decimal.Decimal  `db:"amount"`   //红包总金额
	AmountOne    *decimal.Decimal //`db:"amount_one"`           //单个红包金额，碰运气红包无效
	Quantity     int              `db:"quantity"` //红包总数量
	ExpiredAt    *time.Time       //`db:"expired_at"`           //过期时间
	Status       int              `db:"status"` //红包状态：0红包初始化，1启用，2失效
	OrderType    int              //`db:"order_type"`           //订单类型：发布单、退款单
	PayStatus    int              //`db:"pay_status"`           //支付状态：未支付，支付中，已支付，支付失败
	CreatedAt    *time.Time       `db:"created_at,omitempty"` //创建时间
	UpdatedAt    time.Time        `db:"updated_at,omitempty"` //更新时间
}

//type EnvelopeGoods struct {
//	Inventory
//	EnvelopeNo   *string          `db:"envelope_no,uni"`      //红包编号,红包唯一标识
//	EnvelopeType int              `db:"envelope_type"`        //红包类型：普通红包，碰运气红包
//	Username     string           `db:"username"`             //用户名称
//	UserId       string           `db:"user_id"`              //用户编号, 红包所属用户
//	Blessing     sql.NullString   `db:"blessing"`             //祝福语
//	Amount       decimal.Decimal  `db:"amount"`               //红包总金额
//	AmountOne    *decimal.Decimal `db:"amount_one"`           //单个红包金额，碰运气红包无效
//	Quantity     int              `db:"quantity"`             //红包总数量
//	ExpiredAt    *time.Time       `db:"expired_at"`           //过期时间
//	Status       int              `db:"status"`               //红包状态：0红包初始化，1启用，2失效
//	OrderType    int              `db:"order_type"`           //订单类型：发布单、退款单
//	PayStatus    int              `db:"pay_status"`           //支付状态：未支付，支付中，已支付，支付失败
//	CreatedAt    *time.Time       `db:"created_at,omitempty"` //创建时间
//	UpdatedAt    time.Time        `db:"updated_at,omitempty"` //更新时间
//}

type Inventory struct {
	EnvelopeNo     string          `db:"envelope_no"`     //红包编号,红包唯一标识
	RemainAmount   decimal.Decimal `db:"remain_amount"`   //红包剩余金额额
	RemainQuantity int             `db:"remain_quantity"` //红包剩余数量
}

type Account struct {
	Id           int64           `db:"id,omitempty"`         //账户ID
	AccountNo    string          `db:"account_no,uni"`       //账户编号,账户唯一标识
	AccountName  string          `db:"account_name"`         //账户名称,用来说明账户的简短描述,账户对应的名称或者命名，比如xxx积分、xxx零钱
	AccountType  int             `db:"account_type"`         //账户类型，用来区分不同类型的账户：积分账户、会员卡账户、钱包账户、红包账户
	CurrencyCode string          `db:"currency_code"`        //货币类型编码：CNY人民币，EUR欧元，USD美元 。。。
	UserId       string          `db:"user_id"`              //用户编号, 账户所属用户
	Username     sql.NullString  `db:"username"`             //用户名称
	Balance      decimal.Decimal `db:"balance"`              //账户可用余额
	Status       int             `db:"status"`               //账户状态，账户状态：0账户初始化，1启用，2停用
	CreatedAt    time.Time       `db:"created_at,omitempty"` //创建时间
	UpdatedAt    time.Time       `db:"updated_at,omitempty"` //更新时间
}
