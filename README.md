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

##### Config file

```go
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
	ServerName: "test_server",
	Host:       "localhost",
	Port:       3306,
	User:       "test",
	Password:   "test",
	Db:         "users",

	EnableSqlQueryLogging:       true,
	EnableSqlQueryMetricLogging: true,
}

func main() {

	// Setup DB
	sqlDb, _ := pkg.NewMySQLDbWithoutLogging(testMySQLConfig)
	q := users.New(sqlDb)

	// Persist user
	result, err := q.PersistUser(context.Background(), "Harish")

}
```