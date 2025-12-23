package repository

import (
	"errors"
	"testing"
	"time"

	"dotsat.work/internal/model"
	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

// setupTokenTestDB creates a test database connection and cleans up
func setupTokenTestDB(t *testing.T) *sqlx.DB {
	db, err := sqlx.Connect("pgx", "postgres://postgres:postgres@localhost:5432/dotsat?sslmode=disable")
	if err != nil {
		t.Fatalf("failed to connect to test database: %v", err)
	}

	// Clean up tables before each test
	_, err = db.Exec("TRUNCATE TABLE tokens, users, tenants CASCADE")
	if err != nil {
		t.Fatalf("failed to clean tables: %v", err)
	}

	return db
}

// createTestTenant creates a test tenant
func createTestTenant(t *testing.T, db *sqlx.DB) *model.Tenant {
	tenant := &model.Tenant{
		ID:        uuid.New(),
		Name:      "Test Tenant",
		Subdomain: "test",
		Status:    "active",
		Tier:      "free",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	query := `
		INSERT INTO tenants (id, name, subdomain, status, tier, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := db.Exec(query, tenant.ID, tenant.Name, tenant.Subdomain, tenant.Status, tenant.Tier, tenant.CreatedAt, tenant.UpdatedAt)
	if err != nil {
		t.Fatalf("failed to create test tenant: %v", err)
	}

	return tenant
}

// createTestUser creates a test user
func createTestUser(t *testing.T, db *sqlx.DB, tenantID uuid.UUID) *model.User {
	user := &model.User{
		ID:        uuid.New(),
		TenantID:  tenantID,
		Email:     "test@example.com",
		Role:      "user",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	query := `
		INSERT INTO users (id, tenant_id, email, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := db.Exec(query, user.ID, user.TenantID, user.Email, user.Role, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}

	return user
}

func TestTokenRepository_Create(t *testing.T) {
	db := setupTokenTestDB(t)
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("failed to close database: %v", err)
		}
	}()

	repo := NewTokenRepository(db)

	// Create a test tenant and user first
	tenant := createTestTenant(t, db)
	user := createTestUser(t, db, tenant.ID)

	token := &model.Token{
		UserID:    user.ID,
		Type:      model.TokenTypeMagicLink,
		Token:     "test-token-123",
		ExpiresAt: time.Now().Add(15 * time.Minute),
	}

	err := repo.Create(token)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if token.ID == uuid.Nil {
		t.Error("expected ID to be set")
	}

	if token.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be set")
	}
}

func TestTokenRepository_ConsumeToken(t *testing.T) {
	db := setupTokenTestDB(t)
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("failed to close database: %v", err)
		}
	}()

	repo := NewTokenRepository(db)

	// Create a test tenant and user first
	tenant := createTestTenant(t, db)
	user := createTestUser(t, db, tenant.ID)

	// Create a valid token
	token := &model.Token{
		UserID:    user.ID,
		Type:      model.TokenTypeMagicLink,
		Token:     "valid-token",
		ExpiresAt: time.Now().Add(15 * time.Minute),
	}
	err := repo.Create(token)
	if err != nil {
		t.Fatalf("failed to create token: %v", err)
	}

	// Consume the token
	consumed, err := repo.ConsumeToken(`valid-token`)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if consumed.UsedAt == nil {
		t.Error("expected UsedAt to be set")
	}

	// Try to consume again - should fail
	_, err = repo.ConsumeToken("valid-token")
	if !errors.Is(err, ErrTokenNotFound) {
		t.Errorf("expected ErrTokenNotFound, got %v", err)
	}
}

func TestTokenRepository_ConsumeToken_Expired(t *testing.T) {
	db := setupTokenTestDB(t)
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("failed to close database: %v", err)
		}
	}()

	repo := NewTokenRepository(db)

	// Create a test tenant and user first
	tenant := createTestTenant(t, db)
	user := createTestUser(t, db, tenant.ID)

	// Create an expired token
	token := &model.Token{
		UserID:    user.ID,
		Type:      model.TokenTypeMagicLink,
		Token:     "expired-token",
		ExpiresAt: time.Now().Add(-1 * time.Hour), // Expired 1 hour ago
	}
	err := repo.Create(token)
	if err != nil {
		t.Fatalf("failed to create token: %v", err)
	}

	// Try to consume - should fail
	_, err = repo.ConsumeToken("expired-token")
	if !errors.Is(err, ErrTokenNotFound) {
		t.Errorf("expected ErrTokenNotFound for expired token, got %v", err)
	}
}

func TestTokenRepository_DeleteByUserAndType(t *testing.T) {
	db := setupTokenTestDB(t)
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("failed to close database: %v", err)
		}
	}()

	repo := NewTokenRepository(db)

	// Create a test tenant and user first
	tenant := createTestTenant(t, db)
	user := createTestUser(t, db, tenant.ID)

	// Create multiple tokens
	token1 := &model.Token{
		UserID:    user.ID,
		Type:      model.TokenTypeMagicLink,
		Token:     "token-1",
		ExpiresAt: time.Now().Add(15 * time.Minute),
	}
	token2 := &model.Token{
		UserID:    user.ID,
		Type:      model.TokenTypeMagicLink,
		Token:     "token-2",
		ExpiresAt: time.Now().Add(15 * time.Minute),
	}
	token3 := &model.Token{
		UserID:    user.ID,
		Type:      model.TokenTypePasswordReset,
		Token:     "token-3",
		ExpiresAt: time.Now().Add(15 * time.Minute),
	}

	if err := repo.Create(token1); err != nil {
		t.Fatalf("failed to create token1: %v", err)
	}
	if err := repo.Create(token2); err != nil {
		t.Fatalf("failed to create token2: %v", err)
	}
	if err := repo.Create(token3); err != nil {
		t.Fatalf("failed to create token3: %v", err)
	}

	// Delete magic link tokens
	err := repo.DeleteByUserAndType(user.ID, model.TokenTypeMagicLink)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Magic link tokens should be gone
	_, err = repo.ConsumeToken("token-1")
	if !errors.Is(err, ErrTokenNotFound) {
		t.Error("expected token-1 to be deleted")
	}

	_, err = repo.ConsumeToken("token-2")
	if !errors.Is(err, ErrTokenNotFound) {
		t.Error("expected token-2 to be deleted")
	}

	// Password reset token should still exist
	_, err = repo.ConsumeToken("token-3")
	if err != nil {
		t.Error("expected token-3 to still exist")
	}
}
