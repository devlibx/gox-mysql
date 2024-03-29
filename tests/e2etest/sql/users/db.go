// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.2

package users

import (
	"context"
	"database/sql"
	"fmt"
)

type DBTX interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

func New(db DBTX) *Queries {
	return &Queries{db: db}
}

func Prepare(ctx context.Context, db DBTX) (*Queries, error) {
	q := Queries{db: db}
	var err error
	if q.getUserStmt, err = db.PrepareContext(ctx, GetUser); err != nil {
		return nil, fmt.Errorf("error preparing query GetUser: %w", err)
	}
	if q.getUserByNameAndDepartmentStmt, err = db.PrepareContext(ctx, GetUserByNameAndDepartment); err != nil {
		return nil, fmt.Errorf("error preparing query GetUserByNameAndDepartment: %w", err)
	}
	if q.getUsersStmt, err = db.PrepareContext(ctx, GetUsers); err != nil {
		return nil, fmt.Errorf("error preparing query GetUsers: %w", err)
	}
	if q.persistUserStmt, err = db.PrepareContext(ctx, PersistUser); err != nil {
		return nil, fmt.Errorf("error preparing query PersistUser: %w", err)
	}
	if q.updateUserNameStmt, err = db.PrepareContext(ctx, UpdateUserName); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateUserName: %w", err)
	}
	return &q, nil
}

func (q *Queries) Close() error {
	var err error
	if q.getUserStmt != nil {
		if cerr := q.getUserStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getUserStmt: %w", cerr)
		}
	}
	if q.getUserByNameAndDepartmentStmt != nil {
		if cerr := q.getUserByNameAndDepartmentStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getUserByNameAndDepartmentStmt: %w", cerr)
		}
	}
	if q.getUsersStmt != nil {
		if cerr := q.getUsersStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getUsersStmt: %w", cerr)
		}
	}
	if q.persistUserStmt != nil {
		if cerr := q.persistUserStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing persistUserStmt: %w", cerr)
		}
	}
	if q.updateUserNameStmt != nil {
		if cerr := q.updateUserNameStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateUserNameStmt: %w", cerr)
		}
	}
	return err
}

func (q *Queries) exec(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) (sql.Result, error) {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).ExecContext(ctx, args...)
	case stmt != nil:
		return stmt.ExecContext(ctx, args...)
	default:
		return q.db.ExecContext(ctx, query, args...)
	}
}

func (q *Queries) query(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) (*sql.Rows, error) {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).QueryContext(ctx, args...)
	case stmt != nil:
		return stmt.QueryContext(ctx, args...)
	default:
		return q.db.QueryContext(ctx, query, args...)
	}
}

func (q *Queries) queryRow(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) *sql.Row {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).QueryRowContext(ctx, args...)
	case stmt != nil:
		return stmt.QueryRowContext(ctx, args...)
	default:
		return q.db.QueryRowContext(ctx, query, args...)
	}
}

type Queries struct {
	db                             DBTX
	tx                             *sql.Tx
	getUserStmt                    *sql.Stmt
	getUserByNameAndDepartmentStmt *sql.Stmt
	getUsersStmt                   *sql.Stmt
	persistUserStmt                *sql.Stmt
	updateUserNameStmt             *sql.Stmt
}

func (q *Queries) WithTx(tx *sql.Tx) *Queries {
	return &Queries{
		db:                             tx,
		tx:                             tx,
		getUserStmt:                    q.getUserStmt,
		getUserByNameAndDepartmentStmt: q.getUserByNameAndDepartmentStmt,
		getUsersStmt:                   q.getUsersStmt,
		persistUserStmt:                q.persistUserStmt,
		updateUserNameStmt:             q.updateUserNameStmt,
	}
}
