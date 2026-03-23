package handlers

import (
	"html/template"
	"net/http"
)

// HomeHandler serves the dashboard landing page.
type HomeHandler struct {
	tmpl    *template.Template
	useMock bool
}

// NewHomeHandler creates a new HomeHandler.
func NewHomeHandler(tmpl *template.Template, useMock bool) *HomeHandler {
	return &HomeHandler{tmpl: tmpl, useMock: useMock}
}

// ServeHTTP renders the home/dashboard page.
func (h *HomeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	h.tmpl.ExecuteTemplate(w, "layout", map[string]any{
		"Page":    "home",
		"Title":   "Sigma Compliance Dashboard",
		"UseMock": h.useMock,
	})
}
