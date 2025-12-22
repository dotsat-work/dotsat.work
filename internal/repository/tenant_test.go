package repository

import (
	"testing"
	"time"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"

	"dotsat.work/internal/model"
)

// setupTestDB creates a test database connection
// You'll need a test database running for this
func setupTestDB(t *testing.T) *sqlx.DB {
	// Use your test database connection
	db, err := sqlx.Connect("pgx", "postgres://postgres:postgres@localhost:5432/dotsat?sslmode=disable")
	if err != nil {
		t.Fatalf("failed to connect to test database: %v", err)
	}

	// Clean up the tenants table before each test
	_, err = db.Exec("TRUNCATE TABLE tenants CASCADE")
	if err != nil {
		t.Fatalf("failed to clean tenants table: %v", err)
	}

	return db
}

func TestTenantRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("failed to close database: %v", err)
		}
	}()

	repo := NewTenantRepository(db)

	tenant := &model.Tenant{
		ID:        uuid.New(),
		Name:      "Hewlett Packard",
		Subdomain: "hp",
		Status:    "active",
		Tier:      "premium",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := repo.Create(tenant)
	if err != nil {
		t.Fatalf("failed to create tenant: %v", err)
	}

	// Verify it was created
	found, err := repo.ByID(tenant.ID)
	if err != nil {
		t.Fatalf("failed to find tenant: %v", err)
	}

	if found.Name != tenant.Name {
		t.Errorf("expected name %q, got %q", tenant.Name, found.Name)
	}

	if found.Subdomain != tenant.Subdomain {
		t.Errorf("expected subdomain %q, got %q", tenant.Subdomain, found.Subdomain)
	}
}

func TestTenantRepository_Create_DuplicateSubdomain(t *testing.T) {
	db := setupTestDB(t)
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("failed to close database: %v", err)
		}
	}()

	repo := NewTenantRepository(db)

	tenant1 := &model.Tenant{
		ID:        uuid.New(),
		Name:      "Hewlett Packard",
		Subdomain: "hp",
		Status:    "active",
		Tier:      "premium",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := repo.Create(tenant1)
	if err != nil {
		t.Fatalf("failed to create first tenant: %v", err)
	}

	// Try to create another tenant with the same subdomain
	tenant2 := &model.Tenant{
		ID:        uuid.New(),
		Name:      "HP Inc",
		Subdomain: "hp", // Same subdomain!
		Status:    "active",
		Tier:      "standard",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = repo.Create(tenant2)
	if err != ErrDuplicateSubdomain {
		t.Errorf("expected ErrDuplicateSubdomain, got %v", err)
	}
}

func TestTenantRepository_BySubdomain(t *testing.T) {
	db := setupTestDB(t)
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("failed to close database: %v", err)
		}
	}()

	repo := NewTenantRepository(db)

	tenant := &model.Tenant{
		ID:        uuid.New(),
		Name:      "Honeywell",
		Subdomain: "honeywell",
		Status:    "active",
		Tier:      "standard",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := repo.Create(tenant)
	if err != nil {
		t.Fatalf("failed to create tenant: %v", err)
	}

	// Find by subdomain
	found, err := repo.BySubdomain("honeywell")
	if err != nil {
		t.Fatalf("failed to find tenant by subdomain: %v", err)
	}

	if found.ID != tenant.ID {
		t.Errorf("expected ID %v, got %v", tenant.ID, found.ID)
	}

	if found.Name != tenant.Name {
		t.Errorf("expected name %q, got %q", tenant.Name, found.Name)
	}
}

func TestTenantRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("failed to close database: %v", err)
		}
	}()

	repo := NewTenantRepository(db)

	tenant := &model.Tenant{
		ID:        uuid.New(),
		Name:      "Acme Corp",
		Subdomain: "acme",
		Status:    "active",
		Tier:      "standard",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := repo.Create(tenant)
	if err != nil {
		t.Fatalf("failed to create tenant: %v", err)
	}

	// Update tenant
	tenant.Name = "Acme Corporation"
	tenant.Tier = "premium"
	tenant.UpdatedAt = time.Now()

	err = repo.Update(tenant)
	if err != nil {
		t.Fatalf("failed to update tenant: %v", err)
	}

	// Verify update
	found, err := repo.ByID(tenant.ID)
	if err != nil {
		t.Fatalf("failed to find updated tenant: %v", err)
	}

	if found.Name != "Acme Corporation" {
		t.Errorf("expected name %q, got %q", "Acme Corporation", found.Name)
	}

	if found.Tier != "premium" {
		t.Errorf("expected tier %q, got %q", "premium", found.Tier)
	}
}
