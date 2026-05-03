package domain

import (
	"errors"
	"strings"
	"time"
)

var (
	ErrUserNotFound     = errors.New("user not found")
	ErrInvalidUserID    = errors.New("invalid user id")
	ErrInvalidUserInput = errors.New("invalid user input")
	// 1. Аутентификация
	// - email должен быть уникальным
	ErrEmailAlreadyTaken = errors.New("email already exists")
)

type Roles string

const (
	RolesAdmin Roles = "admin"
	RolesUser  Roles = "user"
	// RolesManager      Roles = "manager"
	// RolesPsychologist Roles = "psychologist"
	// RolesClient       Roles = "client"
	// RolesValue4       Roles = ""
)

type UserContext struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
	Role  Roles  `json:"role"`
}

type User struct {
	ID           int64      `db:"id" json:"id"`
	Email        string     `db:"email" json:"email"`
	PasswordHash string     `db:"password_hash" json:"-"`
	FirstName    string     `db:"first_name" json:"first_name"`
	LastName     string     `db:"last_name" json:"last_name"`
	Role         Roles      `db:"role" json:"role"`
	IsActive     bool       `db:"is_active" json:"is_active"`
	CreatedAt    time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt    *time.Time `db:"deleted_at" json:"deleted_at,omitempty"`
}

type RegisterUserInput struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	IsActive  *bool  `json:"is_active,omitempty"`
}

type UpdateUserInput struct {
	Email        string `json:"email"`
	PasswordHash string `json:"password_hash"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	IsActive     *bool  `json:"is_active,omitempty"`
}

type ListUsersFilter struct {
	Limit          int    `json:"limit"`
	Offset         int    `json:"offset"`
	Query          string `json:"query"`
	IncludeDeleted bool   `json:"include_deleted"`
}

////// accommodating functions
////// accommodating functions
////// accommodating functions
////// accommodating functions
////// accommodating functions
////// accommodating functions
////// accommodating functions
////// accommodating functions
////// accommodating functions
////// accommodating functions
////// accommodating functions

func (in *RegisterUserInput) NormalizeAndValidate() error {
	in.Email = strings.TrimSpace(strings.ToLower(in.Email))
	in.Password = strings.TrimSpace(in.Password)
	in.FirstName = strings.TrimSpace(in.FirstName)
	in.LastName = strings.TrimSpace(in.LastName)

	if in.Email == "" || !strings.Contains(in.Email, "@") {
		return ErrInvalidUserInput
	}
	if in.Password == "" || in.FirstName == "" || in.LastName == "" {
		return ErrInvalidUserInput
	}
	return nil
}

func (in *UpdateUserInput) NormalizeAndValidate() error {
	in.Email = strings.TrimSpace(strings.ToLower(in.Email))
	in.PasswordHash = strings.TrimSpace(in.PasswordHash)
	in.FirstName = strings.TrimSpace(in.FirstName)
	in.LastName = strings.TrimSpace(in.LastName)

	if in.Email == "" || !strings.Contains(in.Email, "@") {
		return ErrInvalidUserInput
	}
	if in.PasswordHash == "" || in.FirstName == "" || in.LastName == "" {
		return ErrInvalidUserInput
	}
	return nil
}

func (f *ListUsersFilter) Normalize() {
	if f.Limit <= 0 || f.Limit > 100 {
		f.Limit = 20
	}
	if f.Offset < 0 {
		f.Offset = 0
	}
	f.Query = strings.TrimSpace(
		strings.ToLower(f.Query))
}
