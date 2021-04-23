package gox_mysql

import (
	"context"
	"database/sql"
	"github.com/harishb2k/gox-base"
	"github.com/harishb2k/gox-base/metrics"
	goxdb "github.com/harishb2k/gox-database"
	"time"
)

type mysqlSelectOp struct {
	*sql.DB
	gox.CrossFunction
	config *goxdb.Config
}

func (m *mysqlSelectOp) Find(metric metrics.LabeledMetric, query string, args ...interface{}) (results []map[string]interface{}, err error) {
	ctx, cancel := withDeadline(m.config, m.CrossFunction)
	defer cancel()
	return m.FindContext(ctx, metric, query, args...)
}

func (m *mysqlSelectOp) FindOne(metric metrics.LabeledMetric, query string, args ...interface{}) (result map[string]interface{}, err error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Duration(m.config.QueryTimout)*time.Millisecond)
	defer cancelFunc()
	return m.FindOneContext(ctx, metric, query, args...)
}

func (m *mysqlSelectOp) FindContext(ctx context.Context, metric metrics.LabeledMetric, query string, args ...interface{}) (results []map[string]interface{}, err error) {
	m.Counter(metric.Name).Inc()

	// Prepare statement
	var statement *sql.Stmt
	statement, err = m.PrepareContext(ctx, query)
	if err != nil {
		m.Counter(metric.NameWithErrorPrefix()).Inc()
		return nil, &goxdb.DatabaseError{Op: goxdb.FindAll, Query: query, Args: args, Err: err}
	}
	defer statement.Close()

	// Query rows
	var rows *sql.Rows
	rows, err = statement.QueryContext(ctx, args...)
	if err != nil {
		m.Counter(metric.NameWithErrorPrefix()).Inc()
		return nil, &goxdb.DatabaseError{Op: goxdb.FindAll, Query: query, Args: args, Err: err}
	}
	defer rows.Close()

	// Make columns to pull data
	cols, _ := rows.Columns()
	columns := make([]interface{}, len(cols))
	columnPointers := make([]interface{}, len(cols))
	for i := 0; i < len(cols); i++ {
		columnPointers[i] = &columns[i]
	}

	// Scan rows and return result
	var resultsToReturn []map[string]interface{}
	for rows.Next() {

		// Scan records for row
		if err = rows.Scan(columnPointers...); err != nil {
			m.Counter(metric.NameWithErrorPrefix()).Inc()
			return nil, &goxdb.DatabaseError{Op: goxdb.FindAll, Query: query, Args: args, Err: err}
		}

		// Build result
		result := make(map[string]interface{})
		for i, colName := range cols {
			val := columnPointers[i].(*interface{})
			result[colName] = *val
		}
		resultsToReturn = append(resultsToReturn, result)
	}

	if len(resultsToReturn) > 0 {
		m.Counter(metric.NameWithSuccessPrefix()).Inc()
		return resultsToReturn, nil
	} else {
		m.Counter(metric.NameWithErrorPrefix() + "_no_record").Inc()
		return nil, &goxdb.NoDatabaseRecordError{DatabaseError: goxdb.DatabaseError{Op: goxdb.FindAll, Query: query, Args: args, Err: nil}}
	}
}

func (m *mysqlSelectOp) FindOneContext(ctx context.Context, metric metrics.LabeledMetric, query string, args ...interface{}) (result map[string]interface{}, err error) {
	m.Counter(metric.Name).Inc()

	// Prepare statement
	var statement *sql.Stmt
	statement, err = m.PrepareContext(ctx, query)
	if err != nil {
		m.Counter(metric.NameWithErrorPrefix()).Inc()
		return nil, &goxdb.DatabaseError{Op: goxdb.Find, Query: query, Args: args, Err: err}
	}
	defer statement.Close()

	// Query rows
	var rows *sql.Rows
	rows, err = statement.QueryContext(ctx, args...)
	if err != nil {
		m.Counter(metric.NameWithErrorPrefix()).Inc()
		return nil, &goxdb.DatabaseError{Op: goxdb.Find, Query: query, Args: args, Err: err}
	}
	defer rows.Close()

	// Make sure we have results
	if !rows.Next() {
		m.Counter(metric.NameWithErrorPrefix() + "_no_record").Inc()
		return nil, &goxdb.NoDatabaseRecordError{DatabaseError: goxdb.DatabaseError{Op: goxdb.Find, Query: query, Args: args, Err: nil}}
	}

	// Make columns to pull data
	cols, _ := rows.Columns()
	columns := make([]interface{}, len(cols))
	columnPointers := make([]interface{}, len(cols))
	for i := 0; i < len(cols); i++ {
		columnPointers[i] = &columns[i]
	}

	// Scan rows and return result
	if err = rows.Scan(columnPointers...); err != nil {
		m.Counter(metric.NameWithErrorPrefix()).Inc()
		return nil, &goxdb.DatabaseError{Op: goxdb.Find, Query: query, Args: args, Err: err}
	}

	// Build result
	result = make(map[string]interface{})
	for i, colName := range cols {
		val := columnPointers[i].(*interface{})
		result[colName] = *val
	}

	m.Counter(metric.NameWithSuccessPrefix()).Inc()
	return result, nil
}
