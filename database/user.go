package database

import (
	"fmt"
	"htmx-go/models"
)

// GetUsers returns all users from the database
func (db *db) GetUsers() ([]*models.User, error) {
	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]*models.User, 0)
	for rows.Next() {
		user := &models.User{}
		err := rows.Scan(
			&user.UID,
			&user.Email,
			&user.Password,
			&user.Name,
			&user.Provider,
			&user.EmailVerified,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.RefreshToken,
			&user.RoleId,
			&user.Avatar,
			&user.TwitterUsername,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

func (db *db) GetUserByEmail(email string) (*models.User, error) {
	user := &models.User{}
	row := db.QueryRow("SELECT * FROM users WHERE email = ?", email)
	err := row.Scan(
		&user.UID,
		&user.Email,
		&user.Password,
		&user.Name,
		&user.Provider,
		&user.EmailVerified,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.RefreshToken,
		&user.RoleId,
		&user.Avatar,
		&user.TwitterUsername,
	)
	if err != nil {
		return nil, fmt.Errorf("GetUserByEmail: %v", err)
	}
	return user, nil
}

func (db *db) GetUserByUID(uid string) (*models.User, error) {
	user := &models.User{}
	row := db.QueryRow("SELECT * FROM users WHERE uid = ?", uid)
	err := row.Scan(
		&user.UID,
		&user.Email,
		&user.Password,
		&user.Name,
		&user.Provider,
		&user.EmailVerified,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.RefreshToken,
		&user.RoleId,
		&user.Avatar,
		&user.TwitterUsername,
	)
	if err != nil {
		return nil, fmt.Errorf("GetUserByUID: %v", err)
	}
	return user, nil
}

func (db *db) CreateUser(user *models.User) error {
	_, err := db.Exec(
		`INSERT INTO users 
        (
            uid,
            email,
            password,
            name,
            provider,
            email_verified,
            created_at,
            updated_at,
            refresh_token,
            role_id
        )
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		user.UID,
		user.Email,
		*user.Password,
		user.Name,
		user.Provider,
		user.EmailVerified,
		user.CreatedAt,
		user.UpdatedAt,
		*user.RefreshToken,
		user.RoleId,
	)
	return err
}

func (db *db) GetUserRefreshToken(uid string) (string, error) {
	var refreshToken string
	row := db.QueryRow("SELECT refresh_token FROM users WHERE uid = ?", uid)
	err := row.Scan(&refreshToken)
	if err != nil {
		return "", err
	}
	return refreshToken, nil
}

func (db *db) UpdateUserRefreshToken(uid, refreshToken string) error {
	_, err := db.Exec("UPDATE users SET refresh_token = ? WHERE uid = ?", refreshToken, uid)
	return err
}
