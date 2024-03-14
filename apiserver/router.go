package apiserver

import (
	"fmt"
	"htmx-go/database"
	"htmx-go/models/permissions"
	"htmx-go/utils"
	"log"
	"net/http"
	"strings"
	"time"
)

type appRouter struct {
	*http.ServeMux
}

type Router interface {
	http.Handler
	GET(string, HandlerFunc)
	POST(string, HandlerFunc)
	PUT(string, HandlerFunc)
	DELETE(string, HandlerFunc)
	PATCH(string, HandlerFunc)
	OPTIONS(string, HandlerFunc)
	HEAD(string, HandlerFunc)
	Handle(string, http.Handler)
	HandleFunc(string, http.HandlerFunc)
}

func newRouter() *appRouter {
	var router = &appRouter{
		ServeMux: http.NewServeMux(),
	}
	return router
}

// HandlerFunc is the type of the function that will be used to handle requests
type HandlerFunc func(http.ResponseWriter, *http.Request, *pageData)

// ServeHTTP implemented to satisfy the http.Handler interface
func (handler HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println(fmt.Sprintf("%s %s", r.Method, r.URL))
	data := &pageData{
		Title:       "Rafael Zasas - Portfolio and Blog",
		Keywords:    "portfolio, software, blog, technology, web development, programming, golang, go, htmx, html, css, javascript, web, development",
		Description: "Personal blog and portfolio of Rafael Zasas, a software developer and owner of Protea Technology Services LLC.",
		Extra:       make(map[string]interface{}),
	}

	authMiddleware(r, w, data)
	adminRouteGuard(r, w, data)

	// set neccessary headers
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	handler(w, r, data)
}

func adminRouteGuard(r *http.Request, w http.ResponseWriter, data *pageData) {

	if strings.Contains(r.URL.Path, "admin") == false {
		return
	}

	if data.User == nil {
		http.Redirect(w, r, "/404", http.StatusSeeOther)
		return
	}

	db := r.Context().Value("db").(database.Database)

	hasPermission, _ := db.HasPermission(data.User.RoleId, permissions.ViewAdmin)
	if !hasPermission {
		w.Header().Set("HX-Redirect", "/404")
		http.Redirect(w, r, "/404", http.StatusSeeOther)
		return
	}
}

func authMiddleware(r *http.Request, w http.ResponseWriter, data *pageData) {
	// Validate User Authentication Status and set data.User if user is authenticated
	uid, err := utils.ValidateJWT(r)

	if err != nil {
		// If the JWT is invalid, check the refresh token
		uid, err = utils.ValidateRefreshToken(r)
		// If the refresh token is valid.
		// Create a new JWT pair and set the cookies
		if err == nil {
			t, rt, err := utils.GenerateJWTokens(uid)
			if err != nil {
				log.Println(fmt.Errorf("could not generate JWT: %s", err))
			}
			http.SetCookie(w, &http.Cookie{
				Name:     "token",
				Value:    t,
				Expires:  time.Now().Add(5 * time.Minute),
				HttpOnly: true,
				Secure:   true,
			})

			http.SetCookie(w, &http.Cookie{
				Name:     "refresh_token",
				Value:    rt,
				Expires:  time.Now().Add(24 * time.Hour * 7),
				HttpOnly: true,
				Secure:   true,
			})
		}
	}

	db := r.Context().Value("db").(database.Database)

	if uid != "" {
		if user, err := db.GetUserByUID(uid); err == nil {
			data.User = user
		}
	}
}

func (r *appRouter) GET(path string, handler HandlerFunc) {
	r.Handle("GET "+path, handler)
}

func (r *appRouter) POST(path string, handler HandlerFunc) {
	r.Handle("POST "+path, handler)
}

func (r *appRouter) PUT(path string, handler HandlerFunc) {
	r.Handle("PUT "+path, handler)
}

func (r *appRouter) DELETE(path string, handler HandlerFunc) {
	r.Handle("DELETE "+path, handler)
}

func (r *appRouter) PATCH(path string, handler HandlerFunc) {
	r.Handle("PATCH "+path, handler)
}

func (r *appRouter) OPTIONS(path string, handler HandlerFunc) {
	r.Handle("OPTIONS "+path, handler)
}

func (r *appRouter) HEAD(path string, handler HandlerFunc) {
	r.Handle("HEAD "+path, handler)
}

func (r *appRouter) Handle(path string, handler http.Handler) {
	r.ServeMux.Handle(path, handler)
}

func (r *appRouter) HandleFunc(path string, handler http.HandlerFunc) {
	r.ServeMux.HandleFunc(path, handler)
}
