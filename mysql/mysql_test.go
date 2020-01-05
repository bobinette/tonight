package mysql

import (
	"database/sql"
	"fmt"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/require"

	"github.com/bobinette/tonight/tonighttest"
)

func TestStores(t *testing.T) {
	db, err := sql.Open("mysql", fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local&timeout=1s",
		"root",
		"root",
		"127.0.0.1",
		"3306",
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
