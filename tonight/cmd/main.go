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

		JWT struct {
			Key string `toml:"key"`
		} `toml:"key"`

		Google struct {
			CookieSecret string `toml:"cookie_secret"`
			ClientID     string `toml:"client_id"`
			ClientSecret string `toml:"client_secret"`
			RedirectURL  string `toml:"redirect_url"`
		} `toml:"google"`

		App struct {
			Dir    string `toml:"dir"`
			Assets string `toml:"assets"`
		} `toml:"app"`
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

	index := &bleve.Index{}
	if err := index.Open(cfg.Bleve.Path); err != nil {
		log.Fatal(err)
	}
	defer index.Close()

	jwtKey := []byte(cfg.JWT.Key)

	// Create server + register routes
	srv := echo.New()

	echo.NotFoundHandler = func(c echo.Context) error {
		return fmt.Errorf("route %s (method %s) not found", c.Request().URL, c.Request().Method)
	}

	srv.HTTPErrorHandler = tonight.HTTPErrorHandler
	srv.Use(middleware.Logger())
	// srv.Use(middleware.Recover())

	// Login handler
	tonight.RegisterLoginHandler(
		srv,
		jwtKey,
		[]byte(cfg.Google.CookieSecret),
		cfg.Google.ClientID,
		cfg.Google.ClientSecret,
		cfg.Google.RedirectURL,
		userRepo,
	)

	// API handler
	tonight.RegisterAPIHandler(srv, jwtKey, taskRepo, index, planningRepo, userRepo)

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

	if err := srv.Start(":9090"); err != nil {
		log.Fatal(err)
	}
}
