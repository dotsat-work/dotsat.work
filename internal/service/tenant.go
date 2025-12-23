package service

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"dotsat.work/internal/model"
	"dotsat.work/internal/repository"
	"github.com/google/uuid"
)

var (
	ErrInvalidSubdomain  = errors.New("invalid subdomain: must be 1-63 characters, lowercase letters, numbers, and hyphens only")
	ErrInvalidTenantName = errors.New("invalid tenant name: must be 1-100 characters")
)

type TenantService struct {
	tenantRepository repository.TenantRepository
}

func NewTenantService(tenantRepository repository.TenantRepository) *TenantService {
	return &TenantService{
		tenantRepository: tenantRepository,
	}
}

// Create creates a new tenant with validation
func (s *TenantService) Create(name, subdomain string) (*model.Tenant, error) {
	// Validate name
	name = strings.TrimSpace(name)
	if len(name) < 1 || len(name) > 100 {
		return nil, ErrInvalidTenantName
	}

	// Validate subdomain
	subdomain = strings.ToLower(strings.TrimSpace(subdomain))
	if err := validateSubdomain(subdomain); err != nil {
		return nil, err
	}

	// Create tenant
	tenant := &model.Tenant{
		ID:        uuid.New(),
		Name:      name,
		Subdomain: subdomain,
		Status:    "active",
		Tier:      "free",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := s.tenantRepository.Create(tenant)
	if err != nil {
		if errors.Is(err, repository.ErrDuplicateSubdomain) {
			return nil, fmt.Errorf("subdomain %q is already taken", subdomain)
		}
		return nil, fmt.Errorf("failed to create tenant: %w", err)
	}

	return tenant, nil
}

// ByID retrieves a tenant by ID
func (s *TenantService) ByID(id uuid.UUID) (*model.Tenant, error) {
	tenant, err := s.tenantRepository.ByID(id)
	if err != nil {
		return nil, err
	}
	return tenant, nil
}

// BySubdomain retrieves a tenant by subdomain
func (s *TenantService) BySubdomain(subdomain string) (*model.Tenant, error) {
	subdomain = strings.ToLower(strings.TrimSpace(subdomain))
	tenant, err := s.tenantRepository.BySubdomain(subdomain)
	if err != nil {
		return nil, err
	}
	return tenant, nil
}

// Update updates a tenant
func (s *TenantService) Update(tenant *model.Tenant) error {
	// Validate name
	tenant.Name = strings.TrimSpace(tenant.Name)
	if len(tenant.Name) < 1 || len(tenant.Name) > 100 {
		return ErrInvalidTenantName
	}

	// Validate subdomain
	tenant.Subdomain = strings.ToLower(strings.TrimSpace(tenant.Subdomain))
	if err := validateSubdomain(tenant.Subdomain); err != nil {
		return err
	}

	tenant.UpdatedAt = time.Now()

	err := s.tenantRepository.Update(tenant)
	if err != nil {
		if errors.Is(err, repository.ErrDuplicateSubdomain) {
			return fmt.Errorf("subdomain %q is already taken", tenant.Subdomain)
		}
		return fmt.Errorf("failed to update tenant: %w", err)
	}

	return nil
}

// Delete deletes a tenant
func (s *TenantService) Delete(id uuid.UUID) error {
	err := s.tenantRepository.Delete(id)
	if err != nil {
		return fmt.Errorf("failed to delete tenant: %w", err)
	}
	return nil
}

// List returns all tenants
func (s *TenantService) List() ([]*model.Tenant, error) {
	return s.tenantRepository.List()
}

// validateSubdomain validates the subdomain format per RFC 1034
// Must be 1-63 characters, lowercase letters, numbers, and hyphens
// Cannot start or end with hyphen
func validateSubdomain(subdomain string) error {
	if len(subdomain) < 1 || len(subdomain) > 63 {
		return ErrInvalidSubdomain
	}

	if subdomain[0] == '-' || subdomain[len(subdomain)-1] == '-' {
		return ErrInvalidSubdomain
	}

	for _, char := range subdomain {
		if !((char >= 'a' && char <= 'z') || (char >= '0' && char <= '9') || char == '-') {
			return ErrInvalidSubdomain
		}
	}

	return nil
}
