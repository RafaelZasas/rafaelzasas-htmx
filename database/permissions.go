package database

import "htmx-go/models"

// HasPermission checks if the role has the specified permission.
func (db *db) HasPermission(roleId int, permission models.Permission) (bool, error) {
	var exists bool = false
	query := `SELECT EXISTS(
        SELECT 1 FROM role_permissions WHERE role_id = ? AND permission_id = (
            SELECT id FROM permissions WHERE name = ?
        )
    )`
	err := db.QueryRow(query, roleId, permission).Scan(&exists)
	return exists, err
}

func (db *db) GetRoleName(roleId int) (string, error) {
	var roleName string
	query := `SELECT name FROM roles WHERE id = ?`
	err := db.QueryRow(query, roleId).Scan(&roleName)
	return roleName, err
}
