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
	queryInterface := users.New(sqlDb)

	// Persist user
	if result, err := queryInterface.PersistUser(context.Background(), "Harish"); err == nil {
		fmt.Println("User saved", result)
	} else {
		fmt.Println("Something is wrong", err)
	}

}
```