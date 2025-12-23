package validation

import (
	"testing"
)

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		name    string
		email   string
		wantErr bool
	}{
		// Valid emails
		{"valid simple", "test@example.com", false},
		{"valid with subdomain", "user@mail.example.com", false},
		{"valid with plus", "user+tag@example.com", false},
		{"valid with dots", "first.last@example.com", false},
		{"valid short", "a@b.c", false},
		{"valid with numbers", "user123@example.com", false},

		// Invalid emails
		{"empty", "", true},
		{"no @", "notanemail", true},
		{"no domain", "user@", true},
		{"no local part", "@example.com", true},
		{"double @", "user@@example.com", true},
		{"spaces", "user @example.com", true},
		{"too long", string(make([]byte, 255)) + "@example.com", true},

		// Note: "user@domain" (without TLD) is technically valid per RFC 5322
		// for local network emails, so net/mail accepts it
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateEmail(tt.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateEmail(%q) error = %v, wantErr %v", tt.email, err, tt.wantErr)
			}
		})
	}
}
