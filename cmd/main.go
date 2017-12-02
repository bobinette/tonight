package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

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
	db, err := sql.Open("mysql", mysqlAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	taskRepo := mysql.NewTaskRepository(db)
	userRepo := mysql.NewUserRepository(db)

	index := &bleve.Index{}
	if err := index.Open("./bleve/index"); err != nil {
		log.Fatal(err)
	}
	defer index.Close()

	jwtKey := []byte("tonight_secret")

	loginService := tonight.LoginService{
		Key:        jwtKey,
		Repository: userRepo,
	}
	uiService := tonight.NewUIService(taskRepo, index, userRepo)

	// Create server + register routes
	srv := echo.New()

	echo.NotFoundHandler = func(c echo.Context) error {
		return fmt.Errorf("route %s (method %s) not found", c.Request().URL, c.Request().Method)
	}

	if err := tonight.RegisterTemplateRenderer(srv, "templates"); err != nil {
		log.Fatal(err)
	}

	srv.HTTPErrorHandler = tonight.HTTPErrorHandler
	srv.Use(middleware.Logger())

	// UI routes
	// -- Home
	srv.GET("/", func(c echo.Context) error { return c.Redirect(http.StatusPermanentRedirect, "/ui/home") })
	srv.GET("/home", func(c echo.Context) error { return c.Redirect(http.StatusPermanentRedirect, "/ui/home") })

	srv.GET("/login", loginService.LoginPage)
	srv.POST("/login", loginService.Login)
	srv.POST("/logout", loginService.Logout)

	uiGroup := srv.Group("/ui")
	uiGroup.Use(tonight.JWTMiddleware(jwtKey))

	uiGroup.GET("/home", uiService.Home)

	// -- Calls serving html to partially update the page
	uiGroup.GET("/tasks", uiService.Search)
	uiGroup.POST("/tasks", uiService.CreateTask)
	uiGroup.POST("/tasks/:id", uiService.Update)
	uiGroup.POST("/tasks/:id/done", uiService.MarkDone)
	uiGroup.DELETE("/tasks/:id", uiService.Delete)
	uiGroup.GET("/done", uiService.DoneTasks)
	uiGroup.POST("/ranks", uiService.UpdateRanks)

	uiGroup.POST("/plan", uiService.Plan)
	uiGroup.GET("/plan", uiService.CurrentPlanning)
	uiGroup.DELETE("/plan", uiService.DismissPlanning)

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
