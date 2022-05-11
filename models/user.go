package models

import "github.com/golang-jwt/jwt"

type User struct {
	Base
	Username string `gorm:"type: varchar(255)" json:"username"`
	Email    string `gorm:"uniqueIndex;type: varchar(255)" json:"email"`
	Password string `gorm:"->;<-;not null" json:"-"`
	Token    string `gorm:"-" json:"token,omitempty"`
}

// UserErrors represent the error format for user routes
type UserErrors struct {
	Err      bool   `json:"error"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Claims represent the structure of the JWT token
type Claims struct {
	jwt.StandardClaims
	ID uint `gorm:"primaryKey;autoIncrement"`
}
