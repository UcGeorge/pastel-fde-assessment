package handlers

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/UcGeorge/pastel-fde-assessment/services"
)

// ScreeningHandler handles PEP and Sanctions screening requests.
type ScreeningHandler struct {
	tmpl    *template.Template
	service *services.ScreeningService
	useMock bool
}

// NewScreeningHandler creates a new ScreeningHandler.
func NewScreeningHandler(tmpl *template.Template, service *services.ScreeningService, useMock bool) *ScreeningHandler {
	return &ScreeningHandler{tmpl: tmpl, service: service, useMock: useMock}
}

// Page renders the PEP/Sanctions screening form.
func (h *ScreeningHandler) Page(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	h.tmpl.ExecuteTemplate(w, "layout", map[string]any{
		"Page":    "screening",
		"Title":   "PEP & Sanctions Screening - Sigma",
		"UseMock": h.useMock,
	})
}

// CheckPEP handles the HTMX form submission for PEP screening.
func (h *ScreeningHandler) CheckPEP(w http.ResponseWriter, r *http.Request) {
	h.doScreening(w, r, "PEP")
}

// CheckSanction handles the HTMX form submission for Sanctions screening.
func (h *ScreeningHandler) CheckSanction(w http.ResponseWriter, r *http.Request) {
	h.doScreening(w, r, "Sanction")
}

func (h *ScreeningHandler) doScreening(w http.ResponseWriter, r *http.Request, screenType string) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	threshold, _ := strconv.ParseFloat(r.FormValue("threshold"), 64)

	input := services.ScreeningInput{
		Name:      r.FormValue("name"),
		Threshold: threshold,
		Country:   r.FormValue("country"),
	}

	var result *services.ScreeningResult
	if screenType == "PEP" {
		result = h.service.CheckPEP(r.Context(), input)
	} else {
		result = h.service.CheckSanction(r.Context(), input)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	h.tmpl.ExecuteTemplate(w, "screening_result", result)
}
