package dbx

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/tietang/dbx/mapping"
	"time"
)

// settings := dbx.Settings{
//        DriverName: "mysql",
//        User:       "root",
//        Password:   "123456",
//        Host:       "192.168.232.175:3306",
//        //Host:            "172.16.1.248:3306",
//        Database:        "po0",
//        MaxOpenConns:    10,
//        MaxIdleConns:    2,
//        ConnMaxLifetime: time.Minute * 30,
//        LoggingEnabled:  true,
//        Options: map[string]string{
//            "charset":   "utf8",
//            "parseTime": "true",
//        },
//    }
//    var err error
//    db, err = dbx.Open(settings)
//    if err != nil {
//        panic(err)
//    }
var _ mapperExecutor = new(Database)
var _ sqlExecutor = new(Database)
var _ mapping.EntityMapper = new(Database)
var _ LoggerSettings = new(Database)

func Open(settings Settings) (db *Database, err error) {
	db = &Database{}
	db.DB, err = sql.Open(settings.DriverName, settings.DataSourceName())
	if err != nil {
		panic(err)
	}
	db.EntityMapper = mapping.NewEntityMapper()

	db.ILogger = &defaultLogger{}
	db.LoggerSettings = defaultLoggerSettings
	db.SetLogging(settings.LoggingEnabled)
	//
	runner := NewRunner(db.DB, db.EntityMapper)
	db.mapperExecutor = runner
	runner.LoggerSettings = db.LoggerSettings
	runner.ILogger = db.ILogger
	runner.EntityMapper = db.EntityMapper

	//设置最大打开连接数
	db.SetMaxOpenConns(settings.MaxOpenConns)
	//设置连接池最大空闲连接数
	db.SetMaxIdleConns(settings.MaxIdleConns)
	//设置连接最大生存时间
	db.SetConnMaxLifetime(settings.ConnMaxLifetime)
	db.DefaultAutoCommit = true
	if v, ok := settings.Options["autocommit"]; ok {
		if v == "false" {
			db.DefaultAutoCommit = false
		}
	}
	db.ping()
	return db, err
}

type Database struct {
	*sql.DB
	mapperExecutor
	mapping.EntityMapper
	ILogger
	LoggerSettings
	DefaultAutoCommit bool
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
		return err
	}
	return err
}

func (r *Database) rollback(tx *sql.Tx) error {
	if !r.DefaultAutoCommit {
		tx.Exec("SET AUTOCOMMIT=1")
	}
	e := tx.Rollback()
	return e
}

func (r *Database) commit(tx *sql.Tx) error {
	if !r.DefaultAutoCommit {
		tx.Exec("SET AUTOCOMMIT=1")
	}
	err := tx.Commit()
	return err
}

func (r *Database) ping() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	status := "up"
	if err := r.PingContext(ctx); err != nil {
		status = "down"
	}
	r.Log(&QueryStatus{
		Message: &status,
	})
}

type Settings struct {
	DriverName      string
	User            string
	Password        string
	Database        string
	Host            string
	Options         map[string]string
	ConnMaxLifetime time.Duration
	MaxOpenConns    int
	MaxIdleConns    int
	LoggingEnabled  bool
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
