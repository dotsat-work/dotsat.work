package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHomeHandler_ServeHTTP(t *testing.T) {
	tests := []struct {
		name             string
		expectedStatus   int
		expectedContent  string
		expectedContains []string
	}{
		{
			name:             "returns welcome page",
			expectedStatus:   http.StatusOK,
			expectedContent:  "text/html; charset=utf-8",
			expectedContains: []string{"Welcome to dotsat.work", "Server is running!"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			handler := NewHomeHandler()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()

			// Act
			handler.ServeHTTP(rec, req)

			// Assert
			assertStatus(t, rec.Code, tt.expectedStatus)
			assertContentType(t, rec.Header().Get("Content-Type"), tt.expectedContent)
			assertBodyContains(t, rec.Body.String(), tt.expectedContains)
		})
	}
}

func assertStatus(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("expected status %d, got %d", want, got)
	}
}

func assertContentType(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("expected Content-Type %q, got %q", want, got)
	}
}

func assertBodyContains(t *testing.T, body string, expectedStrings []string) {
	t.Helper()
	for _, expected := range expectedStrings {
		if !strings.Contains(body, expected) {
			t.Errorf("expected body to contain %q, got %q", expected, body)
		}
	}
}
