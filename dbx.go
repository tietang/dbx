package dbx

import (
	"database/sql"
	"fmt"
	"imooc.com/tietang/red-envelope/infra/dbx/mapping"
	"time"
)

//db, err = dbx.Open("mysql", url)
//	if err != nil {
//		panic(err)
//	}
//	db.RegisterTable(&model{}, "[table name]")
func Open(driverName, dataSourceName string) (db *Database, err error) {
	db = &Database{}
	sqlDB, err := sql.Open(driverName, dataSourceName)
	db.DB = sqlDB
	db.EntityMapper = mapping.NewEntityMapper()
	db.Runner = NewRunner(sqlDB, db.EntityMapper)
	db.ILogger = &defaultLogger{}
	db.Runner.LoggerSettings = defaultLoggerSettings
	return db, err
}

type Database struct {
	*Runner
	*sql.DB
	mapping.EntityMapper
}

func (r *Database) Tx(fn func(run *TxRunner) error) error {
	tx, err := r.DB.Begin()

	if err != nil {
		return err
	}
	runner := NewTxRunner(tx)
	runner.EntityMapper = r.EntityMapper
	runner.LoggerSettings = r.LoggerSettings
	runner.ILogger = r.Logger()
	if err := fn(runner); err != nil {
		e := tx.Rollback()
		if e != nil {
			return err
		}
		return err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}

	return err
}

type Settings struct {
	User            string
	Password        string
	Database        string
	Host            string
	Options         map[string]string
	ConnMaxLifetime time.Duration
	MaxOpenConns    int
	MaxIdleConns    int
	LogginEenabled  bool
}

func (s *Settings) DataSourceName() string {
	queryString := ""
	for key, value := range s.Options {
		queryString += key + "=" + value + "&"
	}
	ustr := fmt.Sprintf("%s:%s@tcp(%s)/%s?%s", s.User, s.Password, s.Host, s.Database, queryString)

	return ustr
}
func (s *Settings) ShortDataSourceName() string {
	queryString := ""
	for key, value := range s.Options {
		queryString += key + "=" + value + "&"
	}
	ustr := fmt.Sprintf("%s:***@tcp(%s)/%s?%s", s.User, s.Host, s.Database, queryString)

	return ustr
}
