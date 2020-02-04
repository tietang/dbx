package it

import (
	"github.com/tietang/dbx"
	"time"
)

var db *dbx.Database

func init() {
	settings := dbx.Settings{
		DriverName: "mysql",
		User:       "root",
		Password:   "111111",
		//Password:   "123456",
		Host: "127.0.0.1:3306",
		//Host:            "172.16.1.248:3306",
		Database:        "po",
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
