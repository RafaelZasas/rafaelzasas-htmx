package permissions

import "htmx-go/models"

const (
	ViewAdmin         models.Permission = "view_admin"
	ReadUsers         models.Permission = "read_users"
	DeleteUser        models.Permission = "delete_user"
	AddBlogPost       models.Permission = "add_blog_post"
	UpdateBlogPost    models.Permission = "update_blog_post"
	DeleteBlogPost    models.Permission = "delete_blog_post"
	AddComment        models.Permission = "add_comment"
	UpdatePermissions models.Permission = "update_permissions"
	ReadPermissions   models.Permission = "read_permissions"
	CreatePermissions models.Permission = "create_permissions"
	DeletePermissions models.Permission = "delete_permissions"
)
