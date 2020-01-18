package main

import (
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"github.com/dveselov/juno/internal/config"
	"github.com/dveselov/juno/internal/routes"
)

func main() {
	db, err := sqlx.Connect("postgres", os.Getenv("APP_DB_URL"))
	if err != nil {
		log.Fatalln(err)
	}
	appConfig := config.GetAppConfig(db)

	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	m, err := migrate.NewWithDatabaseInstance(
		"file://./migrations/",
		"postgres", driver,
	)
	if err != nil {
		log.Fatalln(err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalln(err)
	}

	e := echo.New()

	// Middleware
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			context := &config.AppContext{appConfig, c}
			return next(context)
		}
	})
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/begin", routes.Begin)
	e.POST("/verify", routes.Verify)

	e.Logger.Fatal(e.Start(appConfig.ListenAddr))
}
