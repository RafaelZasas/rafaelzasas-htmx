package database

import (
	"fmt"
	"htmx-go/models"
	"log"
)

// GetBlogTags returns all blog tags from the database.
// NOTE: This should only be called on App Bootstrap, and the result should be cached.
func (db *db) GetBlogTags() ([]*models.BlogTag, error) {
	var tags []*models.BlogTag
	rows, err := db.Query("SELECT * FROM blog_tags;")
	if err != nil {
		log.Println("failed to query db for blog tags")
		return nil, fmt.Errorf("GetBlogTags: %v", err)
	}

	defer rows.Close()

	for rows.Next() {
		var tag models.BlogTag
		err = rows.Scan(&tag.ID, &tag.Name, &tag.Description)
		if err != nil {
			return nil, fmt.Errorf("GetBlogTags: %v", err)
		}
		tags = append(tags, &tag)
	}

	return tags, nil
}

func (d *db) CachedTags() []*models.BlogTag {
	return d.Tags
}

// GetTagsByTopicId returns all tags for a given topic
func (db *db) GetBlogTagsByTopicId(topicID int) ([]*models.BlogTag, error) {
	var tags []*models.BlogTag
	rows, err := db.Query(
		`SELECT bt.id, bt.name, bt.description 
        FROM blog_tags bt
        JOIN blog_tag_topics btt ON bt.id = btt.tag_id
        WHERE btt.topic_id = ?;`, topicID)

	if err != nil {
		return nil, fmt.Errorf("GetTagsForTopic: %v", err)
	}

	defer rows.Close()

	for rows.Next() {
		var tag models.BlogTag
		err = rows.Scan(&tag.ID, &tag.Name, &tag.Description)
		if err != nil {
			return nil, fmt.Errorf("GetTagsForTopic: %v", err)
		}
		tags = append(tags, &tag)
	}

	return tags, nil
}

// GetBlogTagsByPostId returns all tags for a given post
func (db *db) GetBlogTagsByPostId(postID int) ([]*models.BlogTag, error) {
	var tags []*models.BlogTag
	rows, err := db.Query(
		`SELECT bt.id, bt.name, bt.description 
        FROM blog_tags bt
        JOIN blog_post_tags bpt ON bt.id = bpt.tag_id
        WHERE bpt.post_id = ?;`, postID)

	if err != nil {
		return nil, fmt.Errorf("GetTagsForPost: %v", err)
	}

	defer rows.Close()

	for rows.Next() {
		var tag models.BlogTag
		err = rows.Scan(&tag.ID, &tag.Name, &tag.Description)
		if err != nil {
			return nil, fmt.Errorf("GetTagsForPost: %v", err)
		}
		tags = append(tags, &tag)
	}

	return tags, nil
}
