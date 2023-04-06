package main

import (
	"context"
	"fmt"
	"github.com/devlibx/gox-base"
	"github.com/devlibx/gox-base/serialization"
	"github.com/devlibx/gox-mysql/database"
	"github.com/devlibx/gox-mysql/tests/e2etest/sql/users"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/opentracer"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
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

	// Datadog tracer enabled - you can use any other
	// How to run DD agent in your local??
	// docker run -it -p 8125:8125 -p 8126:8126 -v /var/run/docker.sock:/var/run/docker.sock:ro -v /proc/:/host/proc/:ro -v /sys/fs/cgroup/:/host/sys/fs/cgroup:ro -e DD_API_KEY=<YOUR_API_KEY>  gcr.io/datadoghq/agent:7
	if true {
		agentAddr := fmt.Sprintf("%s:%d", "localhost", 8126)
		t := opentracer.New(tracer.WithAgentAddr(agentAddr), tracer.WithServiceName("hbtest"), tracer.WithEnv("stage"))
		opentracing.SetGlobalTracer(t)
	}

	sqlDb, _ := database.NewMySQLDb(NewCrossFunctionProvider(), testMySQLConfig)
	// sqlDb, _ := database.NewMySQLDbWithoutLogging(testMySQLConfig)
	q := users.New(sqlDb)

	sqlDb.RegisterPostCallbackFunc(func(data database.PostCallbackData) {
		fmt.Println("PostCallbackData=", serialization.StringifySuppressError(data, "na"))

		if data.TimeTaken > 10 && data.Name == "users.(*Queries).PersistUser" {
			span, _ := opentracing.StartSpanFromContext(data.Ctx, data.GetDbCallNameForTracing())
			defer span.Finish()
			span.SetTag("error", true)
			span.SetTag("time_taken", data.TimeTaken)
		}
	})

	var count int32 = 0
	for i := 0; i < 10; i++ {
		go func(index int) {
			for {
				atomic.AddInt32(&count, 1)
				do(count, q)
			}
		}(i)
	}

	go func() {
		time.Sleep(1000 * time.Second)
		sqlDb.Close()
	}()

	<-sqlDb.WaitCloseChannel()
}

func do(count int32, q *users.Queries) {
	span, ctx := opentracing.StartSpanFromContext(context.Background(), "Api_call_To_Save")
	defer span.Finish()

	time.Sleep(10 * time.Millisecond)

	name := fmt.Sprintf("Harish_%d", count)
	persist(ctx, q, name)

	time.Sleep(20 * time.Millisecond)
	get(ctx, q, name)
}

func persist(ctx context.Context, q *users.Queries, name string) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Api_call_To_Persist")
	defer span.Finish()

	result, err := q.PersistUser(ctx, users.PersistUserParams{Name: "Harish", Department: "tech"})
	var id, rows int64
	if err == nil {
		id, _ = result.LastInsertId()
		rows, _ = result.RowsAffected()
	} else {
		fmt.Println(err)
	}
	_, _ = id, rows
	time.Sleep(1000 * time.Millisecond)
}

func get(ctx context.Context, q *users.Queries, name string) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Api_call_To_Get_And_Verify")
	defer span.Finish()

	result, err := q.GetUserByNameAndDepartment(ctx, users.GetUserByNameAndDepartmentParams{Name: "Harish", Department: "tech"})
	_, _ = result, err
	time.Sleep(1000 * time.Millisecond)
}

func NewCrossFunctionProvider() gox.CrossFunction {
	var loggerConfig zap.Config
	loggerConfig = zap.NewDevelopmentConfig()
	logger, _ := loggerConfig.Build()
	return gox.NewCrossFunction(logger)
}
