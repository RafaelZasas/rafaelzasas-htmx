package apiserver

import (
	"context"
	"embed"
	"htmx-go/database"
	"net"
	"net/http"
	"time"

	"github.com/ProteaTech/bok"
)

var _ API = (*api)(nil)

type API interface {
	// System Methods
	Bootstrap(database.Database) error
	IsBootstrapped() bool
	Serve() error
	Shutdown(ctx context.Context) error

	// Templates
	newTemplate(string) Template

	basePageData(ctx context.Context) *pageData

	// Grant access to the database
	DbFromContext(context.Context) (database.Database, error)

	// internals
	bindStaticRoutes()
	bindPageRoutes()
	bindAuthRoutes()
	bindBlogFeedRoutes()
	bindBlogPostRoutes()
	bindAdminRoutes()
}

type api struct {
	server   *http.Server
	router   bok.Router
	baseTmpl *Template
	fs       *embed.FS
}

func New(ctx context.Context, db database.Database, fs *embed.FS) API {
	api := api{
		fs:     fs,
		router: bok.NewRouter(),
	}

	api.server = &http.Server{
		BaseContext: func(listener net.Listener) context.Context {
			return context.WithValue(ctx, "db", db)
		},
		ConnContext: func(ctx context.Context, c net.Conn) context.Context {
			return context.WithValue(ctx, "conn", c)
		},
		ReadTimeout:       10 * time.Minute,
		ReadHeaderTimeout: 30 * time.Second,
		// WriteTimeout: 60 * time.Second, // breaks sse!
		Handler: api.router,
		Addr:    "0.0.0.0:8080",
	}
	return &api
}
