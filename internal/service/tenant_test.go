package service

import (
	"testing"
)

func TestValidateSubdomain(t *testing.T) {
	tests := []struct {
		name      string
		subdomain string
		wantErr   bool
	}{
		// Valid subdomains
		{"single char", "a", false},
		{"two chars", "hp", false},
		{"three chars", "ibm", false},
		{"with numbers", "tech123", false},
		{"with hyphens", "my-company", false},
		{"max length", "a12345678901234567890123456789012345678901234567890123456789012", false}, // 63 chars
		{"lowercase", "acme", false},

		// Invalid subdomains
		{"empty", "", true},
		{"too long", "a123456789012345678901234567890123456789012345678901234567890123", true}, // 64 chars
		{"uppercase", "Acme", true},
		{"starts with hyphen", "-acme", true},
		{"ends with hyphen", "acme-", true},
		{"underscore", "acme_corp", true},
		{"space", "acme corp", true},
		{"special chars", "acme@corp", true},
		{"dot", "acme.corp", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateSubdomain(tt.subdomain)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateSubdomain(%q) error = %v, wantErr %v", tt.subdomain, err, tt.wantErr)
			}
		})
	}
}
