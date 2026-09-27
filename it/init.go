package it

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/tietang/dbx"
)

var db *dbx.Database

// 测试库/表自动创建：先连无库名 DSN 建库，再执行建表 SQL
const (
	itUser     = "root"
	itPassword = "NewRootPass@2026"
	itHost     = "192.168.2.215:3306"
	itDatabase = "po"
)

// createDatabaseIfNotExists 连接不含库名的 DSN，执行 CREATE DATABASE IF NOT EXISTS
func createDatabaseIfNotExists() error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/?charset=utf8&parseTime=true", itUser, itPassword, itHost)
	raw, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	defer raw.Close()
	_, err = raw.Exec("CREATE DATABASE IF NOT EXISTS `" + itDatabase + "` DEFAULT CHARSET utf8")
	return err
}

func init() {
	if err := createDatabaseIfNotExists(); err != nil {
		panic(fmt.Sprintf("auto create database %s failed: %v", itDatabase, err))
	}
	settings := dbx.Settings{
		DriverName:      "mysql",
		User:            itUser,
		Password:        itPassword,
		Host:            itHost,
		Database:        itDatabase,
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

	if err := createTables(); err != nil {
		panic(fmt.Sprintf("auto create tables failed: %v", err))
	}
}
