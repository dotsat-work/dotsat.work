package model

import (
	"time"

	"github.com/google/uuid"
)

type Tenant struct {
	ID        uuid.UUID `db:"id"`
	Name      string    `db:"name"`
	Subdomain string    `db:"subdomain"`
	Status    string    `db:"status"` // active, suspended, inactive
	Tier      string    `db:"tier"`   // standard, premium, enterprise
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// IsActive returns true if the tenant is active
func (t *Tenant) IsActive() bool {
	return t.Status == "active"
}

// IsSuspended returns true if the tenant is suspended
func (t *Tenant) IsSuspended() bool {
	return t.Status == "suspended"
}
