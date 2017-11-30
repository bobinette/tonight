package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo"

	"github.com/bobinette/tonight"
	"github.com/bobinette/tonight/bleve"
	"github.com/bobinette/tonight/mysql"
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
	fmt.Println(mysqlAddr)
	taskRepo, err := mysql.NewTaskRepository(mysqlAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer taskRepo.Close()

	index := &bleve.Index{}
	if err := index.Open("./bleve/index"); err != nil {
		log.Fatal(err)
	}
	defer index.Close()

	uiService := tonight.NewUIService(taskRepo, index)

	// Create server + register routes
	srv := echo.New()

	echo.NotFoundHandler = func(c echo.Context) error {
		return fmt.Errorf("route %s (method %s) not found", c.Request().URL, c.Request().Method)
	}

	if err := tonight.RegisterTemplateRenderer(srv, "templates"); err != nil {
		log.Fatal(err)
	}

	// UI routes
	// -- Home
	srv.GET("/", func(c echo.Context) error { return c.Redirect(http.StatusPermanentRedirect, "/home") })
	srv.GET("/home", uiService.Home)

	// -- Calls serving html to partially update the page
	srv.GET("/ui/tasks", uiService.Search)
	srv.POST("/ui/tasks", uiService.CreateTask)
	srv.POST("/ui/tasks/:id", uiService.Update)
	srv.POST("/ui/tasks/:id/done", uiService.MarkDone)
	srv.DELETE("/ui/tasks/:id", uiService.Delete)
	srv.GET("/ui/done", uiService.DoneTasks)
	srv.POST("/ui/ranks", uiService.UpdateRanks)

	srv.POST("/ui/plan", uiService.Plan)
	srv.GET("/ui/plan", uiService.CurrentPlanning)
	srv.DELETE("/ui/plan", uiService.DismissPlanning)

	// Ping
	srv.GET("/api/ping", tonight.Ping)

	// API
	indexer := tonight.Indexer{
		Repository: taskRepo,
		Index:      index,
	}
	srv.POST("/api/reindex", indexer.IndexAll)

	// Assets
	srv.Static("/assets", "assets")
	srv.Static("/external", "external")
	srv.Static("/fonts", "fonts")

	if err := srv.Start(":9090"); err != nil {
		log.Fatal(err)
	}
}
