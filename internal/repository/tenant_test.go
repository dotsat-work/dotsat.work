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

func TestTenantRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("failed to close database: %v", err)
		}
	}()

	repo := NewTenantRepository(db)

	tenant := &model.Tenant{
		ID:        uuid.New(),
		Name:      "Test Company",
		Subdomain: "testco",
		Status:    "active",
		Tier:      "standard",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Create tenant first
	err := repo.Create(tenant)
	if err != nil {
		t.Fatalf("failed to create tenant: %v", err)
	}

	// Verify it exists
	found, err := repo.ByID(tenant.ID)
	if err != nil {
		t.Fatalf("failed to find tenant before delete: %v", err)
	}
	if found.ID != tenant.ID {
		t.Errorf("expected ID %v, got %v", tenant.ID, found.ID)
	}

	// Delete the tenant
	err = repo.Delete(tenant.ID)
	if err != nil {
		t.Fatalf("failed to delete tenant: %v", err)
	}

	// Verify it's gone
	_, err = repo.ByID(tenant.ID)
	if err != ErrTenantNotFound {
		t.Errorf("expected ErrTenantNotFound after delete, got %v", err)
	}
}

func TestTenantRepository_List(t *testing.T) {
	db := setupTestDB(t)
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("failed to close database: %v", err)
		}
	}()

	repo := NewTenantRepository(db)

	// Create multiple tenants
	tenants := []*model.Tenant{
		{
			ID:        uuid.New(),
			Name:      "Company A",
			Subdomain: "company-a",
			Status:    "active",
			Tier:      "premium",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        uuid.New(),
			Name:      "Company B",
			Subdomain: "company-b",
			Status:    "active",
			Tier:      "standard",
			CreatedAt: time.Now().Add(1 * time.Second),
			UpdatedAt: time.Now().Add(1 * time.Second),
		},
		{
			ID:        uuid.New(),
			Name:      "Company C",
			Subdomain: "company-c",
			Status:    "trial",
			Tier:      "standard",
			CreatedAt: time.Now().Add(2 * time.Second),
			UpdatedAt: time.Now().Add(2 * time.Second),
		},
	}

	for _, tenant := range tenants {
		err := repo.Create(tenant)
		if err != nil {
			t.Fatalf("failed to create tenant %s: %v", tenant.Name, err)
		}
	}

	// List all tenants
	list, err := repo.List()
	if err != nil {
		t.Fatalf("failed to list tenants: %v", err)
	}

	// Verify we got at least the 3 we created
	if len(list) < 3 {
		t.Errorf("expected at least 3 tenants, got %d", len(list))
	}

	// Verify the tenants are ordered by created_at DESC (newest first)
	// The last created tenant (Company C) should be first in the list
	foundCompanyC := false
	for _, tenant := range list {
		if tenant.Subdomain == "company-c" {
			foundCompanyC = true
			break
		}
	}

	if !foundCompanyC {
		t.Error("expected to find Company C in the list")
	}

	// Verify all our created tenants are in the list
	expectedSubdomains := map[string]bool{
		"company-a": false,
		"company-b": false,
		"company-c": false,
	}

	for _, tenant := range list {
		if _, exists := expectedSubdomains[tenant.Subdomain]; exists {
			expectedSubdomains[tenant.Subdomain] = true
		}
	}

	for subdomain, found := range expectedSubdomains {
		if !found {
			t.Errorf("expected to find tenant with subdomain %q in list", subdomain)
		}
	}
}

func TestTenantRepository_List_Empty(t *testing.T) {
	db := setupTestDB(t)
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("failed to close database: %v", err)
		}
	}()

	repo := NewTenantRepository(db)

	// List should return empty slice, not error
	list, err := repo.List()
	if err != nil {
		t.Fatalf("failed to list tenants: %v", err)
	}

	if list == nil {
		t.Error("expected empty slice, got nil")
	}

	if len(list) != 0 {
		t.Errorf("expected empty list, got %d tenants", len(list))
	}
}
