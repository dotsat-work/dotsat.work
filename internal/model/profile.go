package model

import (
	"time"

	"github.com/google/uuid"
)

type Profile struct {
	ID        uuid.UUID `db:"id"`
	UserID    uuid.UUID `db:"user_id"`
	Name      string    `db:"name"`
	Bio       *string   `db:"bio"`
	Phone     *string   `db:"phone"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// HasName returns true if the profile has a name set
func (p *Profile) HasName() bool {
	return p.Name != ""
}
