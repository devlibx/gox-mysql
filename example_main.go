package main

import (
	"context"
	"fmt"
	"github.com/devlibx/gox-base"
	"github.com/devlibx/gox-mysql/pkg"
	"github.com/devlibx/gox-mysql/tests/e2etest/sql/users"
	"go.uber.org/zap"
	"sync/atomic"
	"time"
)

var testMySQLConfig = &pkg.MySQLConfig{
	ServerName:  "test_server",
	Host:        "localhost",
	Port:        3306,
	User:        "test",
	Password:    "test",
	Db:          "users",
	TablePrefix: "integrating_tests",

	EnableSqlQueryLogging:       true,
	EnableSqlQueryMetricLogging: true,
}

func main() {
	sqlDb, _ := pkg.NewMySQLDb(NewCrossFunctionProvider(), testMySQLConfig)
	q := users.New(sqlDb)
	var count int32 = 0
	for i := 0; i < 10; i++ {
		go func(index int) {
			for {
				atomic.AddInt32(&count, 1)
				result, err := q.PersistUser(context.Background(), "Harish")
				var id, rows int64
				if err == nil {
					id, _ = result.LastInsertId()
					rows, _ = result.RowsAffected()
				} else {
					fmt.Println(err)
				}
				if count%100 == 0 {
					// fmt.Println("Index=", index, "Id=", id, "Rows=", rows, "Count=", count)
				}
				_, _ = id, rows
				time.Sleep(1000 * time.Millisecond)
			}
		}(i)
	}
	time.Sleep(10 * time.Minute)
}

func NewCrossFunctionProvider() gox.CrossFunction {
	var loggerConfig zap.Config
	loggerConfig = zap.NewDevelopmentConfig()
	logger, _ := loggerConfig.Build()
	return gox.NewCrossFunction(logger)
}
