package e2etest

import (
	"context"
	"fmt"
	"github.com/devlibx/gox-base"
	"github.com/devlibx/gox-mysql/database"
	"github.com/devlibx/gox-mysql/tests/e2etest/sql/users"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"os"
	"testing"
)

var testMySQLConfig = &database.MySQLConfig{
	ServerName:  "test_server",
	Host:        "localhost",
	Port:        3306,
	User:        "test",
	Password:    "test",
	Db:          "users",
	TablePrefix: "integrating_tests",
}

func TestSimpleTestCase(t *testing.T) {
	if os.Getenv("E2E_TESTS") != "true" {
		t.Skip("Enable end-to-end test by setting E2E_TESTS=true")
	}

	sqlDb, err := database.NewMySQLDb(NewCrossFunctionProvider(), testMySQLConfig)
	assert.NoError(t, err)

	t.Run("Insert a new user", func(t *testing.T) {
		q := users.New(sqlDb)
		result, err := q.PersistUser(context.Background(), "Harish")
		assert.NoError(t, err)
		id, _ := result.LastInsertId()
		rows, _ := result.RowsAffected()
		fmt.Println("Id=", id, "Rows=", rows)

		user, err := q.GetUser(context.Background(), "Harish")
		assert.NoError(t, err)
		fmt.Println(user)

	})

}

func NewCrossFunctionProvider() gox.CrossFunction {
	var loggerConfig zap.Config
	loggerConfig = zap.NewDevelopmentConfig()
	logger, _ := loggerConfig.Build()
	return gox.NewCrossFunction(logger)
}
