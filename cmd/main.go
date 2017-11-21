package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo"

	"github.com/bobinette/tonight"
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

	uiService := tonight.NewUIService(taskRepo)

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
	srv.POST("/ui/tasks", uiService.CreateTask)
	srv.POST("/ui/tasks/:id", uiService.MarkDone)
	srv.DELETE("/ui/tasks/:id", uiService.Delete)
	srv.GET("/ui/done", uiService.DoneTasks)
	srv.POST("/ui/ranks", uiService.UpdateRanks)

	srv.POST("/ui/plan", uiService.Plan)

	// Ping
	srv.GET("/api/ping", tonight.Ping)

	// Assets
	srv.Static("/assets", "assets")
	srv.Static("/fonts", "fonts")

	if err := srv.Start(":9090"); err != nil {
		log.Fatal(err)
	}
}
