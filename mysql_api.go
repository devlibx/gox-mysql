package gox_mysql

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/harishb2k/gox-base"
	goxdb "github.com/harishb2k/gox-database"
	"time"
)

func init() {
	goxdb.RegisterDatabasePlugin("mysql", func(config *goxdb.Config, function gox.CrossFunction) (goxdb.Db, error) {
		return NewMySqlDb(config, function)
	})
}

type mysqlDb struct {
	mysqlInsertOp
	mysqlSelectOp
	mysqlExecuteOp
}

func NewMySqlDb(config *goxdb.Config, cf gox.CrossFunction) (goxdb.Db, error) {
	if config.Port == 0 {
		config.Port = 3306
	}
	if len(config.Url) == 0 {
		config.Url = []string{"127.0.0.1"}
	}
	if config.QueryTimout <= 0 {
		config.QueryTimout = 1000
	}
	url := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", config.User, config.Password, config.Url[0], config.Port, config.Db)
	if db, err := sql.Open("mysql", url); err != nil {
		return nil, &goxdb.DatabaseError{Op: goxdb.Open, Query: "", Args: nil, Err: err}
	} else {
		return &mysqlDb{
			mysqlInsertOp:  mysqlInsertOp{DB: db, CrossFunction: cf, config: config},
			mysqlSelectOp:  mysqlSelectOp{DB: db, CrossFunction: cf, config: config},
			mysqlExecuteOp: mysqlExecuteOp{DB: db, CrossFunction: cf, config: config},
		}, nil
	}
}

func withDeadline(config *goxdb.Config, cf gox.CrossFunction) (context.Context, context.CancelFunc) {
	return context.WithDeadline(context.Background(), cf.Now().Add(time.Duration(config.QueryTimout)*time.Millisecond))
}
