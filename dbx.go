package dbx

import (
    "database/sql"
    "fmt"
    "github.com/tietang/dbx/mapping"
    "time"
)

func Open(settings Settings) (db *Database, err error) {
    db = &Database{}
    db.DB, err = sql.Open(settings.DriverName, settings.DataSourceName())
    if err != nil {
        panic(err)
    }
    db.EntityMapper = mapping.NewEntityMapper()
    db.Runner = NewRunner(db.DB, db.EntityMapper)
    db.ILogger = &defaultLogger{}
    db.Runner.LoggerSettings = defaultLoggerSettings
    db.SetLogging(settings.LoggingEenabled)
    //设置最大打开连接数
    db.SetMaxOpenConns(settings.MaxOpenConns)
    //设置连接池最大空闲连接数
    db.SetMaxIdleConns(settings.MaxIdleConns)
    //设置连接最大生存时间
    db.SetConnMaxLifetime(settings.ConnMaxLifetime)
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
    DriverName      string
    User            string
    Password        string
    Database        string
    Host            string
    Options         map[string]string
    ConnMaxLifetime time.Duration
    MaxOpenConns    int
    MaxIdleConns    int
    LoggingEenabled bool
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
