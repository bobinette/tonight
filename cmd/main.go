package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"github.com/bobinette/tonight/tonight"
	"github.com/bobinette/tonight/tonight/bleve"
	"github.com/bobinette/tonight/tonight/mysql"
)

func main() {
	mysqlConfig := struct {
		User     string
		Password string
		Host     string
		Port     string
		Database string
	}{
		User:     "root",
		Password: "root",
		Host:     "192.168.50.4",
		Port:     "3306",
		Database: "tonight",
	}

	if user := os.Getenv("MYSQL_USER"); user != "" {
		mysqlConfig.User = user
	}

	if password := os.Getenv("MYSQL_PASSWORD"); password != "" {
		mysqlConfig.Password = password
	}

	if host := os.Getenv("MYSQL_HOST"); host != "" {
		mysqlConfig.Host = host
	}

	if port := os.Getenv("MYSQL_PORT"); port != "" {
		mysqlConfig.Port = port
	}

	if database := os.Getenv("MYSQL_DATABASE"); database != "" {
		mysqlConfig.Database = database
	}

	mysqlAddr := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		mysqlConfig.User,
		mysqlConfig.Password,
		mysqlConfig.Host,
		mysqlConfig.Port,
		mysqlConfig.Database,
	)
	db, err := sql.Open("mysql", mysqlAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	taskRepo := mysql.NewTaskRepository(db)
	userRepo := mysql.NewUserRepository(db)
	planningRepo := mysql.NewPlanningRepository(db, taskRepo)

	index := &bleve.Index{}
	if err := index.Open("./tonight/bleve/index"); err != nil {
		log.Fatal(err)
	}
	defer index.Close()

	jwtKey := []byte("tonight_secret")

	// Create server + register routes
	srv := echo.New()

	echo.NotFoundHandler = func(c echo.Context) error {
		return fmt.Errorf("route %s (method %s) not found", c.Request().URL, c.Request().Method)
	}

	srv.HTTPErrorHandler = tonight.HTTPErrorHandler
	srv.Use(middleware.Logger())
	srv.Use(middleware.Recover())

	// Login handler
	tonight.RegisterLoginHandler(srv, jwtKey, userRepo)

	// API handler
	tonight.RegisterAPIHandler(srv, jwtKey, taskRepo, index, planningRepo, userRepo)

	// Ping
	srv.GET("/api/ping", tonight.Ping)

	// API
	indexer := tonight.Indexer{
		Repository: taskRepo,
		Index:      index,
	}
	srv.POST("/api/reindex", indexer.IndexAll)

	// Assets
	srv.Static("/css", "front/static/css")
	srv.Static("/fonts", "front/static/fonts")
	srv.Static("/images", "front/static/images")
	srv.Static("/js", "front/static/js")

	srv.Static("/", "app/dist")
	srv.Static("/static", "app/dist/static")

	if err := srv.Start(":9090"); err != nil {
		log.Fatal(err)
	}
}
