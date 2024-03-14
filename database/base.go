package database

import (
	"database/sql"
	"fmt"
	"htmx-go/models"
	"log"
	"os"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
	_ "modernc.org/sqlite"
)

var _ Database = (*db)(nil)

type Database interface {
	// Database Methods
	Bootstrap() error
	IsBootstrapped() bool
	runMigrations() error
	Close() error

	// User Methods
	CreateUser(*models.User) error
	GetUsers() ([]*models.User, error)
	GetUserByEmail(string) (*models.User, error)
	GetUserByUID(string) (*models.User, error)
	GetUserRefreshToken(string) (string, error)
	UpdateUserRefreshToken(string, string) error

	// Permission Methods
	HasPermission(int, models.Permission) (bool, error)
	GetRoleName(int) (string, error)

	// Blog Posts Preview Methods
	GetPostPreviews(int, int) ([]*models.BlogPost, error)
	GetPostPreviewsByTagId(int, int, int) ([]*models.BlogPost, error)
	GetPostPreviewById(int) (*models.BlogPost, error)
	GetPostPreviewsByTopicId(int) ([]*models.BlogPost, error)
	SearchPosts(string) ([]*models.BlogPost, error)

	// Blog Post Methods
	CreatePost(*models.BlogPost) (int, error)
	GetPostBySlug(string) (*models.BlogPost, error)
	UpdatePost(*models.BlogPost) error
	UpdatePostViewCount(int) error
	DeletePost(int) error

	// Blog Tag Methods
	GetBlogTags() ([]*models.BlogTag, error)
	CachedTags() []*models.BlogTag // Returns the cached tags
	GetBlogTagsByTopicId(int) ([]*models.BlogTag, error)
	GetBlogTagsByPostId(int) ([]*models.BlogTag, error)

	// Blog Topic Methods
	GetBlogTopics() ([]*models.BlogTopic, error)
	CachedTopics() []*models.BlogTopic // Returns the cached topics
	GetBlogTopicByPostId(int) (*models.BlogTopic, error)
}

type db struct {
	*sql.DB
	Tags   []*models.BlogTag
	Topics []*models.BlogTopic
}

// New returns a new instance of a Database
// NOTE: The underlying struct type is a pointer to DB
func New() (Database, error) {
	log.Println("🗃️ Connecting to database...")

	tursoUrl := os.Getenv("TURSO_URL")
	tursoToken := os.Getenv("TURSO_TOKEN")

	var url string

	if os.Getenv("ENV") == "production" {
		url = fmt.Sprintf("%s?authToken=%s", tursoUrl, tursoToken)
	} else {
		url = "file:htmx-go.db"
	}

	sqlDb, err := sql.Open("libsql", url)
	if err != nil {
		return nil, fmt.Errorf("database.Connect: %q\n", err)
	}

	baseDb := db{DB: sqlDb}

	return &baseDb, nil
}

func (db *db) Bootstrap() error {

	// Keeping this as a private function to avoid fuck ups
	if err := db.runMigrations(); err != nil {
		return fmt.Errorf("database.New: %q\n", err)
	}

	// cache tags and topics
	tags, err := db.GetBlogTags()
	if err != nil {
		return fmt.Errorf("database.New: %q\n", err)
	}
	db.Tags = tags

	topics, err := db.GetBlogTopics()
	if err != nil {
		return fmt.Errorf("database.New: %q\n", err)
	}
	db.Topics = topics
	return nil
}

func (db *db) Close() error {
	log.Println("🗃️ Closing database...")
	if err := db.DB.Close(); err != nil {
		return fmt.Errorf("database.Close: %q", err)
	}
	log.Println("✅ Database closed")
	return nil
}
