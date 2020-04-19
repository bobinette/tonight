package mysql

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/require"

	"github.com/bobinette/tonight/tonighttest"
)

func orString(s, def string) string {
	if s != "" {
		return s
	}
	return def
}

func setUp(t *testing.T) (*sql.DB, func()) {
	host := orString(os.Getenv("MYSQL_HOST"), "127.0.0.1")
	port := orString(os.Getenv("MYSQL_PORT"), "3307")
	user := orString(os.Getenv("MYSQL_USER"), "root")
	password := orString(os.Getenv("MYSQL_PASSWORD"), "root")
	database := orString(os.Getenv("MYSQL_DB"), "tonight_v2_test")

	db, err := sql.Open("mysql", fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local&timeout=1s",
		user,
		password,
		host,
		port,
		database,
	))
	require.NoError(t, err)

	return db, func() {
		require.NoError(t, db.Close())
	}
}

func TestStores(t *testing.T) {
	db, tearDown := setUp(t)
	defer tearDown()
	defer func() {
		db.Exec("DELETE FROM projects")
		db.Exec("DELETE FROM users")
	}()

	projectStore := NewProjectStore(db)
	taskStore := NewTaskStore(db)
	tonighttest.TestStores(t, projectStore, taskStore)
}
