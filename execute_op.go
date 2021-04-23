package gox_mysql

import (
	"context"
	"database/sql"
	"github.com/harishb2k/gox-base"
	"github.com/harishb2k/gox-base/metrics"
	goxdb "github.com/harishb2k/gox-database"
)

type mysqlExecuteOp struct {
	*sql.DB
	gox.CrossFunction
	config *goxdb.Config
}

func (m *mysqlExecuteOp) Execute(metric metrics.LabeledMetric, query string, args ...interface{}) (int, error) {
	ctx, cancel := withDeadline(m.config, m.CrossFunction)
	defer cancel()
	return m.ExecuteContext(ctx, metric, query, args...)
}

func (m *mysqlExecuteOp) ExecuteContext(ctx context.Context, metric metrics.LabeledMetric, query string, args ...interface{}) (int, error) {
	m.Counter(metric.Name).Inc()

	statement, err := m.PrepareContext(ctx, query)
	if err != nil {
		m.Counter(metric.NameWithErrorPrefix()).Inc()
		return 0, &goxdb.DatabaseError{Op: goxdb.Execute, Query: query, Args: args, Err: err}
	}
	defer statement.Close()

	r, err := statement.Exec(args...)
	if err != nil {
		m.Counter(metric.NameWithErrorPrefix()).Inc()
		return 0, &goxdb.DatabaseError{Op: goxdb.Execute, Query: query, Args: args, Err: err}
	}

	result, err := r.RowsAffected()
	if err != nil {
		m.Counter(metric.NameWithErrorPrefix()).Inc()
		return int(result), &goxdb.DatabaseError{Op: goxdb.Execute, Query: query, Args: args, Err: err}
	}

	m.Counter(metric.NameWithSuccessPrefix()).Inc()
	return int(result), nil
}
