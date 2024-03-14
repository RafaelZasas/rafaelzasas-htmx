package apiserver

import (
	"context"
	"embed"
	"fmt"
	"html/template"
	"htmx-go/database"
	"htmx-go/models"
	"htmx-go/utils"
	"log"
	"net"
	"net/http"
	"time"
)

var _ API = (*api)(nil)

type API interface {
	// System Methods
	Bootstrap() error
	IsBootstrapped() bool
	Serve() error
	Shutdown(ctx context.Context) error

	// Templates
	newTemplate(string) Template

	// Grant access to the database
	DB() database.Database

	// internals
	bindStaticRoutes()
	bindPageRoutes()
	bindAuthRoutes()
	bindBlogFeedRoutes()
	bindBlogPostRoutes()
	bindAdminRoutes()
}

// PageData is the structure that will be passed to every handler
type pageData struct {
	Title       string
	Keywords    string
	Description string
	User        *models.User
	Extra       map[string]interface{} // Extra data that can be passed to the template
}

func (p *pageData) Add(key string, value interface{}) {
	p.Extra[key] = value
}

type Template struct {
	tmpl *template.Template
}

func (a *api) newTemplate(tmplName string) Template {
	newT := template.Must(a.baseTmpl.tmpl.Clone())
	newT = template.Must(newT.ParseFS(*a.fs, tmplName))
	return Template{tmpl: newT}
}

func (t *Template) Render(w http.ResponseWriter, name string, data interface{}) {
	err := t.tmpl.ExecuteTemplate(w, name, data)
	if err != nil {
		log.Println("could not render %s: %s", name, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

type api struct {
	server   *http.Server
	router   Router
	baseTmpl *Template
	db       database.Database
	fs       *embed.FS
}

func New(ctx context.Context, db database.Database, fs *embed.FS) API {
	api := api{
		db:     db,
		fs:     fs,
		router: newRouter(),
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

func (api *api) Bootstrap() error {
	// Set up the base template with the necessary functions
	base := template.New("base")
	base.Funcs(template.FuncMap{
		"hasPermission": api.db.HasPermission,
		"formatDate": func(date time.Time) string {
			return date.Format("January 2, 2006")
		},
		"withComponentData": utils.WithComponentData,
	})
	// Parse the base template and neccessary partials and components
	base = template.Must(base.ParseFS(
		*api.fs,
		"views/base.html",
		"views/header.html",
		"views/components/*.html",
		"views/layouts/*.html",
	))

	// Initialize the base template
	api.baseTmpl = &Template{tmpl: base}

	api.bindStaticRoutes()
	api.bindPageRoutes()
	api.bindAuthRoutes()
	api.bindBlogFeedRoutes()
	api.bindBlogPostRoutes()
	api.bindAdminRoutes()
	return nil
}

func (api *api) DB() database.Database {
	return api.db
}

func (api *api) IsBootstrapped() bool {
	if api.fs == nil {
		log.Println("api filesystem is nil")
		return false
	}
	if api.db == nil {
		log.Println("api database is nil")
		return false
	}
	if api.baseTmpl == nil {
		log.Println("api base template is nil")
		return false
	}
	if api.router == nil {
		log.Println("api router is nil")
		return false
	}
	if api.server == nil {
		log.Println("api server is nil")
		return false
	}
	return true
}

func (api *api) Serve() error {
	log.Println("🌐 Starting server...")

	go func() {
		err := api.server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("server failed to start: %s", err)
		}
	}()
	log.Println("🍻 Server started on port 8080")
	return nil
}

func (api *api) Shutdown(ctx context.Context) error {
	log.Println("⚰️  Server shutting down...")
	if err := api.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("server shutdown failed: %s", err)
	}
	log.Println("✅ Server shutdown successful")
	return nil
}
