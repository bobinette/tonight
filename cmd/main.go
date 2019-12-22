package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/BurntSushi/toml"
	_ "github.com/go-sql-driver/mysql"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"github.com/bobinette/tonight"
	"github.com/bobinette/tonight/mysql"
)

func main() {
	// Load configuration
	var cfg struct {
		Web struct {
			Bind string `toml:"bind"`
		} `toml:"web"`

		MySQL struct {
			User     string `toml:"user"`
			Password string `toml:"password"`
			Host     string `toml:"host"`
			Port     string `toml:"port"`
			Database string `toml:"database"`
		} `toml:"mysql"`

		FrontEnd struct {
			Mode     string `toml:"mode"`
			ProxyURL string `toml:"proxyUrl"`
			Dir      string `toml:"dir"`
		} `toml:"front-end"`
	}

	if _, err := toml.DecodeFile("config.toml", &cfg); err != nil {
		log.Fatal(err)
	}
	// Load configuration -- end

	// MySQL and stores
	db, err := sql.Open("mysql", fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local&timeout=1s",
		cfg.MySQL.User,
		cfg.MySQL.Password,
		cfg.MySQL.Host,
		cfg.MySQL.Port,
		cfg.MySQL.Database,
	))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	ctx := context.Background()
	if err := db.PingContext(ctx); err != nil {
		log.Fatal(err)
	}
	// MySQL and stores -- end

	// HTTP server via echo
	srv := echo.New()
	srv.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))

	srv.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))

	echo.NotFoundHandler = func(c echo.Context) error {
		return fmt.Errorf("route %s (method %s) not found", c.Request().URL, c.Request().Method)
	}

	srv.HTTPErrorHandler = echo.HTTPErrorHandler(func(err error, c echo.Context) {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
	})
	// HTTP server via echo -- env

	// Register and start
	eventStore := mysql.NewEventStore(db)
	taskStore := mysql.NewTaskStore(db)
	tonight.RegisterHTTP(
		srv.Group("/api"),
		eventStore,
		taskStore,
	)

	// @TODO: not prod ready. Use the config to determine what should be used
	if cfg.FrontEnd.Mode == "proxy" {
		proxyURL, err := url.Parse(cfg.FrontEnd.ProxyURL)
		if err != nil {
			log.Fatal(err)
		}
		targets := []*middleware.ProxyTarget{
			{
				URL: proxyURL,
			},
		}
		srv.Use(middleware.ProxyWithConfig(middleware.ProxyConfig{
			Skipper: func(c echo.Context) bool {
				return strings.HasPrefix(c.Path(), "/api")
			},
			Balancer: middleware.NewRandomBalancer(targets),
		}))
	} else {
		srv.Static("/", cfg.FrontEnd.Dir)
	}

	if err := srv.Start(cfg.Web.Bind); err != nil {
		log.Fatal(err)
	}
}
