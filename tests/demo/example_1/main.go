package main

import (
	"context"
	"fmt"
	"github.com/devlibx/gox-mysql/pkg"
	"github.com/devlibx/gox-mysql/tests/e2etest/sql/users"
)

var testMySQLConfig = &pkg.MySQLConfig{
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
	sqlDb, err := pkg.NewMySQLDbWithoutLogging(testMySQLConfig)
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
