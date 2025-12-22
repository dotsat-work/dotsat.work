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
	ErrTenantNotFound     = errors.New("tenant not found")
	ErrDuplicateSubdomain = errors.New("subdomain already exists")
)

type TenantRepository interface {
	Create(tenant *model.Tenant) error
	ByID(id uuid.UUID) (*model.Tenant, error)
	BySubdomain(subdomain string) (*model.Tenant, error)
	Update(tenant *model.Tenant) error
	Delete(id uuid.UUID) error
	List() ([]*model.Tenant, error)
}

type tenantRepository struct {
	db *sqlx.DB
}

func NewTenantRepository(db *sqlx.DB) TenantRepository {
	return &tenantRepository{db: db}
}

func (r *tenantRepository) Create(tenant *model.Tenant) error {
	query := `
		INSERT INTO tenants (id, name, subdomain, status, tier, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.db.Exec(
		query,
		tenant.ID,
		tenant.Name,
		tenant.Subdomain,
		tenant.Status,
		tenant.Tier,
		tenant.CreatedAt,
		tenant.UpdatedAt,
	)
	if err != nil {
		// Check for unique constraint violation
		if strings.Contains(err.Error(), "duplicate key value") {
			return ErrDuplicateSubdomain
		}
		return err
	}

	return nil
}

func (r *tenantRepository) ByID(id uuid.UUID) (*model.Tenant, error) {
	tenant := &model.Tenant{}
	query := `SELECT * FROM tenants WHERE id = $1`

	err := r.db.Get(tenant, query, id)
	if err == sql.ErrNoRows {
		return nil, ErrTenantNotFound
	}

	return tenant, err
}

func (r *tenantRepository) BySubdomain(subdomain string) (*model.Tenant, error) {
	tenant := &model.Tenant{}
	query := `SELECT * FROM tenants WHERE subdomain = $1`

	err := r.db.Get(tenant, query, subdomain)
	if err == sql.ErrNoRows {
		return nil, ErrTenantNotFound
	}

	return tenant, err
}

func (r *tenantRepository) Update(tenant *model.Tenant) error {
	query := `
		UPDATE tenants
		SET name = $1, subdomain = $2, status = $3, tier = $4, updated_at = $5
		WHERE id = $6
	`

	_, err := r.db.Exec(
		query,
		tenant.Name,
		tenant.Subdomain,
		tenant.Status,
		tenant.Tier,
		tenant.UpdatedAt,
		tenant.ID,
	)
	return err
}

func (r *tenantRepository) Delete(id uuid.UUID) error {
	query := `DELETE FROM tenants WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *tenantRepository) List() ([]*model.Tenant, error) {
	var tenants []*model.Tenant
	query := `SELECT * FROM tenants ORDER BY created_at DESC`

	err := r.db.Select(&tenants, query)
	return tenants, err
}
