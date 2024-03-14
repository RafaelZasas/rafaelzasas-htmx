package database

import (
	"database/sql"
	"fmt"
	"htmx-go/models"
)

// CreatePost creates a new blog post in the database
func (db *db) CreatePost(post *models.BlogPost) (int, error) {
	stmt := `
        INSERT INTO blog_posts
        (author_uid,
        title,
        slug,
        excerpt,
        content,
        published_at,
        updated_at,
        is_published,
        image_url)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?);
    `
	result, err := db.Exec(
		stmt,
		post.Author.UID,
		post.Title,
		post.Slug,
		post.Excerpt,
		post.Content,
		post.PublishedAt,
		post.UpdatedAt,
		post.IsPublished,
		post.ImageURL,
	)
	if err != nil {
		return 0, fmt.Errorf("CreatePost: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("CreatePost: %v", err)
	}
	return int(id), nil
}

// GetPost returns a blog post by its slug
//
// GetPost will return a fully populated BlogPost struct,
// including the author, tags, and topic
func (db *db) GetPostBySlug(slug string) (*models.BlogPost, error) {
	var post models.BlogPost
	var authorUID string
	err := db.QueryRow(
		`SELECT
        id, author_uid, title, slug,
        excerpt, content,
        published_at, updated_at,
        is_published, image_url, view_count
        FROM blog_posts
        WHERE slug = ?;`, slug).Scan(
		&post.ID,
		&authorUID,
		&post.Title,
		&post.Slug,
		&post.Excerpt,
		&post.Content,
		&post.PublishedAt,
		&post.UpdatedAt,
		&post.IsPublished,
		&post.ImageURL,
		&post.ViewCount,
	)

	if err == sql.ErrNoRows {
		return nil, err
	}

	if err != nil {
		return nil, fmt.Errorf("GetPostBySlug: %v", err)
	}

	if user, err := db.GetUserByUID(authorUID); err != nil {
		return nil, fmt.Errorf("GetPostBySlug: %v", err)
	} else {
		post.Author = user
	}

	// Get the tags for the post
	if tags, err := db.GetBlogTagsByPostId(post.ID); err != nil {
		return nil, fmt.Errorf("GetPostBySlug: %v", err)
	} else {
		post.Tags = tags
	}

	// Get the topic for the post
	if topic, err := db.GetBlogTopicByPostId(post.ID); err != nil {
		return nil, fmt.Errorf("GetPostBySlug: %v", err)
	} else {
		post.Topic = topic
	}

	return &post, nil
}

// UpdatePost updates a blog post in the database
func (db *db) UpdatePost(post *models.BlogPost) error {
	stmt := `
        UPDATE blog_posts
        SET title = ?, slug = ?, content = ?,
        published_at = ?, updated_at = ?,
        is_published = ?, image_url = ?
        WHERE id = ?;
    `
	_, err := db.Exec(
		stmt,
		post.Title,
		post.Slug,
		post.Content,
		post.PublishedAt,
		post.UpdatedAt,
		post.IsPublished,
		post.ImageURL,
		post.ID,
	)
	if err != nil {
		return fmt.Errorf("UpdatePost: %v", err)
	}
	return nil
}

// UpdatePostViewCount updates the view count for a blog post
func (db *db) UpdatePostViewCount(id int) error {
	stmt := `
        UPDATE blog_posts
        SET view_count = view_count + 1
        WHERE id = ?;
    `
	_, err := db.Exec(stmt, id)
	if err != nil {
		return fmt.Errorf("UpdatePostViewCount: %v", err)
	}
	return nil
}

// DeletePost deletes a blog post from the database
func (db *db) DeletePost(id int) error {
	stmt := `DELETE FROM blog_posts WHERE id = ?;`
	_, err := db.Exec(stmt, id)
	if err != nil {
		return fmt.Errorf("DeletePost: %v", err)
	}
	return nil
}
