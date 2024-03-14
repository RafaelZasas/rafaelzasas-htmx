package apiserver

import (
	"context"
	"fmt"
	"html/template"
	"htmx-go/models"
	"htmx-go/models/permissions"
	"net/http"
)

type AdminLink struct {
	URL        string
	Label      string
	IsActive   bool
	Permission models.Permission
}

var adminLinks = []AdminLink{
	{
		URL: "/admin", Label: "Dashboard",
		Permission: permissions.ViewAdmin,
		IsActive:   true,
	},
	{
		URL: "/admin/users", Label: "Users",
		Permission: permissions.ReadUsers,
		IsActive:   false,
	},
	{
		URL: "/admin/blog", Label: "Blog",
		Permission: permissions.UpdateBlogPost,
		IsActive:   false,
	},
	{
		URL: "/admin/permissions", Label: "Permissions",
		Permission: permissions.ReadPermissions,
		IsActive:   false,
	},
	{
		URL: "/admin/roles", Label: "Roles",
		Permission: permissions.ReadPermissions,
		IsActive:   false,
	},
}

func (api *api) bindAdminRoutes() {
	router := api.router.WithMiddleware(adminRouteGuard())

	// Authentication
	router.GET("/admin", api.handleAdminPage)
	router.GET("/admin/users", api.handleAdminUsersPage)
	router.DELETE("/admin/users/{uid}", api.handleDeleteUser)
}

func (api *api) newAdminTemplate(ctx context.Context, tmplName, route string) (Template, []AdminLink) {
	db, _ := api.DbFromContext(ctx)
	newT := template.Must(api.baseTmpl.tmpl.Clone())
	newT.Funcs(template.FuncMap{
		"getRoleName": db.GetRoleName,
	})
	newT = template.Must(newT.ParseFS(*api.fs, tmplName))

	for i := range adminLinks {
		adminLinks[i].IsActive = false
		if adminLinks[i].URL == route {
			adminLinks[i].IsActive = true
		}
	}
	return Template{tmpl: newT}, adminLinks
}

func (api *api) handleAdminPage(w http.ResponseWriter, r *http.Request) {

	data := api.basePageData(r.Context())

	var tmpl, links = api.newAdminTemplate(
		r.Context(), "views/admin/dashboard.html", "/admin",
	)

	fmt.Println(data.User)
	data.Extra["AdminNavLinks"] = links
	tmpl.Render(w, "adminLayout", data)
	return
}

func (api *api) handleAdminUsersPage(w http.ResponseWriter, r *http.Request) {
	db, _ := api.DbFromContext(r.Context())

	data := api.basePageData(r.Context())

	var tmpl, links = api.newAdminTemplate(
		r.Context(), "views/admin/users.html", "/admin/users",
	)

	users, err := db.GetUsers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data.Title = "Admin | Go-Htmx Example"
	data.Extra["AdminNavLinks"] = links
	data.Extra["Users"] = users
	tmpl.Render(w, "adminLayout", data)
	return
}

func (api *api) handleDeleteUser(w http.ResponseWriter, r *http.Request) {

	uid := r.PathValue("uid")
	if uid == "" {
		http.Error(w, "uid is required", http.StatusBadRequest)
		return
	}
	//
	//   err := api.db.DeleteUser(uid)
	//   if err != nil {
	//       http.Error(w, err.Error(), http.StatusInternalServerError)
	//       return
	//   }
	w.WriteHeader(http.StatusAccepted)
}
