package model

import (
	"time"

	"github.com/google/uuid"
)

type Token struct {
	ID        uuid.UUID  `db:"id"`
	UserID    uuid.UUID  `db:"user_id"`
	Type      string     `db:"type"` // "email_verify", "password_reset", "magic_link", "email_change"
	Token     string     `db:"token"`
	ExpiresAt time.Time  `db:"expires_at"`
	UsedAt    *time.Time `db:"used_at"`
	CreatedAt time.Time  `db:"created_at"`
}

const (
	TokenTypeEmailVerify   = "email_verify"
	TokenTypePasswordReset = "password_reset"
	TokenTypeEmailChange   = "email_change"
	TokenTypeMagicLink     = "magic_link"
)

// IsExpired returns true if the token has expired
func (t *Token) IsExpired() bool {
	return time.Now().After(t.ExpiresAt)
}

// IsUsed returns true if the token has been used
func (t *Token) IsUsed() bool {
	return t.UsedAt != nil
}

// IsValid returns true if the token is not expired and not used
func (t *Token) IsValid() bool {
	return !t.IsExpired() && !t.IsUsed()
}
