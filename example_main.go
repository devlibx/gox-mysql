package main

import (
	"context"
	"fmt"
	"github.com/devlibx/gox-base"
	"github.com/devlibx/gox-base/serialization"
	"github.com/devlibx/gox-mysql/database"
	"github.com/devlibx/gox-mysql/tests/e2etest/sql/users"
	"go.uber.org/zap"
	"sync/atomic"
	"time"
)

var testMySQLConfig = &database.MySQLConfig{
	ServerName:  "test_server",
	Host:        "localhost",
	Port:        3306,
	User:        "test",
	Password:    "test",
	Db:          "users",
	TablePrefix: "integrating_tests",

	EnableSqlQueryLogging:       true,
	EnableSqlQueryMetricLogging: true,
	MetricDumpIntervalSec:       1,
	MetricResetAfterEveryNSec:   10,
}

// Main is a sample where we insert to MySQL user table. It also shows we dump the metric every 10 sec
// (enabled by MetricResetAfterEveryNSec)
func main() {
	sqlDb, _ := database.NewMySQLDb(NewCrossFunctionProvider(), testMySQLConfig)
	// sqlDb, _ := database.NewMySQLDbWithoutLogging(testMySQLConfig)
	q := users.New(sqlDb)

	sqlDb.RegisterPostCallbackFunc(func(data database.PostCallbackData) {
		fmt.Println("PostCallbackData=", serialization.StringifySuppressError(data, "na"))
	})

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

	go func() {
		time.Sleep(30 * time.Second)
		sqlDb.Close()
	}()

	<-sqlDb.WaitCloseChannel()
}

func NewCrossFunctionProvider() gox.CrossFunction {
	var loggerConfig zap.Config
	loggerConfig = zap.NewDevelopmentConfig()
	logger, _ := loggerConfig.Build()
	return gox.NewCrossFunction(logger)
}
