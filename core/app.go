// Package Core serves as the backbone of the application.
//
// It defines the main App interface and its base implementation.
package core

import (
	"context"
	"embed"
	"fmt"
	"htmx-go/apiserver"
	"htmx-go/database"
	"log"
	"os"

	"github.com/joho/godotenv"
)

var _ App = (*BaseApp)(nil)

type App interface {
	DB() database.Database
	API() apiserver.API
	FS() *embed.FS
	Context() context.Context
	IsBootstrapped() bool
	Bootstrap() error
	Start(ctx context.Context) error
	Shutdown(ctx context.Context) error
}

type BaseApp struct {
	db  database.Database
	fs  *embed.FS
	api apiserver.API
	ctx context.Context
}

func NewApp(fs *embed.FS) *BaseApp {
	return &BaseApp{
		fs: fs,
	}
}

func (app *BaseApp) DB() database.Database {
	return app.db
}

func (app *BaseApp) API() apiserver.API {
	return app.api
}

func (app *BaseApp) FS() *embed.FS {
	return app.fs
}

func (app *BaseApp) Context() context.Context {
	return app.ctx
}

func (app *BaseApp) IsBootstrapped() bool {
	if app.DB() == nil || app.API() == nil || app.FS() == nil || app.Context() == nil {
		return false
	}
	if !app.db.IsBootstrapped() || !app.api.IsBootstrapped() {
		return false
	}
	return true
}

func (app *BaseApp) Bootstrap() error {
	log.Println("🚨 App bootstrap initiated")

	app.ctx = context.Background()

	env := os.Getenv("ENV")
	if env != "production" {
		log.Println("🚧 running in development mode")
		if err := godotenv.Load(); err != nil {
			log.Println("failed to load .env file")
			return err
		}
	} else {
		log.Println("🚀 running in production mode")
	}

	db, err := database.New()
	if err != nil {
		log.Println("failed to connect to database")
		return err
	}

	app.db = db
	if err := app.db.Bootstrap(); err != nil {
		log.Println("failed to bootstrap database")
		return fmt.Errorf("bootstrap: %s", err)
	}

	api := apiserver.New(app.ctx, app.db, app.FS())
	app.api = api
	if err := app.api.Bootstrap(app.db); err != nil {
		log.Println("failed to bootstrap api")
		return fmt.Errorf("bootstrap: %s", err)
	}

	log.Println("🎉 App bootstrap successful")
	return nil
}

func (app *BaseApp) Start(ctx context.Context) error {
	if err := app.api.Serve(); err != nil {
		return fmt.Errorf("start: %s", err)
	}
	return nil
}

func (app *BaseApp) Shutdown(ctx context.Context) error {
	log.Println("🚨 App shutdown initiated")

	if err := app.db.Close(); err != nil {
		log.Println("failed to close database connection")
		return fmt.Errorf("shutdown: %s", err)
	}

	if err := app.api.Shutdown(ctx); err != nil {
		log.Println("failed to shutdown api")
		return fmt.Errorf("shutdown: %s", err)
	}
	return nil
}
