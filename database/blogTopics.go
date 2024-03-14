package database

import (
	"fmt"
	"htmx-go/models"
)

// GetBlogTopics returns all blog topics from the database.
// NOTE: This should only be called on App Bootstrap, and the result should be cached.
func (db *db) GetBlogTopics() ([]*models.BlogTopic, error) {
	var topics []*models.BlogTopic
	rows, err := db.Query("SELECT * FROM blog_topics;")
	if err != nil {
		return nil, fmt.Errorf("GetBlogTopics: %v", err)
	}

	defer rows.Close()

	for rows.Next() {
		var topic models.BlogTopic
		err = rows.Scan(&topic.ID, &topic.Name, &topic.Description)
		if err != nil {
			return nil, fmt.Errorf("GetBlogTopics: %v", err)
		}
		topics = append(topics, &topic)
	}

	return topics, nil
}

// CachedTopics returns the cached topics
func (d *db) CachedTopics() []*models.BlogTopic {
	return d.Topics
}

// GetBlogTopicByPostId returns the blog topic for the given post ID.
func (db *db) GetBlogTopicByPostId(postID int) (*models.BlogTopic, error) {
	var topic models.BlogTopic

	stmt := `
    SELECT bt.id, bt.name, bt.description 
    FROM blog_topics bt
    JOIN blog_post_topics bpt
    ON bt.id = bpt.topic_id
    WHERE bpt.post_id = ?;`

	err := db.QueryRow(stmt, postID).Scan(
		&topic.ID, &topic.Name, &topic.Description,
	)
	if err != nil {
		return nil, fmt.Errorf("GetBlogTopicByPostId: %v", err)
	}
	return &topic, nil
}
