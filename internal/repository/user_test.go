package repository

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"

	"dotsat.work/internal/model"
)

// setupUserTestDB creates a test database connection and cleans up users and tenants
func setupUserTestDB(t *testing.T) (*sqlx.DB, uuid.UUID) {
	db, err := sqlx.Connect("pgx", "postgres://postgres:postgres@localhost:5432/dotsat?sslmode=disable")
	if err != nil {
		t.Fatalf("failed to connect to test database: %v", err)
	}

	// Clean up users and tenants tables before each test
	_, err = db.Exec("TRUNCATE TABLE users CASCADE")
	if err != nil {
		t.Fatalf("failed to clean users table: %v", err)
	}

	_, err = db.Exec("TRUNCATE TABLE tenants CASCADE")
	if err != nil {
		t.Fatalf("failed to clean tenants table: %v", err)
	}

	// Create a test tenant for users to belong to
	tenantID := uuid.New()
	_, err = db.Exec(
		"INSERT INTO tenants (id, name, subdomain, status, tier, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		tenantID, "Test Tenant", "test-tenant", "active", "standard", time.Now(), time.Now(),
	)
	if err != nil {
		t.Fatalf("failed to create test tenant: %v", err)
	}

	return db, tenantID
}

func TestUserRepository_Create(t *testing.T) {
	db, tenantID := setupUserTestDB(t)
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("failed to close database: %v", err)
		}
	}()

	repo := NewUserRepository(db)

	passwordHash := "hashed_password_123"
	user := &model.User{
		ID:           uuid.New(),
		TenantID:     tenantID,
		Email:        "john.doe@example.com",
		PasswordHash: &passwordHash,
		Role:         "user",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	err := repo.Create(user)
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	// Verify it was created
	found, err := repo.ByID(user.ID)
	if err != nil {
		t.Fatalf("failed to find user: %v", err)
	}

	if found.Email != user.Email {
		t.Errorf("expected email %q, got %q", user.Email, found.Email)
	}

	if found.TenantID != user.TenantID {
		t.Errorf("expected tenant_id %v, got %v", user.TenantID, found.TenantID)
	}

	if found.Role != user.Role {
		t.Errorf("expected role %q, got %q", user.Role, found.Role)
	}
}

func TestUserRepository_Create_DuplicateEmail(t *testing.T) {
	db, tenantID := setupUserTestDB(t)
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("failed to close database: %v", err)
		}
	}()

	repo := NewUserRepository(db)

	passwordHash := "hashed_password_123"
	user1 := &model.User{
		ID:           uuid.New(),
		TenantID:     tenantID,
		Email:        "duplicate@example.com",
		PasswordHash: &passwordHash,
		Role:         "user",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	err := repo.Create(user1)
	if err != nil {
		t.Fatalf("failed to create first user: %v", err)
	}

	// Try to create another user with same email in same tenant
	user2 := &model.User{
		ID:           uuid.New(),
		TenantID:     tenantID,                // Same tenant!
		Email:        "duplicate@example.com", // Same email!
		PasswordHash: &passwordHash,
		Role:         "admin",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	err = repo.Create(user2)
	if !errors.Is(err, ErrDuplicateEmail) {
		t.Errorf("expected ErrDuplicateEmail, got %v", err)
	}
}

func TestUserRepository_Create_SameEmailDifferentTenant(t *testing.T) {
	db, tenantID1 := setupUserTestDB(t)
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("failed to close database: %v", err)
		}
	}()

	// Create a second tenant
	tenantID2 := uuid.New()
	_, err := db.Exec(
		"INSERT INTO tenants (id, name, subdomain, status, tier, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		tenantID2, "Second Tenant", "second-tenant", "active", "standard", time.Now(), time.Now(),
	)
	if err != nil {
		t.Fatalf("failed to create second tenant: %v", err)
	}

	repo := NewUserRepository(db)

	passwordHash := "hashed_password_123"

	// Create user in first tenant
	user1 := &model.User{
		ID:           uuid.New(),
		TenantID:     tenantID1,
		Email:        "same@example.com",
		PasswordHash: &passwordHash,
		Role:         "user",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	err = repo.Create(user1)
	if err != nil {
		t.Fatalf("failed to create user in first tenant: %v", err)
	}

	// Create user with same email in second tenant - should succeed!
	user2 := &model.User{
		ID:           uuid.New(),
		TenantID:     tenantID2,          // Different tenant!
		Email:        "same@example.com", // Same email!
		PasswordHash: &passwordHash,
		Role:         "admin",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	err = repo.Create(user2)
	if err != nil {
		t.Fatalf("expected success for same email in different tenant, got error: %v", err)
	}

	// Verify both users exist
	found1, err := repo.ByID(user1.ID)
	if err != nil {
		t.Fatalf("failed to find user1: %v", err)
	}

	found2, err := repo.ByID(user2.ID)
	if err != nil {
		t.Fatalf("failed to find user2: %v", err)
	}

	if found1.Email != found2.Email {
		t.Error("expected both users to have same email")
	}

	if found1.TenantID == found2.TenantID {
		t.Error("expected users to belong to different tenants")
	}
}

func TestUserRepository_ByEmail(t *testing.T) {
	db, tenantID := setupUserTestDB(t)
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("failed to close database: %v", err)
		}
	}()

	repo := NewUserRepository(db)

	passwordHash := "hashed_password_123"
	user := &model.User{
		ID:           uuid.New(),
		TenantID:     tenantID,
		Email:        "findme@example.com",
		PasswordHash: &passwordHash,
		Role:         "user",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	err := repo.Create(user)
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	// Find by email
	found, err := repo.ByEmail("findme@example.com")
	if err != nil {
		t.Fatalf("failed to find user by email: %v", err)
	}

	if found.ID != user.ID {
		t.Errorf("expected ID %v, got %v", user.ID, found.ID)
	}

	if found.Email != user.Email {
		t.Errorf("expected email %q, got %q", user.Email, found.Email)
	}
}

func TestUserRepository_ByTenantID(t *testing.T) {
	db, tenantID := setupUserTestDB(t)
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("failed to close database: %v", err)
		}
	}()

	repo := NewUserRepository(db)

	passwordHash := "hashed_password_123"

	// Create multiple users in the same tenant
	users := []*model.User{
		{
			ID:           uuid.New(),
			TenantID:     tenantID,
			Email:        "user1@example.com",
			PasswordHash: &passwordHash,
			Role:         "admin",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
		{
			ID:           uuid.New(),
			TenantID:     tenantID,
			Email:        "user2@example.com",
			PasswordHash: &passwordHash,
			Role:         "user",
			CreatedAt:    time.Now().Add(1 * time.Second),
			UpdatedAt:    time.Now().Add(1 * time.Second),
		},
		{
			ID:           uuid.New(),
			TenantID:     tenantID,
			Email:        "user3@example.com",
			PasswordHash: &passwordHash,
			Role:         "viewer",
			CreatedAt:    time.Now().Add(2 * time.Second),
			UpdatedAt:    time.Now().Add(2 * time.Second),
		},
	}

	for _, user := range users {
		err := repo.Create(user)
		if err != nil {
			t.Fatalf("failed to create user %s: %v", user.Email, err)
		}
	}

	// Get all users for this tenant
	foundUsers, err := repo.ByTenantID(tenantID)
	if err != nil {
		t.Fatalf("failed to get users by tenant: %v", err)
	}

	if len(foundUsers) != 3 {
		t.Errorf("expected 3 users, got %d", len(foundUsers))
	}

	// Verify all emails are present
	expectedEmails := map[string]bool{
		"user1@example.com": false,
		"user2@example.com": false,
		"user3@example.com": false,
	}

	for _, user := range foundUsers {
		if _, exists := expectedEmails[user.Email]; exists {
			expectedEmails[user.Email] = true
		}
	}

	for email, found := range expectedEmails {
		if !found {
			t.Errorf("expected to find user with email %q", email)
		}
	}
}

func TestUserRepository_Update(t *testing.T) {
	db, tenantID := setupUserTestDB(t)
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("failed to close database: %v", err)
		}
	}()

	repo := NewUserRepository(db)

	passwordHash := "hashed_password_123"
	user := &model.User{
		ID:           uuid.New(),
		TenantID:     tenantID,
		Email:        "update@example.com",
		PasswordHash: &passwordHash,
		Role:         "user",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	err := repo.Create(user)
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	// Update user
	newPasswordHash := "new_hashed_password_456"
	user.Email = "updated@example.com"
	user.PasswordHash = &newPasswordHash
	user.Role = "admin"
	user.UpdatedAt = time.Now()

	err = repo.Update(user)
	if err != nil {
		t.Fatalf("failed to update user: %v", err)
	}

	// Verify update
	found, err := repo.ByID(user.ID)
	if err != nil {
		t.Fatalf("failed to find updated user: %v", err)
	}

	if found.Email != "updated@example.com" {
		t.Errorf("expected email %q, got %q", "updated@example.com", found.Email)
	}

	if found.Role != "admin" {
		t.Errorf("expected role %q, got %q", "admin", found.Role)
	}

	if found.PasswordHash == nil || *found.PasswordHash != newPasswordHash {
		t.Error("expected password hash to be updated")
	}
}

func TestUserRepository_Update_NotFound(t *testing.T) {
	db, tenantID := setupUserTestDB(t)
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("failed to close database: %v", err)
		}
	}()

	repo := NewUserRepository(db)

	// Try to update a user that doesn't exist
	passwordHash := "hashed_password_123"
	nonExistentUser := &model.User{
		ID:           uuid.New(), // This ID doesn't exist in DB
		TenantID:     tenantID,
		Email:        "nonexistent@example.com",
		PasswordHash: &passwordHash,
		Role:         "admin",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	err := repo.Update(nonExistentUser)
	if !errors.Is(err, ErrUserNotFound) {
		t.Errorf("expected ErrUserNotFound when updating non-existent user, got %v", err)
	}
}

func TestUserRepository_Delete(t *testing.T) {
	db, tenantID := setupUserTestDB(t)
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("failed to close database: %v", err)
		}
	}()

	repo := NewUserRepository(db)

	passwordHash := "hashed_password_123"
	user := &model.User{
		ID:           uuid.New(),
		TenantID:     tenantID,
		Email:        "delete@example.com",
		PasswordHash: &passwordHash,
		Role:         "user",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Create user first
	err := repo.Create(user)
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	// Verify it exists
	found, err := repo.ByID(user.ID)
	if err != nil {
		t.Fatalf("failed to find user before delete: %v", err)
	}
	if found.ID != user.ID {
		t.Errorf("expected ID %v, got %v", user.ID, found.ID)
	}

	// Delete the user
	err = repo.Delete(user.ID)
	if err != nil {
		t.Fatalf("failed to delete user: %v", err)
	}

	// Verify it's gone
	_, err = repo.ByID(user.ID)
	if !errors.Is(err, ErrUserNotFound) {
		t.Errorf("expected ErrUserNotFound after delete, got %v", err)
	}
}

func TestUserRepository_Delete_NotFound(t *testing.T) {
	db, _ := setupUserTestDB(t)
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("failed to close database: %v", err)
		}
	}()

	repo := NewUserRepository(db)

	// Try to delete non-existent user
	nonExistentID := uuid.New()
	err := repo.Delete(nonExistentID)
	if !errors.Is(err, ErrUserNotFound) {
		t.Errorf("expected ErrUserNotFound when deleting non-existent user, got %v", err)
	}
}
