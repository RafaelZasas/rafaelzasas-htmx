package database

import (
	"fmt"
	"htmx-go/models"
)

// GetPostPreviewById returns a blog post by its ID
//
// NOTE: GetPostPreview will *not* return a fully populated BlogPost struct,
// it will only return the fields necessary for a preview on the blog feed
func (db *db) GetPostPreviewById(id int) (*models.BlogPost, error) {
	var post models.BlogPost
	err := db.QueryRow(
		`SELECT
        id, title, slug, content, excerpt,
        published_at, updated_at,
        is_published, image_url
        FROM blog_posts
        WHERE id = ?;`, id).Scan(
		&post.ID,
		&post.Title,
		&post.Slug,
		&post.Content,
		&post.Excerpt,
		&post.PublishedAt,
		&post.UpdatedAt,
		&post.IsPublished,
		&post.ImageURL,
	)

	if err != nil {
		return nil, fmt.Errorf("GetPostById: %v", err)
	}

	return &post, nil
}

// GetPostPreviews returns a slice of blog post previews
func (db *db) GetPostPreviews(page, pageSize int) ([]*models.BlogPost, error) {
	var posts []*models.BlogPost
	rows, err := db.Query(
		`SELECT
        id, title, slug, excerpt, published_at, updated_at, is_published
        FROM blog_posts
        WHERE is_published = true
        ORDER BY published_at DESC
        LIMIT ? OFFSET ?;`, pageSize, (page-1)*pageSize)

	if err != nil {
		return nil, fmt.Errorf("GetPostPreviews: %v", err)
	}

	defer rows.Close()

	for rows.Next() {
		var post models.BlogPost
		err = rows.Scan(
			&post.ID,
			&post.Title,
			&post.Slug,
			&post.Excerpt,
			&post.PublishedAt,
			&post.UpdatedAt,
			&post.IsPublished,
		)
		if err != nil {
			return nil, fmt.Errorf("GetPostPreviews: %v", err)
		}
		posts = append(posts, &post)
	}

	return posts, nil
}

// GetPostPreviewsByTagId returns all posts that are tagged with the given tag
func (db *db) GetPostPreviewsByTagId(tagID, startIndex, pageSize int) ([]*models.BlogPost, error) {
	var posts []*models.BlogPost
	rows, err := db.Query(
		`SELECT 
        bp.id, bp.title, bp.slug, bp.excerpt, bp.published_at,
        bp.updated_at, bp.is_published
        FROM blog_posts bp
        JOIN blog_post_tags bpt ON bp.id = bpt.post_id
        WHERE bpt.tag_id = ?
        ORDER BY bp.published_at DESC
        LIMIT ? OFFSET ?;`, tagID, pageSize, startIndex)

	if err != nil {
		return nil, fmt.Errorf("GetPostsByTagId: %v", err)
	}

	defer rows.Close()

	for rows.Next() {
		var post models.BlogPost
		err = rows.Scan(
			&post.ID,
			&post.Title,
			&post.Slug,
			&post.Excerpt,
			&post.PublishedAt,
			&post.UpdatedAt,
			&post.IsPublished,
		)
		if err != nil {
			return nil, fmt.Errorf("GetPostsByTagId: %v", err)
		}
		posts = append(posts, &post)
	}

	return posts, nil
}

// GetPostPreviewsByTopicId returns all posts that are tagged with the given topic
func (db *db) GetPostPreviewsByTopicId(topicID int) ([]*models.BlogPost, error) {
	var posts []*models.BlogPost
	rows, err := db.Query(
		`SELECT 
        bp.id, bp.title, bp.slug, bp.excerpt, bp.published_at,
        bp.updated_at, bp.is_published
        FROM blog_posts bp
        JOIN blog_post_tags bpt ON bp.id = bpt.post_id
        JOIN tags t ON bpt.tag_id = t.id
        WHERE t.topic_id = ?;`, topicID)

	if err != nil {
		return nil, fmt.Errorf("GetPostsByTopicId: %v", err)
	}

	defer rows.Close()

	for rows.Next() {
		var post models.BlogPost
		err = rows.Scan(
			&post.ID,
			&post.Title,
			&post.Slug,
			&post.Excerpt,
			&post.PublishedAt,
			&post.UpdatedAt,
			&post.IsPublished,
		)
		if err != nil {
			return nil, fmt.Errorf("GetPostsByTopicId: %v", err)
		}
		posts = append(posts, &post)
	}

	return posts, nil
}

// SearchPosts uses FTS5 to search for posts that contain the given term
func (db *db) SearchPosts(searchTerm string) ([]*models.BlogPost, error) {
	var posts []*models.BlogPost
	rows, err := db.Query(
		`SELECT
        rowID
        FROM blog_posts_fts
        WHERE blog_posts_fts MATCH ?
        ORDER BY RANK;`, searchTerm)

	if err != nil {
		return nil, fmt.Errorf("SearchPosts: %v", err)
	}

	defer rows.Close()

	for rows.Next() {
		var post *models.BlogPost
		var id int
		err = rows.Scan(
			&id,
		)
		if err != nil {
			return nil, fmt.Errorf("SearchPosts: %v", err)
		}
		post, err = db.GetPostPreviewById(id)
		if err != nil {
			return nil, fmt.Errorf("SearchPosts: %v", err)
		}
		posts = append(posts, post)
	}

	return posts, nil

}
