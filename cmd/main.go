package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo"

	"github.com/bobinette/tonight"
	"github.com/bobinette/tonight/bleve"
	"github.com/bobinette/tonight/mysql"
)

func main() {
	mysqlAddr := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		"root",         //username
		"root",         //password
		"192.168.50.4", //host
		"3306",         //port
		"tonight",      //database
	)
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
