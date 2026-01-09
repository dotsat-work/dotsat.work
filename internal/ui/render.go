package ui

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/a-h/templ"
)

// Render renders a templ component to the response writer.
// If rendering fails, it logs the error and returns a 500 Internal Server Error.
func Render(w http.ResponseWriter, r *http.Request, c templ.Component) {
	err := c.Render(r.Context(), w)
	if err != nil {
		slog.Error("render failed", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

// RenderFragment renders specific fragments of a templ component.
// This is useful for HTMX partial updates where you only want to render
// specific parts of a component identified by their fragment IDs.
func RenderFragment(w http.ResponseWriter, r *http.Request, c templ.Component, fragmentIDs ...any) {
	err := templ.RenderFragments(r.Context(), w, c, fragmentIDs...)
	if err != nil {
		slog.Error("render fragment failed", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

// RenderOOB renders a component wrapped in an HTMX out-of-band swap container.
// This allows updating elements outside the main target element in HTMX responses.
//
// Example:
//
//	RenderOOB(w, r, myComponent, "innerHTML:#sidebar")
func RenderOOB(w http.ResponseWriter, r *http.Request, c templ.Component, target string) {
	// Write OOB wrapper start
	_, err := fmt.Fprintf(w, `<div hx-swap-oob="%s">`, target)
	if err != nil {
		slog.Error("render oob write wrapper start failed", "error", err)
		return
	}

	// Render component
	err = c.Render(r.Context(), w)
	if err != nil {
		slog.Error("render oob component render failed", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	// Write OOB wrapper end
	_, err = w.Write([]byte(`</div>`))
	if err != nil {
		slog.Error("render oob write wrapper end failed", "error", err)
	}
}
