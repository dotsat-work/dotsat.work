package service

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"dotsat.work/internal/model"
	"dotsat.work/internal/repository"
	"dotsat.work/internal/validation"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidRole            = errors.New("invalid role: must be 'admin', 'user', or 'viewer'")
	ErrInvalidCurrentPassword = errors.New("current password is incorrect")
)

type UserService struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepository repository.UserRepository) *UserService {
	return &UserService{
		userRepository: userRepository,
	}
}

// Create creates a new user with password hashing
func (s *UserService) Create(tenantID uuid.UUID, email, password, role string) (*model.User, error) {
	// Validate email
	email = strings.ToLower(strings.TrimSpace(email))
	if err := validation.ValidateEmail(email); err != nil {
		return nil, err
	}

	// Validate role
	if !isValidRole(role) {
		return nil, ErrInvalidRole
	}

	// Hash password (if provided)
	var passwordHash *string
	if password != "" {
		if err := validation.ValidatePassword(password); err != nil {
			return nil, err
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return nil, fmt.Errorf("failed to hash password: %w", err)
		}
		hashStr := string(hash)
		passwordHash = &hashStr
	}

	// Create user
	user := &model.User{
		ID:           uuid.New(),
		TenantID:     tenantID,
		Email:        email,
		PasswordHash: passwordHash,
		Role:         role,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	err := s.userRepository.Create(user)
	if err != nil {
		if errors.Is(err, repository.ErrDuplicateEmail) {
			return nil, fmt.Errorf("email %q is already registered in this organization", email)
		}
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

// ByID retrieves a user by ID
func (s *UserService) ByID(id uuid.UUID) (*model.User, error) {
	return s.userRepository.ByID(id)
}

// ByEmail retrieves a user by email
func (s *UserService) ByEmail(email string) (*model.User, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	return s.userRepository.ByEmail(email)
}

// ByTenantID retrieves all users for a tenant
func (s *UserService) ByTenantID(tenantID uuid.UUID) ([]*model.User, error) {
	return s.userRepository.ByTenantID(tenantID)
}

// Update updates a user
func (s *UserService) Update(user *model.User) error {
	// Validate email
	user.Email = strings.ToLower(strings.TrimSpace(user.Email))
	if err := validation.ValidateEmail(user.Email); err != nil {
		return err
	}

	// Validate role
	if !isValidRole(user.Role) {
		return ErrInvalidRole
	}

	user.UpdatedAt = time.Now()

	err := s.userRepository.Update(user)
	if err != nil {
		if errors.Is(err, repository.ErrDuplicateEmail) {
			return fmt.Errorf("email %q is already registered in this organization", user.Email)
		}
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

// UpdatePassword updates a user's password
func (s *UserService) UpdatePassword(userID uuid.UUID, currentPassword, newPassword string) error {
	user, err := s.userRepository.ByID(userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Check if the user has a password
	if user.PasswordHash == nil {
		return fmt.Errorf("passwordless accounts cannot update password")
	}

	// Verify the current password
	err = bcrypt.CompareHashAndPassword([]byte(*user.PasswordHash), []byte(currentPassword))
	if err != nil {
		return ErrInvalidCurrentPassword
	}

	// Validate new password
	if err := validation.ValidatePassword(newPassword); err != nil {
		return err
	}

	// Hash new password
	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	hashStr := string(hash)
	user.PasswordHash = &hashStr

	err = s.userRepository.Update(user)
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	return nil
}

// Delete deletes a user
func (s *UserService) Delete(id uuid.UUID) error {
	err := s.userRepository.Delete(id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

// isValidRole checks if role is valid
func isValidRole(role string) bool {
	switch role {
	case "admin", "user", "viewer":
		return true
	default:
		return false
	}
}
