package main

import (
	"fmt"
	_ "github.com/glebarez/go-sqlite"
	"github.com/tietang/dbx"
)

var db *dbx.Database

func init() {
	settings := dbx.Settings{
		DriverName:     "sqlite",
		DataSourceName: "some.db?_pragma=busy_timeout(5000)&_pragma=journal_mode(WAL)",
		LoggingEnabled: true,
	}
	sqlDb, err := dbx.Open(settings)
	if err != nil {
		panic(err)
	}
	db = sqlDb
	db.SetLogging(true)
	//db.RegisterTable(&EnvelopeGoods{}, "red_envelope_goods")
}

func main() {
	// connect
	//db, err := sql.Open("sqlite", ":memory:")
	//
	//db, err := sql.Open("sqlite", "./path/to/some.db")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//_ = db.QueryRow("select sqlite_version()")

	// get SQLite version
	//var v string
	//err := db.GetValue(&v, "select sqlite_version()")
	//v, err := db.GetString("select sqlite_version()")
	v, err := db.GetString(" select datetime(CURRENT_TIMESTAMP,'localtime')")
	fmt.Println(err)
	fmt.Println(v)

}
