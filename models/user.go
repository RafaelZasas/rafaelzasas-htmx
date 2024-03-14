package models

import "time"

type User struct {
	UID             string    `json:"uid" db:"uid"`
	Email           string    `json:"email" db:"email"`
	Name            string    `json:"name" db:"name"`
	Password        *string   `json:"password,omitEmpty" db:"password"`
	Provider        string    `json:"provider" db:"provider"`
	EmailVerified   bool      `json:"emailVerified" db:"email_verified"`
	CreatedAt       time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt       time.Time `json:"updatedAt" db:"updated_at"`
	RefreshToken    *string   `json:"-" db:"refresh_token"`
	RoleId          int       `json:"roleId" db:"role_id"`
	Avatar          *string   `json:"avatar" db:"avatar"`
	TwitterUsername *string   `json:"twitterUsername" db:"twitter_username"`
}
