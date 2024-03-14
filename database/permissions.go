package database

import "htmx-go/models"

// HasPermission checks if the role has the specified permission.
func (db *db) HasPermission(uid string, p models.Permission) (bool, error) {
	userRoleId, err := db.GetUserRole(uid)
	if err != nil {
		return false, err
	}

	var exists bool = false
	query := `SELECT EXISTS(
        SELECT 1 FROM role_permissions WHERE role_id = ? AND permission_id = (
            SELECT id FROM permissions WHERE name = ?
        )
    )`
	err = db.QueryRow(query, userRoleId, p).Scan(&exists)
	return exists, err
}

func (db *db) GetRoleName(roleId int) (string, error) {
	var roleName string
	query := `SELECT name FROM roles WHERE id = ?`
	err := db.QueryRow(query, roleId).Scan(&roleName)
	return roleName, err
}
