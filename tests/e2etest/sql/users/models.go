// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.2

package users

import (
	"database/sql"
)

type IntegratingTestsUser struct {
	ID      int32         `json:"id"`
	Name    string        `json:"name"`
	Deleted sql.NullInt32 `json:"deleted"`
}