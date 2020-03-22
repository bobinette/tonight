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

func TestStores(t *testing.T) {
	host := orString(os.Getenv("MYSQL_HOST"), "127.0.0.1")
	port := orString(os.Getenv("MYSQL_PORT"), "3306")
	user := orString(os.Getenv("MYSQL_USER"), "root")
	password := orString(os.Getenv("MYSQL_PASSWORD"), "root")
	db, err := sql.Open("mysql", fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local&timeout=1s",
		user,
		password,
		host,
		port,
		"tonight_v2_test",
	))
	require.NoError(t, err)
	defer func() {
		var err error
		_, err = db.Exec("DELETE FROM projects")
		require.NoError(t, err)
		_, err = db.Exec("DELETE FROM users")
		require.NoError(t, err)
		require.NoError(t, db.Close())
	}()

	projectStore := NewProjectStore(db)
	taskStore := NewTaskStore(db)
	userStore := NewUserStore(db)
	tonighttest.TestStores(t, projectStore, taskStore, userStore)
}
