package apiserver

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"
)

func (api *api) bindBlogPostRoutes() {
	api.router.GET("/blog/{slug}", api.handleBlogPostPage)
	api.router.GET("/blog/edit/{slug}", api.handleBlogEditPage)
	api.router.PATCH("/blog/{slug}", api.handleEditPost)
}

func (api *api) handleBlogPostPage(w http.ResponseWriter, r *http.Request) {
	data := api.basePageData(r.Context())
	db, err := api.DbFromContext(r.Context())
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	tmpl := api.newTemplate("views/blog/blogPage.html")
	slug := r.PathValue("slug")
	post, err := db.GetPostBySlug(slug)

	if err == sql.ErrNoRows {
		http.Redirect(w, r, "/404", http.StatusSeeOther)
		return
	}

	if err != nil {
		log.Println(fmt.Errorf("api.handleBlogPostPage: %w", err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	var keywords = ""
	for _, tag := range post.Tags {
		keywords = fmt.Sprintf("%s, %s", keywords, tag.Name)
	}

	if err := db.UpdatePostViewCount(post.ID); err != nil {
		log.Println(fmt.Errorf("api.handleBlogPostPage: %w", err))
	}

	// make sure that html content is rendered and not escaped

	content := template.HTML(post.Content)

	data.Title = post.Title
	data.Description = post.Excerpt
	data.Keywords = keywords
	data.Add("Post", post)
	data.Add("Content", content)

	tmpl.Render(w, "base", data)
}

// (api *api) handleBlogEditPage handles the edit post page
func (api *api) handleBlogEditPage(w http.ResponseWriter, r *http.Request) {
	data := api.basePageData(r.Context())

	db, err := api.DbFromContext(r.Context())
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	tmpl := api.newTemplate("views/blog/editPost.html")
	slug := r.PathValue("slug")
	post, err := db.GetPostBySlug(slug)

	if err == sql.ErrNoRows {
		http.Redirect(w, r, "/404", http.StatusSeeOther)
		return
	}

	if err != nil {
		log.Println(fmt.Errorf("api.handleBlogEditPage: %w", err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data.Title = "Edit Post"
	data.Add("Post", post)

	tmpl.Render(w, "base", data)
}

func (api *api) handleEditPost(w http.ResponseWriter, r *http.Request) {
	db, err := api.DbFromContext(r.Context())
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	slug := r.PathValue("slug")
	post, err := db.GetPostBySlug(slug)

	if err == sql.ErrNoRows {
		http.Redirect(w, r, "/404", http.StatusSeeOther)
		return
	}

	if err != nil {
		log.Println(fmt.Errorf("api.handleEditPost: %w", err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if err := r.ParseForm(); err != nil {
		log.Println(fmt.Errorf("api.handleEditPost: %w", err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	post.Title = r.FormValue("Title")
	post.Content = r.FormValue("Content")
	post.Excerpt = r.FormValue("Excerpt")
	post.UpdatedAt = time.Now()

	fmt.Println(post)
	if err := db.UpdatePost(post); err != nil {
		log.Println(fmt.Errorf("api.handleEditPost: %w", err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Add("HX-Redirect", fmt.Sprintf("/blog/%s", post.Slug))
	w.WriteHeader(http.StatusCreated)
}
