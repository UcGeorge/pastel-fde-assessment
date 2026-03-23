package handlers

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/UcGeorge/pastel-fde-assessment/services"
)

// AdverseMediaHandler handles adverse media screening requests.
type AdverseMediaHandler struct {
	tmpl    *template.Template
	service *services.AdverseMediaService
}

// NewAdverseMediaHandler creates a new AdverseMediaHandler.
func NewAdverseMediaHandler(tmpl *template.Template, service *services.AdverseMediaService) *AdverseMediaHandler {
	return &AdverseMediaHandler{tmpl: tmpl, service: service}
}

// Page renders the adverse media screening form.
func (h *AdverseMediaHandler) Page(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	h.tmpl.ExecuteTemplate(w, "layout", map[string]any{
		"Page":  "adverse_media",
		"Title": "Adverse Media Screening — Sigma",
	})
}

// Check handles the HTMX form submission for adverse media screening.
func (h *AdverseMediaHandler) Check(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	limit, _ := strconv.Atoi(r.FormValue("limit"))
	if limit <= 0 {
		limit = 10
	}

	input := services.AdverseMediaInput{
		Query: r.FormValue("query"),
		Limit: limit,
	}

	result := h.service.Check(r.Context(), input)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	h.tmpl.ExecuteTemplate(w, "adverse_media_result", result)
}
