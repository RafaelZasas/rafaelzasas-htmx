package apiserver

import (
	"context"
	"errors"
	"fmt"
	"html/template"
	"htmx-go/database"
	"htmx-go/models"
	"htmx-go/utils"
	"log"
	"log/slog"
	"net/http"
	"time"

	"github.com/ProteaTech/bok/middleware"
)

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
		slog.Error("could not render %s: %s", name, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (api *api) Bootstrap(db database.Database) error {

	api.router = api.router.WithMiddleware(
		authMiddleware(),
		middleware.Logger(slog.Default()),
		middleware.RecoverPanic(),
		middleware.Timeout(time.Minute*5),
	)

	// Set up the base template with the necessary functions
	base := template.New("base")
	base.Funcs(template.FuncMap{
		"hasPermission": db.HasPermission,
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

func (api *api) DbFromContext(ctx context.Context) (database.Database, error) {
	db := ctx.Value("db")
	if db == nil {
		return nil, errors.New("db not in context")
	}
	return db.(database.Database), nil
}

func (api *api) basePageData(ctx context.Context) *pageData {

	defaultPageData := &pageData{
		Title:       "Rafael Zasas - Portfolio and Blog",
		Keywords:    "portfolio, software, blog, technology, web development, programming, golang, go, htmx, html, css, javascript, web, development",
		Description: "Personal blog and portfolio of Rafael Zasas, a software developer and owner of Protea Technology Services LLC.",
		Extra:       make(map[string]interface{}),
	}

	uid := ctx.Value("uid")
	if uid != nil {
		db, err := api.DbFromContext(ctx)
		if err != nil {
			slog.Error("could not get db from context", "err", err)
		}
		user, err := db.GetUserByUID(uid.(string))
		if err != nil {
			slog.Error("could not get user from context", "err", err)
		}
		defaultPageData.User = user
	}

	return defaultPageData

}

func (api *api) IsBootstrapped() bool {
	if api.fs == nil {
		log.Println("api filesystem is nil")
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
