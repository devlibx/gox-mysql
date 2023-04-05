package database

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/devlibx/gox-base"
	"github.com/devlibx/gox-base/errors"
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

	callbacks *Callbacks
}

type Callbacks struct {
	PostCallbackFunc PostCallbackFunc
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

	// Add callback functions
	callbacks := &Callbacks{
		PostCallbackFunc: func(data PostCallbackData) {},
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
		callbacks:            callbacks,
	}, nil
}

func (d *DB) RegisterPostCallbackFunc(function PostCallbackFunc) {
	d.callbacks.PostCallbackFunc = function
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

func (d *DB) Exec(query string, args ...interface{}) (result sql.Result, err error) {
	defer d.buildNewLogInf(context.Background(), query).done(err, args...)
	result, err = d.db.Exec(query, args...)
	return result, err
}

func (d *DB) ExecContext(ctx context.Context, query string, args ...interface{}) (result sql.Result, err error) {
	defer d.buildNewLogInf(ctx, query).done(err, args...)
	result, err = d.db.ExecContext(ctx, query, args...)
	return result, err
}

func (d *DB) PrepareContext(ctx context.Context, query string) (result *sql.Stmt, err error) {
	defer d.buildNewLogInf(ctx, query).done(err)
	result, err = d.db.PrepareContext(ctx, query)
	return result, err
}

func (d *DB) QueryContext(ctx context.Context, query string, args ...interface{}) (result *sql.Rows, err error) {
	defer d.buildNewLogInf(ctx, query).done(err, args...)
	result, err = d.db.QueryContext(ctx, query, args...)
	return result, err
}

func (d *DB) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	defer d.buildNewLogInf(ctx, query).done(nil, args...)
	return d.db.QueryRowContext(ctx, query, args...)
}

func (d *DB) buildNewLogInf(ctx context.Context, query string) logInfo {
	return newLogInf(ctx, query, d.logger, d.config.EnableSqlQueryLogging, d.config.EnableSqlQueryMetricLogging, d.callbacks)
}
