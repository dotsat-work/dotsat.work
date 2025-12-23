package repository

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"dotsat.work/internal/model"
)

var (
	ErrUserNotFound   = errors.New("user not found")
	ErrDuplicateEmail = errors.New("email already exists in this tenant")
)

type UserRepository interface {
	Create(user *model.User) error
	ByID(id uuid.UUID) (*model.User, error)
	ByEmail(email string) (*model.User, error)
	ByTenantID(tenantID uuid.UUID) ([]*model.User, error)
	Update(user *model.User) error
	Delete(id uuid.UUID) error
}

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *model.User) error {
	query := `
		INSERT INTO users (id, tenant_id, email, password_hash, role, pending_email, email_verified_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err := r.db.Exec(
		query,
		user.ID,
		user.TenantID,
		user.Email,
		user.PasswordHash,
		user.Role,
		user.PendingEmail,
		user.EmailVerifiedAt,
		user.CreatedAt,
		user.UpdatedAt,
	)
	if err != nil {
		// Check for unique constraint violation
		if strings.Contains(err.Error(), "duplicate key value") {
			return ErrDuplicateEmail
		}
		return err
	}

	return nil
}

func (r *userRepository) ByID(id uuid.UUID) (*model.User, error) {
	user := &model.User{}
	query := `SELECT * FROM users WHERE id = $1`

	err := r.db.Get(user, query, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrUserNotFound
	}

	return user, err
}

func (r *userRepository) ByEmail(email string) (*model.User, error) {
	user := &model.User{}
	query := `SELECT * FROM users WHERE email = $1`

	err := r.db.Get(user, query, email)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrUserNotFound
	}

	return user, err
}

func (r *userRepository) ByTenantID(tenantID uuid.UUID) ([]*model.User, error) {
	var users []*model.User
	query := `SELECT * FROM users WHERE tenant_id = $1 ORDER BY created_at DESC`

	err := r.db.Select(&users, query, tenantID)
	return users, err
}

func (r *userRepository) Update(user *model.User) error {
	query := `
		UPDATE users
		SET email = $1, password_hash = $2, role = $3, pending_email = $4, email_verified_at = $5, updated_at = $6
		WHERE id = $7
	`

	result, err := r.db.Exec(
		query,
		user.Email,
		user.PasswordHash,
		user.Role,
		user.PendingEmail,
		user.EmailVerifiedAt,
		user.UpdatedAt,
		user.ID,
	)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrUserNotFound
	}

	return nil
}

func (r *userRepository) Delete(id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrUserNotFound
	}

	return nil
}
