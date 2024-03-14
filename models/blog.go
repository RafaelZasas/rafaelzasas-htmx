package models

import (
	"time"
)

// Post represents a blog post.
type BlogPost struct {
	ID      int    `db:"id"`
	Title   string `db:"title"`
	Slug    string `db:"slug"`
	Content string `db:"content"`
	Excerpt string `db:"excerpt"`
	// NOTE: *In the db this is author_uid of type TEXT NOT NULL,
	// but in the struct it is Author of type *models.User
	Author      *User
	PublishedAt time.Time `db:"published_at"`
	UpdatedAt   time.Time `db:"updated_at"`
	IsPublished bool      `db:"is_published"`
	ViewCount   int       `db:"view_count"`
	ImageURL    string    `db:"image_url"`
	// NOTE: This is not present in the blog_posts table
	// it will be populated from memory during the database query
	Tags []*BlogTag
	// NOTE: This too is not present in the blog_posts table
	// it will also be populated from memory during the database query
	Topic *BlogTopic
}

// BlogTag represents a tag for categorizing posts.
type BlogTag struct {
	ID          int    `db:"id"`
	Name        string `db:"name"`
	Description string `db:"description"`
}

// BlogTopic represents a group of tags.
type BlogTopic struct {
	ID          int    `db:"id"`
	Name        string `db:"name"`
	Description string `db:"description"`
}
