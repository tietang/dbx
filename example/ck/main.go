package main

import (
	"fmt"
	_ "github.com/mailru/go-clickhouse/v2"
	"github.com/tietang/dbx"
	"time"
)

type DubboSlowConsumer struct {
	Service   string  `db:"service"`
	Method    string  `db:"method"`
	ResTimeMs float64 `db:"responseTime_MS"`
}

func main() {
	//conn, err := sql.Open("chhttp",
	//"http://admin:N+JeozDXx5b3JxQrLQmTsLbXqo8N0YjS@10.99.74.78::81234/mo?read_timeout=10s&write_timeout=20s")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//if err := conn.Ping(); err != nil {
	//	log.Fatal(err)
	//}

	settings := dbx.Settings{
		DriverName:      "chhttp",
		Protocol:        "http",
		User:            "admin",
		Password:        "N+JeozDXx5b3JxQrLQmTsLbXqo8N0YjS",
		Host:            "10.99.74.78:8124",
		Database:        "mo",
		MaxOpenConns:    10,
		MaxIdleConns:    2,
		ConnMaxLifetime: time.Minute * 30,
		Options: map[string]string{
			"read_timeout":  "10s",
			"write_timeout": "20s",
		},
	}
	db, err := dbx.Open(settings)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	fmt.Println(err)
	sql := "SELECT   service,   method,   sum(elapsed) / 1000000 as responseTime_MS " +
		" FROM m_dubbo " +
		" WHERE dt >= toDateTime(?) AND dt <= toDateTime(?) and side = 'consumer' " +
		" GROUP BY   service,   method " +
		" HAVING responseTime_MS > 1000"

	var consumers []*DubboSlowConsumer
	end := time.Now().Unix()
	start := end - 5000
	err = db.Find(&consumers, sql, start, end)
	fmt.Println(err)
	for _, consumer := range consumers {
		fmt.Println(consumer)
	}

}
