package apiserver

import (
	"bytes"
	"fmt"
	"html/template"
	"htmx-go/models"
	"log"
	"net/http"
	"strconv"
	"strings"
)

const pageSize = 5

func (api *api) bindBlogFeedRoutes() {
	api.router.GET("/blog", api.handleBlogFeedPage)
	api.router.GET("/blog/previews", api.handleBlogPreviews)
	api.router.GET("/blog/search", api.handleBlogSearch)
	api.router.GET("/blog/topics", api.handleBlogTopics)
	api.router.GET("/blog/tags", api.handleBlogTags)
}

func (api *api) handleBlogFeedPage(w http.ResponseWriter, r *http.Request, data *pageData) {
	tmpl := api.newTemplate("views/blog/blogFeed.html")
	data.Title = "Blog | Rafael Zasas"
	data.Description = "A blog created using Go on the streets and HTMX in the sheets."
	data.Keywords = "blog, go, golang, htmx, html, css, javascript, web, development"
	// cache for a week since I dont think I will be posting more than that
	w.Header().Add("cache-control", "max-age=604800, public")
	tmpl.Render(w, "base", data)
}

func (api *api) handleBlogPreviews(w http.ResponseWriter, r *http.Request, _ *pageData) {
	baseTemplate := api.newTemplate("views/blog/blogFeed.html")

	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 1 // default to the first page if no page is specified or if there is an error
	}

	tagId, err := strconv.Atoi(r.URL.Query().Get("tag"))
	if err != nil {
		tagId = 0
	}

	startIndex := (page - 1) * pageSize

	// Get the posts for the current page
	var posts []*models.BlogPost
	if tagId > 0 {
		posts, err = api.DB().GetPostPreviewsByTagId(tagId, startIndex, pageSize)
		if err != nil {
			log.Println(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	} else {
		posts, err = api.DB().GetPostPreviews(startIndex, pageSize)
		if err != nil {
			log.Println(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	if len(posts) == 0 {
		baseTemplate.Render(w, "noPosts", nil)
		return
	}

	var buff bytes.Buffer
	for i, post := range posts {
		tmpl := template.Must(baseTemplate.tmpl.Parse(`{{template "snippet" .}}`))
		tmpl = template.Must(tmpl.Clone())

		if err := tmpl.Execute(&buff, post); err != nil {
			log.Println(fmt.Errorf("Error executing snippet template: %w", err))
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Add the element that will trigger the next page load
		if i == len(posts)-2 {
			fmt.Fprintf(&buff,
				`<span 
    id="nextPage"
    hx-get="/blog/previews?page=%d&tag=%d"
    hx-target="#feed"
    hx-swap="beforeend"
    hx-trigger="intersect once"
    hx-indicator="#blocks-shuffle-indicator"
    ></span>`, page+1, tagId)
		}
	}

	tmpl, err := template.New("").Parse(buff.String())
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
	return
}

func (api *api) handleBlogSearch(w http.ResponseWriter, r *http.Request, _ *pageData) {
	baseTemplate := api.newTemplate("views/blog/blogFeed.html")
	searchQuery := r.URL.Query().Get("q")

	if searchQuery == "" {
		posts, err := api.DB().GetPostPreviews(0, pageSize)
		if err != nil {
			log.Println(fmt.Errorf("(api *api) handleBlogSearch: %w", err))
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		var buff bytes.Buffer
		for _, post := range posts {
			tmpl := template.Must(baseTemplate.tmpl.Parse(`{{template "snippet" .}}`))
			tmpl = template.Must(tmpl.Clone())

			if err := tmpl.Execute(&buff, post); err != nil {
				log.Println(fmt.Errorf("Error executing snippet template: %w", err))
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

		}

		tmpl, err := template.New("").Parse(buff.String())
		if err != nil {
			log.Println(fmt.Errorf("(api *api) handleBlogSearch: %w", err))
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
	} else {
		posts, err := api.DB().SearchPosts(searchQuery)
		if err != nil {
			log.Println(fmt.Errorf("Error executing snippet template: %w", err))
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		if len(posts) == 0 {
			baseTemplate.Render(w, "noPosts", nil)
			return
		}

		var buff bytes.Buffer
		for _, post := range posts {
			tmpl := template.Must(baseTemplate.tmpl.Parse(`{{template "snippet" .}}`))
			tmpl = template.Must(tmpl.Clone())

			if err := tmpl.Execute(&buff, post); err != nil {
				log.Println(fmt.Errorf("Error executing snippet template: %w", err))
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
		}

		tmpl, err := template.New("").Parse(buff.String())
		if err != nil {
			log.Println(fmt.Errorf("Error executing snippet template: %w", err))
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		if err := tmpl.Execute(w, nil); err != nil {
			log.Println(fmt.Errorf("(api *api) handleBlogSearch: %w", err))
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}

	}
}

func (api *api) handleBlogTopics(w http.ResponseWriter, r *http.Request, data *pageData) {
	var optionTags = make([]string, len(api.DB().CachedTopics()))
	for i, topic := range api.DB().CachedTopics() {
		optionTags[i] = fmt.Sprintf(
			"<option value=\"%d\">%s</option>",
			topic.ID, topic.Name)
	}
	html := template.HTML(strings.Join(optionTags, ""))

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, "%s", html)
	return
}

func (api *api) handleBlogTags(w http.ResponseWriter, r *http.Request, data *pageData) {
	tags := api.DB().CachedTags()

	// Filter tags if a topicId is provided
	topicIdStr := r.URL.Query().Get("topic")
	topicId, err := strconv.Atoi(topicIdStr)
	if err == nil {
		if t, err := api.DB().GetBlogTagsByTopicId(topicId); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		} else {
			tags = t
		}
	}

	var optionTags = make([]string, len(tags)+1)
	optionTags[0] = "<option value=\"0\">All Tags</option>"
	for i, tag := range tags {
		optionTags[i+1] = fmt.Sprintf(
			"<option value=\"%d\">%s</option>",
			tag.ID, tag.Name)
	}
	html := template.HTML(strings.Join(optionTags, ""))

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, "%s", html)
	return
}
