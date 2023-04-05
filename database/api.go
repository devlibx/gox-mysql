package database

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/devlibx/gox-base"
	"github.com/devlibx/gox-base/errors"
	"github.com/devlibx/gox-base/util"
	"github.com/opentracing/opentracing-go"
	_ "github.com/rcrowley/go-metrics"
	"go.uber.org/zap"
	"sync"
	"time"
)

const methodDepthToGetFunctionName = 3

type DB struct {
	db     *sql.DB
	dsn    string
	logger *zap.Logger
	config *MySQLConfig

	stopOnce             *sync.Once
	stopCompleteChan     chan bool
	stoppedCtx           context.Context
	stoppedCtxCancelFunc context.CancelFunc
}

func NewMySQLDbWithoutLogging(config *MySQLConfig) (*DB, error) {
	return NewMySQLDb(gox.NewNoOpCrossFunction(), config)
}

func NewMySQLDb(cf gox.CrossFunction, config *MySQLConfig) (*DB, error) {
	config.SetupDefaults()
	_db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", config.User, config.Password, config.Host, config.Port, config.Db))
	if err != nil {
		return nil, errors.Wrap(err, "failed to open SQL masterDb")
	}

	// To stop everything from DB side
	stoppedCtx, stoppedCtxCancelFunc := context.WithCancel(context.Background())

	// Start metric logging if enabled
	if config.EnableSqlQueryMetricLogging {
		startMetricDump(stoppedCtx, config)
	}

	return &DB{
		db:                   _db,
		dsn:                  "",
		logger:               cf.Logger().Named("db"),
		config:               config,
		stopOnce:             &sync.Once{},
		stopCompleteChan:     make(chan bool, 2),
		stoppedCtx:           stoppedCtx,
		stoppedCtxCancelFunc: stoppedCtxCancelFunc,
	}, nil
}

func (d *DB) WaitCloseChannel() chan bool {
	return d.stopCompleteChan
}

func (d *DB) Close() {
	d.stopOnce.Do(func() {
		d.stoppedCtxCancelFunc()
		time.Sleep(1 * time.Second)
		d.stopCompleteChan <- true
		close(d.stopCompleteChan)
	})
}

func (d *DB) Exec(query string, args ...interface{}) (sql.Result, error) {

	// Log query details and metrics
	defer newLogInf(query, d.logger, d.config.EnableSqlQueryLogging, d.config.EnableSqlQueryMetricLogging).done(args...)

	return d.db.Exec(query, args...)
}

func (d *DB) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {

	// Log span and trace
	span, ctx := opentracing.StartSpanFromContext(ctx, util.GetMethodName(methodDepthToGetFunctionName))
	defer span.Finish()

	// Log query details and metrics
	defer newLogInf(query, d.logger, d.config.EnableSqlQueryLogging, d.config.EnableSqlQueryMetricLogging).done(args...)

	return d.db.ExecContext(ctx, query, args...)
}

func (d *DB) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {

	// Log span and trace
	span, ctx := opentracing.StartSpanFromContext(ctx, util.GetMethodName(methodDepthToGetFunctionName))
	defer span.Finish()

	// Log query details and metrics
	defer newLogInf(query, d.logger, d.config.EnableSqlQueryLogging, d.config.EnableSqlQueryMetricLogging).done()

	return d.db.PrepareContext(ctx, query)
}

func (d *DB) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {

	// Log span and trace
	span, ctx := opentracing.StartSpanFromContext(ctx, util.GetMethodName(methodDepthToGetFunctionName))
	defer span.Finish()

	// Log query details and metrics
	defer newLogInf(query, d.logger, d.config.EnableSqlQueryLogging, d.config.EnableSqlQueryMetricLogging).done(args...)

	return d.db.QueryContext(ctx, query, args...)
}

func (d *DB) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {

	// Log span and trace
	span, ctx := opentracing.StartSpanFromContext(ctx, util.GetMethodName(methodDepthToGetFunctionName))
	defer span.Finish()

	// Log query details and metrics
	defer newLogInf(query, d.logger, d.config.EnableSqlQueryLogging, d.config.EnableSqlQueryMetricLogging).done(args...)

	return d.db.QueryRowContext(ctx, query, args...)
}
