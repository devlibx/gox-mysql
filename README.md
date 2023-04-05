### About

This library gives convenient access to MySQL.

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
)

var testMySQLConfig = &database.MySQLConfig{
	ServerName:                  "test_server",
	Host:                        "localhost",
	Port:                        3306,
	User:                        "test",
	Password:                    "test",
	Db:                          "users",
	EnableSqlQueryLogging:       true,
	EnableSqlQueryMetricLogging: true,
}

func main() {
	// Setup DB
	sqlDb, err := database.NewMySQLDbWithoutLogging(testMySQLConfig)
	if err != nil {
		panic(err)
	}
	
	// This is a callback (Optional)
	// It tell you time takne, when this DB call started, ended etc
	// You can use it to alert if some specific query take some time (you get the name of the query in the payload)
	sqlDb.RegisterPostCallbackFunc(func(data database.PostCallbackData) {
		fmt.Println("PostCallbackData=", serialization.StringifySuppressError(data, "na"))

		// We will get the callback which contains total time taken for debuting
		// Note you can add your own alerts e.g. If this is "PersistUser" and take more than 20 ms
		// then do something
		// >> PostCallbackData= {"name":"users.(*Queries).PersistUser","start_time":1680709127885,"end_time":1680709127898,"time_taken":13}
	})

	queryInterface := users.New(sqlDb)

	// Persist user
	if result, err := queryInterface.PersistUser(context.Background(), "Harish"); err == nil {
		fmt.Println("User saved", result)
	} else {
		fmt.Println("Something is wrong", err)
	}
}
```