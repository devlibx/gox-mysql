package gox_mysql

import (
	"context"
	"database/sql"
	"github.com/harishb2k/gox-base"
	"github.com/harishb2k/gox-base/metrics"
	goxdb "github.com/harishb2k/gox-database"
)

type mysqlInsertOp struct {
	*sql.DB
	gox.CrossFunction
	config *goxdb.Config
}

func (m *mysqlInsertOp) Persist(metric metrics.LabeledMetric, query string, args ...interface{}) (interface{}, error) {
	ctx, cancel := withDeadline(m.config, m.CrossFunction)
	defer cancel()
	return m.PersistContext(ctx, metric, query, args...)
}

func (m *mysqlInsertOp) PersistContext(ctx context.Context, metric metrics.LabeledMetric, query string, args ...interface{}) (interface{}, error) {
	m.Counter(metric.Name).Inc()

	statement, err := m.PrepareContext(ctx, query)
	if err != nil {
		m.Counter(metric.NameWithErrorPrefix()).Inc()
		return nil, &goxdb.DatabaseError{Op: goxdb.Insert, Query: query, Args: args, Err: err}
	}
	defer statement.Close()

	_, err = statement.Exec(args...)
	if err != nil {
		m.Counter(metric.NameWithErrorPrefix()).Inc()
		return nil, &goxdb.DatabaseError{Op: goxdb.Insert, Query: query, Args: args, Err: err}
	}

	m.Counter(metric.NameWithSuccessPrefix()).Inc()
	return nil, nil
}
