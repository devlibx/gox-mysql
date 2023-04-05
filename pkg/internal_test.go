package pkg

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCleanQuery(t *testing.T) {
	query := `-- name: PersistUser :execresult
INSERT INTO integrating_tests_users (name)
VALUES (?)
`
	out := cleanQuery(query)
	assert.Equal(t, "INSERT INTO integrating_tests_users (name) VALUES (?)", out)
}
