package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/BurntSushi/toml"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"github.com/bobinette/tonight/tonight"
	"github.com/bobinette/tonight/tonight/bleve"
	"github.com/bobinette/tonight/tonight/mysql"
)

func main() {
	var env string
	flag.StringVar(&env, "env", "dev", "")
	flag.Parse()

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

		Bleve struct {
			Path string `toml:"path"`
		} `toml:"bleve"`

		App struct {
			Dir    string `toml:"dir"`
			Assets string `toml:"assets"`
		} `toml:"app"`

		Dev struct {
			ForwardedEmail string `toml:"forwarded_email"`
		} `toml:"dev"`
	}

	if _, err := toml.DecodeFile(fmt.Sprintf("config.%s.toml", env), &cfg); err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("mysql", mysql.Format(
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

	taskRepo := mysql.NewTaskRepository(db)
	userRepo := mysql.NewUserRepository(db)
	planningRepo := mysql.NewPlanningRepository(db, taskRepo)
	tagReader := mysql.NewTagReader(db)

	index := &bleve.Index{}
	if err := index.Open(cfg.Bleve.Path); err != nil {
		log.Fatal(err)
	}
	defer index.Close()

	// Create server + register routes
	srv := echo.New()

	echo.NotFoundHandler = func(c echo.Context) error {
		return fmt.Errorf("route %s (method %s) not found", c.Request().URL, c.Request().Method)
	}

	srv.HTTPErrorHandler = tonight.HTTPErrorHandler
	srv.Use(middleware.Logger())

	if cfg.Dev.ForwardedEmail != "" {
		srv.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				if c.Request().Header.Get("X-Forwarded-Email") == "" {
					c.Request().Header.Set("X-Forwarded-Email", cfg.Dev.ForwardedEmail)
				}
				return next(c)
			}
		})
	}

	// API handler
	tonight.RegisterAPIHandler(srv, taskRepo, index, planningRepo, userRepo, tagReader)

	// Ping
	srv.GET("/api/ping", tonight.Ping)

	// Assets
	srv.Static("/", cfg.App.Dir)
	srv.Static("/static", cfg.App.Assets)

	clis := tonight.RegisterCLI(taskRepo, index)

	if len(os.Args) > 1 {
		if fn, ok := clis[os.Args[len(os.Args)-1]]; ok {
			fmt.Printf("Running cli: %s\n", os.Args[1])
			fn()
			return
		}
	}

	if err := srv.Start(cfg.Web.Bind); err != nil {
		log.Fatal(err)
	}
}
