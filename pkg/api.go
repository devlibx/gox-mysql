package pkg

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/devlibx/gox-base"
	"github.com/devlibx/gox-base/errors"
	"go.uber.org/zap"
	"regexp"
	"strings"
	"time"
)

var regexToCleanQueryToDump = regexp.MustCompile(`^[^\n]+\n`)

type DB struct {
	db     *sql.DB
	dsn    string
	logger *zap.Logger
}

func NewMySQLDb(cf gox.CrossFunction, config *MySQLConfig) (*DB, error) {
	config.SetupDefaults()
	_db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", config.User, config.Password, config.Host, config.Port, config.Db))
	if err != nil {
		return nil, errors.Wrap(err, "failed to open SQL masterDb")
	}
	return &DB{
		db:     _db,
		dsn:    "",
		logger: cf.Logger().Named("db"),
	}, nil
}

func (d DB) Exec(query string, args ...interface{}) (sql.Result, error) {
	defer newLogInf("Exec", query, d.logger).done(args...)
	return d.db.Exec(query, args...)
}

func (d DB) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	defer newLogInf("ExecContext", query, d.logger).done(args...)
	return d.db.ExecContext(ctx, query, args...)
}

func (d DB) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	defer newLogInf("PrepareContext", query, d.logger).done()
	return d.db.PrepareContext(ctx, query)
}

func (d DB) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	defer newLogInf("QueryContext", query, d.logger).done()
	return d.db.QueryContext(ctx, query, args...)
}

func (d DB) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	defer newLogInf("QueryRowContext", query, d.logger).done()
	return d.db.QueryRowContext(ctx, query, args...)
}

type logInfo struct {
	name      string
	startTime int64
	query     string
	logger    *zap.Logger
}

func (l logInfo) dump(args ...interface{}) {
	query := regexToCleanQueryToDump.ReplaceAllString(l.query, "")
	query = strings.ReplaceAll(query, "\n", " ")
	l.logger.Info(l.name, zap.Int64("time", time.Now().UnixMilli()-l.startTime), zap.String("query", query), zap.Any("args", args))
}

func (l logInfo) done(args ...interface{}) {
	l.dump(args...)
}

func newLogInf(name string, query string, logger *zap.Logger) logInfo {
	return logInfo{startTime: time.Now().UnixMilli(), name: name, query: query, logger: logger}
}
