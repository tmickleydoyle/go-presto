package presto

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestQuery(t *testing.T) {
	presto := os.Getenv("PRESTO")
	if presto == "" {
		t.Skip("Skipping tests because PRESTO isn't set in the environment")
	}

	q, err := NewQuery(presto, "test", "go-presto-test", "sys", "", "select count(*) from sys.node")
	require.NoError(t, err)

	rows := make([][]interface{}, 0, 1)
	for {
		row, err := q.Next()
		require.NoError(t, err)
		if row == nil {
			break
		}

		rows = append(rows, row)
	}

	assert.Equal(t, 1, len(rows))
	assert.Equal(t, 1, len(rows[0]))
}
