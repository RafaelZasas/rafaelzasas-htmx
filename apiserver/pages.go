package apiserver

import (
	"fmt"
	"net/http"
	"strings"
)

func (api *api) bindPageRoutes() {
	// Routes
	api.router.GET("/{$}", api.handleIndex)
	api.router.GET("/about", api.handleAbout)
	api.router.GET("/contact", api.handleContactView)
	api.router.GET("/login", api.handleLoginView)
	api.router.GET("/register", api.handleSignUpView)
	api.router.GET("/profile", api.handleProfileView)
	api.router.GET("/services", api.handleServicesView)
	// this matches any route tapi.hat is not found
	api.router.GET("/", api.handleNotFound)
	// explicit 404 for route protection ambiguity
	api.router.GET("/404", api.handleNotFound)

}

func (api *api) handleIndex(w http.ResponseWriter, r *http.Request, data *pageData) {
	fmt.Println("handleIndex")
	var tmpl = api.newTemplate("views/index.html")
	tmpl.Render(w, "base", data)
}

func (api *api) handleNotFound(w http.ResponseWriter, r *http.Request, data *pageData) {
	fmt.Println("handleNotFound")
	var tmpl = api.newTemplate("views/404.html")
	data.Title = "404 | Rafael Zasas"
	w.WriteHeader(http.StatusNotFound)
	// set the max cache control since this is a static page
	w.Header().Add("cache-control", "max-age=31536000, public")
	tmpl.Render(w, "base", data)
}

func (api *api) handleAbout(w http.ResponseWriter, r *http.Request, data *pageData) {
	fmt.Println("handleAbout")
	var tmpl = api.newTemplate("views/about.html")
	data.Title = "About | Rafael Zasas"
	tmpl.Render(w, "base", data)
}

func (api *api) handleServicesView(w http.ResponseWriter, r *http.Request, data *pageData) {
	fmt.Println("handleServicesView")
	var tmpl = api.newTemplate("views/services.html")
	data.Title = "Services | Rafael Zasas"
	tmpl.Render(w, "base", data)
}

func (api *api) handleContactView(w http.ResponseWriter, r *http.Request, data *pageData) {
	fmt.Println("handleContactView")
	var tmpl = api.newTemplate("views/contact.html")
	data.Title = "Contact | Rafael Zasas"
	tmpl.Render(w, "base", data)
}

func (api *api) handleLoginView(w http.ResponseWriter, r *http.Request, data *pageData) {
	if data.User != nil {
		w.Header().Add("HX-Redirect", "/")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	var tmpl = api.newTemplate("views/login.html")
	data.Title = "Login | Rafael Zasas"
	referer := r.Header.Get("HX-Current-Url")
	// make sure that the referer is this server, and not another site
	if !strings.Contains(referer, r.Host) {
		referer = "/"
	}
	data.Add("Redirect", referer)
	tmpl.Render(w, "base", data)
}

func (api *api) handleSignUpView(w http.ResponseWriter, r *http.Request, data *pageData) {
	fmt.Println("handleSignUpView")
	var tmpl = api.newTemplate("views/register.html")
	data.Title = "Register | Rafael Zasas"
	tmpl.Render(w, "base", data)
}

func (api *api) handleProfileView(w http.ResponseWriter, r *http.Request, data *pageData) {
	if data.User == nil {
		w.Header().Add("HX-Redirect", "/login")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	var tmpl = api.newTemplate("views/profile.html")
	data.Title = "Profile | Rafael Zasas"
	w.Header().Add("cache-control", "no-store, no-cache, private")
	tmpl.Render(w, "base", data)
}
