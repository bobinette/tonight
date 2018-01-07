package mysql

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bobinette/tonight/tonight/tests"
)

func TestTaskRepository(t *testing.T) {
	mysqlAddr := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		"root",
		"root",
		"192.168.50.4",
		"3306",
		"tonight_test",
	)
	db, err := sql.Open("mysql", mysqlAddr)
	require.NoError(t, err)

	_, err = db.Exec("DELETE FROM tasks")
	require.NoError(t, err)
	defer func() {
		_, err := db.Exec("DELETE FROM tasks")
		assert.NoError(t, err)
		assert.NoError(t, db.Close())
	}()

	taskRepo := NewTaskRepository(db)
	tests.TestTaskRepository(t, taskRepo)
}
