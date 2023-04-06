### About

This library gives convenient access to MySQL. It comes with few defaults:

1. SqlC support which allows toy to generate the Go typesafe code from SQL statements
2. OpenTracing - if OpenTracing is configured then it will add the time taken by DB call in the Span
3. After each DB call, it will call your PostCallbackFunc function -> if registered by you,
   You can log slow query, Or you can add it as an error in your span (given in example below)
4. SQL statement and SQL payload can be turned on to debug. Set EnableSqlQueryLogging=true
5. You can enable histogram during perf to how each query is doing. It is dumped to the console every N sec
   This will help in perf testing

##### Try out sample application

```
# DB name is users for this setup

CREATE TABLE IF NOT EXISTS `integrating_tests_users`
(
    `id`      int(11) unsigned NOT NULL AUTO_INCREMENT,
    `name`    varchar(265) NOT NULL DEFAULT '',
    `deleted` int          DEFAULT 0,
    PRIMARY KEY (`id`)
);
```

Run ```example_main.go```. Update the user/password in the file if needed.

---

### Example code

You can run it using ```go run examples/example_1/main.go```

```go
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

	queryInterface := users.New(sqlDb)

	// Persist user
	if result, err := queryInterface.PersistUser(context.Background(), users.PersistUserParams{Name: "Harish", Department: "tech"}); err == nil {
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
```

## Observability

All the following are part of the lib, you need to switch on/off as you need

```shell
 EnableSqlQueryLogging         => To print "Query, Time Taken, Params"
 EnableSqlQueryMetricLogging   => To print "Histogram"
 OpenTracking is automitaclly enabled, and it sends traces to whatever OpenTracing framework you have configured
```

#### Query, Time Taken, Params

This lib outputs logs which shows how much a query took, with SQL and the parameter used. It also provides a hook,
which is called at the end of each DB call which also exposes all this info.
You can use it to log slow queries and generate alerts

```shell
2023-04-06T12:00:38.232+0530    INFO    db      database/internal.go:45 users.(*Queries).PersistUser    
{"time": 128, "query": "INSERT INTO integrating_tests_users (name, department) VALUES (?, ?)", "args": ["Harish","tech"]}
```

#### Histogram

This lib has inbuilt support to dump histogram of each query (If enabled). This histogram is printed for each query
and gives details.

```shell
metrics: 11:55:47.745129 histogram INSERT INTO integrating_tests_users (name, department) VALUES (?, ?)
metrics: 11:55:47.745141   count:              32
metrics: 11:55:47.745145   min:                 1
metrics: 11:55:47.745149   max:               163
metrics: 11:55:47.745153   mean:               25.50
metrics: 11:55:47.745157   stddev:             34.65
metrics: 11:55:47.745160   median:             15.00
metrics: 11:55:47.745164   75%:                29.75
metrics: 11:55:47.745171   95%:               131.80
metrics: 11:55:47.745175   99%:               163.00
metrics: 11:55:47.745181   99.9%:             163.00
```

#### OpenTracing integration

Sample trace to show an API call and the automatically added ```DB_Call_GetUser``` by the lib

![alt text](https://github.com/devlibx/images/blob/master/Full_API_Call_with_Db_Traces.png?raw=true)

Sample trace to show an API call and to mark your custom trace ```Slow_Query_Trace__PersistUser``` when timeTake > Nms.

Sample code which you can use. You can do any other actions if needed.
```go
sqlDb.RegisterPostCallbackFunc(func (data database.PostCallbackData) {
if data.TimeTaken > 1 {
    span, _ := opentracing.StartSpanFromContext(data.Ctx, data.GetDbCallNameForTracing())
    defer span.Finish()
    span.SetTag("error", true)
    span.SetTag("time_taken", data.TimeTaken)
    }
})
```

![alt text](https://github.com/devlibx/images/blob/master/DB_Traces_with_Error.png)



