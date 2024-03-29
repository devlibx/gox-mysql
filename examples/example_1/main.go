package main

import (
	"context"
	"fmt"
	"github.com/devlibx/gox-base/serialization"
	"github.com/devlibx/gox-mysql/database"
	"github.com/devlibx/gox-mysql/tests/e2etest/sql/users"
	"github.com/opentracing/opentracing-go"
)

var testMySQLConfig = &database.MySQLConfig{
	ServerName:                  "test_server",
	Host:                        "localhost",
	Port:                        3306,
	User:                        "test",
	Password:                    "test",
	Db:                          "users",
	EnableSqlQueryLogging:       false,
	EnableSqlQueryMetricLogging: true,
}

func main() {
	// Setup DB
	sqlDb, err := database.NewMySQLDbWithoutLogging(testMySQLConfig)
	if err != nil {
		panic(err)
	}

	// Start: ====================== This is Optional and added to show how can you add trace to slow query ============
	// This is a callback (Optional)
	// It tell you time taken, when this DB call started, ended etc.
	// You can use it to alert if some specific query take some time (you get the name of the query in the payload)
	sqlDb.RegisterPostCallbackFunc(func(data database.PostCallbackData) {
		fmt.Println("PostCallbackData=", serialization.StringifySuppressError(data, "na"))

		// We will get the callback which contains total time taken for debugging
		if data.TimeTaken > 1 {
			span, _ := opentracing.StartSpanFromContext(data.Ctx, data.GetDbCallNameForTracing())
			defer span.Finish()
			span.SetTag("error", true)
			span.SetTag("time_taken", data.TimeTaken)
			fmt.Printf("Something is wrong it took very long: data=%s \n", serialization.StringifySuppressError(data, "na"))

			// >> You will see following if time > 1ms
			// Something is wrong it took very long: data={"name":"users.(*Queries).PersistUser","start_time":1680761853659,"end_time":1680761853672,"time_taken":13,"error":null}
		}
	})
	// End: ====================== This is Optional and added to show how can you add trace to slow query ============

	queryInterface := users.New(sqlDb)

	// Persist user
	if result, err := queryInterface.PersistUser(
		context.Background(), users.PersistUserParams{Name: "Harish", Department: "tech"},
	); err == nil {
		fmt.Println("User saved", result)
	} else {
		fmt.Println("Something is wrong", err)
	}

	if users, err := queryInterface.GetUserByNameAndDepartment(
		context.Background(),
		users.GetUserByNameAndDepartmentParams{Name: "Harish", Department: "tech"},
	); err == nil {
		for _, u := range users {
			fmt.Println("Users: ID=", u.ID, "Name=", u.Name, "Department=", u.Department)
		}
	} else {
		fmt.Println("Something is wrong", err)
	}
}
