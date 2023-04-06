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