package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Bayan2019/rbk-it-school-hw-4/internal/auth"
	"github.com/Bayan2019/rbk-it-school-hw-4/internal/domain"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

////// methods
////// methods
////// methods

func (r *UserRepository) Create(ctx context.Context, input domain.RegisterUserInput) (domain.User, error) {
	query := `
		INSERT INTO users (email, password_hash, first_name, last_name, is_active)
		VALUES (:email, :password_hash, :first_name, :last_name, :is_active)
		RETURNING id, email, password_hash, first_name, last_name, is_active, created_at, updated_at, deleted_at
	`

	// 6. Безопасность
	// - bcrypt для паролей
	// - не хранить пароли в plain text
	hashPassword, err := auth.HashPassword(input.Password)
	if err != nil {
		return domain.User{}, err
	}
	args := map[string]any{
		"email":         input.Email,
		"password_hash": hashPassword,
		"first_name":    input.FirstName,
		"last_name":     input.LastName,
		"is_active":     boolValue(input.IsActive, true),
	}

	rows, err := r.db.NamedQueryContext(ctx, query, args)
	if err != nil {
		if isUniqueViolation(err) {
			return domain.User{}, domain.ErrEmailAlreadyTaken
		}
		return domain.User{}, err
	}
	defer rows.Close()

	if rows.Next() {
		var user domain.User
		if err := rows.StructScan(&user); err != nil {
			return domain.User{}, err
		}
		return user, nil
	}

	return domain.User{}, errors.New("failed to create user")
}

func (r *UserRepository) List(ctx context.Context, filter domain.ListUsersFilter) ([]domain.User, error) {
	builder := strings.Builder{}
	builder.WriteString(`
		SELECT id, email, password_hash, first_name, last_name, is_active, created_at, updated_at, deleted_at
		FROM users
		WHERE 1=1
	`)

	args := map[string]any{
		"limit":  filter.Limit,
		"offset": filter.Offset,
	}

	if !filter.IncludeDeleted {
		fmt.Println("IncludeDeleted")
		builder.WriteString(" AND deleted_at IS NULL")
	}

	if filter.Query != "" {
		builder.WriteString(" AND (LOWER(email) LIKE :query OR LOWER(first_name) LIKE :query OR LOWER(last_name) LIKE :query)")
		args["query"] = "%" + filter.Query + "%"
	}

	builder.WriteString(" ORDER BY id LIMIT :limit OFFSET :offset")

	query, queryArgs, err := sqlx.Named(builder.String(), args)
	if err != nil {
		return nil, err
	}
	query = r.db.Rebind(query)

	var users []domain.User
	if err := r.db.SelectContext(ctx, &users, query, queryArgs...); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id int64, includeDeleted bool) (domain.User, error) {
	query := `
		SELECT id, email, password_hash, first_name, last_name, role, is_active, created_at, updated_at, deleted_at
		FROM users
		WHERE id = $1
	`
	if !includeDeleted {
		query += " AND deleted_at IS NULL"
	}

	var user domain.User
	if err := r.db.GetContext(ctx, &user, query, id); err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return domain.User{}, domain.ErrUserNotFound
		}
		return domain.User{}, err
	}

	return user, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string, includeDeleted bool) (domain.User, error) {
	query := `
		SELECT id, email, password_hash, first_name, last_name, role, is_active, created_at, updated_at, deleted_at
		FROM users
		WHERE email = $1
	`
	if !includeDeleted {
		query += " AND deleted_at IS NULL"
	}

	var user domain.User
	if err := r.db.GetContext(ctx, &user, query, email); err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return domain.User{}, domain.ErrUserNotFound
		}
		return domain.User{}, err
	}

	return user, nil
}

func (r *UserRepository) Update(ctx context.Context, id int64, input domain.UpdateUserInput) (domain.User, error) {
	query := `
		UPDATE users
		SET email = :email,
		    password_hash = :password_hash,
		    first_name = :first_name,
		    last_name = :last_name,
		    is_active = :is_active
		WHERE id = :id AND deleted_at IS NULL
		RETURNING id, email, password_hash, first_name, last_name, is_active, created_at, updated_at, deleted_at
	`

	args := map[string]any{
		"id":            id,
		"email":         input.Email,
		"password_hash": input.PasswordHash,
		"first_name":    input.FirstName,
		"last_name":     input.LastName,
		"is_active":     boolValue(input.IsActive, true),
	}

	rows, err := r.db.NamedQueryContext(ctx, query, args)
	if err != nil {
		if isUniqueViolation(err) {
			return domain.User{}, domain.ErrEmailAlreadyTaken
		}
		return domain.User{}, err
	}
	defer rows.Close()

	if rows.Next() {
		var user domain.User
		if err := rows.StructScan(&user); err != nil {
			return domain.User{}, err
		}
		return user, nil
	}

	return domain.User{}, domain.ErrUserNotFound
}

func (r *UserRepository) Delete(ctx context.Context, id int64) error {
	result, err := r.db.ExecContext(ctx, `
		UPDATE users
		SET deleted_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return domain.ErrUserNotFound
	}

	return nil
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

func boolValue(v *bool, fallback bool) bool {
	if v == nil {
		return fallback
	}
	return *v
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}
	return strings.Contains(strings.ToLower(err.Error()), "duplicate key")
}
