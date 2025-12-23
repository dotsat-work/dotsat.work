package validation

import (
	"strings"
	"testing"
)

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		// Valid passwords - clean, no common patterns
		{"valid strong passphrase", "correct-horse-battery-staple", false},
		{"valid long password", "MySecretPassword2024!", false},
		{"valid with spaces", "my super secret phrase 2024", false},
		{"valid 12 chars exactly", "abcDEF123!@#", false},
		{"valid unicode", "пароль-секретный-2024", false},
		{"valid numbers only", "123409876543", false},

		// Too short
		{"too short 11 chars", "Short123!@#", true},
		{"too short 8 chars", "Pass123!", true},
		{"too short 1 char", "a", true},
		{"empty", "", true},

		// Too long (bcrypt limit)
		{"too long 73 chars", strings.Repeat("a", 73), true},
		{"too long 100 chars", strings.Repeat("x", 100), true},
		{"exactly 72 chars mixed", "abcdefghij" + strings.Repeat("klmnopqrst", 6) + "uv", false}, // 72 chars, no repetition

		// Exact match - common passwords (blocked)
		{"exact password", "password", true},
		{"exact passw0rd", "passw0rd", true},
		{"exact admin", "admin", true},
		{"exact letmein", "letmein", true},
		{"exact qwerty", "qwerty", true},
		{"exact 123456", "123456", true},
		{"exact 12345678", "12345678", true},
		{"exact 123456789", "123456789", true},
		{"exact 1234567890", "1234567890", true},
		{"exact qwertyuiop", "qwertyuiop", true},
		{"exact asdfghjkl", "asdfghjkl", true},
		{"exact welcome", "welcome", true},
		{"exact monkey", "monkey", true},

		// Common password + numbers (blocked)
		{"password123", "password123", true},
		{"password1234567", "password1234567", true},
		{"admin123", "admin123", true},
		{"admin456789", "admin456789", true},
		{"letmein1", "letmein1", true},
		{"letmein667845", "letmein667845", true},
		{"qwerty123", "qwerty123", true},
		{"qwerty999", "qwerty999", true},
		{"welcome123", "welcome123", true},
		{"monkey456", "monkey456", true},
		{"1234567890123", "1234567890123", true}, // 1234567890 + 123

		// Numbers + common password (blocked)
		{"123password", "123password", true},
		{"456admin", "456admin", true},
		{"999letmein", "999letmein", true},
		{"123qwerty", "123qwerty", true},
		{"789welcome", "789welcome", true},
		{"000monkey", "000monkey", true},

		// Common password + special chars (blocked if the pattern matches)
		{"password!", "password!", true},
		{"password!@#", "password!@#", true},
		{"admin!!", "admin!!", true},
		{"letmein@#$", "letmein@#$", true},

		// Special chars + common password (blocked if the pattern matches)
		{"!password", "!password", true},
		{"@admin", "@admin", true},
		{"#letmein", "#letmein", true},

		// Valid - common words in strong passwords
		{"password in middle", "my-password-is-secret-2024", false},
		{"admin in middle", "admin-account-management-system", false},
		{"letmein in middle", "letmein-but-make-it-secure", false},
		{"qwerty in middle", "my-qwerty-keyboard-layout", false},
		{"123456 in middle", "room-123456-floor-9", false},

		// Valid - contains common word but not exact pattern
		{"pass not password", "pass-phrase-2024", false},
		{"adm not admin", "admiration-society", false},
		{"let not letmein", "letter-to-santa", false},
		{"123 not 123456", "room-123-floor-4", false},

		// Excessive repetition (6+ same chars) - blocked
		{"6 repeated a", "aaaaaa-rest-of-password", true},
		{"7 repeated 1", "1111111-rest", true},
		{"10 repeated x", "xxxxxxxxxx-rest", true},
		{"all repeated", "aaaaaaaaaaaa", true},

		// Acceptable repetition (5 or fewer) - allowed
		{"5 repeated a ok", "aaaaa-rest-of-password", false},
		{"3 repeated ok", "aaa-bbb-ccc-ddd-eee", false},
		{"no repetition", "abcdefghijkl", false},

		// Case-insensitive - common passwords
		{"uppercase PASSWORD", "PASSWORD", true},
		{"uppercase PASSWORD123", "PASSWORD123", true},
		{"mixed case PaSsWoRd", "PaSsWoRd", true},
		{"mixed case PaSsWoRd123", "PaSsWoRd123", true},
		{"uppercase ADMIN456", "ADMIN456", true},
		{"mixed case QwErTy789", "QwErTy789", true},

		// Case-insensitive - but valid when in the middle
		{"uppercase in middle", "my-PASSWORD-is-secure", false},
		{"mixed case in middle", "QwErTy-keyboard-layout", false},

		// Edge cases
		{"passw0rd variation", "passw0rd123", true}, // Leetspeak caught
		{"valid similar to weak", "passwords-are-important", false},
		{"valid with numbers", "my-secure-pass-2024", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePassword(tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePassword(%q) error = %v, wantErr %v", tt.password, err, tt.wantErr)
			}
		})
	}
}

func TestHasExcessiveRepetition(t *testing.T) {
	tests := []struct {
		name     string
		password string
		want     bool
	}{
		// Excessive repetition (6+ consecutive chars)
		{"exactly 6 repeated a", "aaaaaa", true},
		{"7 repeated 1", "1111111", true},
		{"10 repeated x", "xxxxxxxxxx", true},
		{"6 in middle", "abc-aaaaaa-def", true},
		{"6 at start", "aaaaaabcdefgh", true},
		{"6 at end", "abcdefaaaaaa", true},

		// Acceptable repetition (5 or fewer)
		{"exactly 5 repeated a", "aaaaa", false},
		{"3 repeated", "aaa", false},
		{"2 repeated", "aa", false},
		{"1 char", "a", false},
		{"no repetition", "abcdefghijkl", false},
		{"multiple groups of 3", "aaa-bbb-ccc-ddd", false},
		{"alternating", "ababababab", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := hasExcessiveRepetition(tt.password)
			if got != tt.want {
				t.Errorf("hasExcessiveRepetition(%q) = %v, want %v", tt.password, got, tt.want)
			}
		})
	}
}
