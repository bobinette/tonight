package mysql

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	// "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bobinette/tonight/tonight/tests"
)

func TestTagReader(t *testing.T) {
	mysqlAddr := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		"root",
		"root",
		"192.168.50.4",
		"3306",
		"tonight_test",
	)

	if os.Getenv("TRAVIS") == "true" {
		mysqlAddr = fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
			"tonight",
			"tonight",
			"127.0.0.1",
			"3307",
			"tonight_test",
		)
	}

	db, err := sql.Open("mysql", mysqlAddr)
	require.NoError(t, err)

	_, err = db.Exec("DELETE FROM tasks")
	_, err = db.Exec("DELETE FROM users")
	require.NoError(t, err)
	// defer func() {
	// 	_, err := db.Exec("DELETE FROM tasks")
	// 	_, err = db.Exec("DELETE FROM users")
	// 	assert.NoError(t, err)
	// 	assert.NoError(t, db.Close())
	// }()

	taskRepo := NewTaskRepository(db)
	userRepo := NewUserRepository(db)
	tagReader := NewTagReader(db)
	tests.TestTagReader(t, tagReader, taskRepo, userRepo)
}
