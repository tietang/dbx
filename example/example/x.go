package main

import (
	"fmt"
	"github.com/tietang/dbx/mapping"
)

func main() {
	m := mapping.NewEntityMapper()
	e, r := m.GetEntity(&Order{})
	fmt.Println(e, r)
}

type Order struct {
	Id      int64  `db:"id,omitempty"` // 自增id
	TradeNo string `db:"trade_no,uni"` // 交易单号
}
