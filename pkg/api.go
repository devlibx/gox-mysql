package pkg

import (
	"database/sql"
	"fmt"
	"github.com/devlibx/gox-base/errors"
)

func NewMySQLDb(config *MySQLConfig) (*sql.DB, error) {
	config.SetupDefaults()
	_db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", config.User, config.Password, config.Host, config.Port, config.Db))
	if err != nil {
		return nil, errors.Wrap(err, "failed to open SQL masterDb")
	}
	return _db, nil
}
