package handlers

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/UcGeorge/pastel-fde-assessment/services"
)

// TransactionHandler handles transaction monitoring requests.
type TransactionHandler struct {
	tmpl    *template.Template
	service *services.TransactionService
}

// NewTransactionHandler creates a new TransactionHandler.
func NewTransactionHandler(tmpl *template.Template, service *services.TransactionService) *TransactionHandler {
	return &TransactionHandler{tmpl: tmpl, service: service}
}

// Page renders the transaction monitoring form page.
func (h *TransactionHandler) Page(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	h.tmpl.ExecuteTemplate(w, "layout", map[string]any{
		"Page":  "transaction",
		"Title": "Transaction Monitoring — Sigma",
	})
}

// Submit handles the HTMX form submission and returns the result partial.
func (h *TransactionHandler) Submit(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	amount, _ := strconv.ParseFloat(r.FormValue("amount"), 64)

	input := services.TransactionInput{
		Reference:       r.FormValue("reference"),
		Amount:          amount,
		Currency:        r.FormValue("currency"),
		SenderAccount:   r.FormValue("sender_account"),
		ReceiverAccount: r.FormValue("receiver_account"),
		Channel:         r.FormValue("channel"),
		Type:            r.FormValue("type"),
		Narration:       r.FormValue("narration"),
		UniqueID:        r.FormValue("unique_id"),
		Country:         r.FormValue("country"),
		DeviceID:        r.FormValue("device_id"),
	}

	result := h.service.SubmitTransaction(r.Context(), input)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	h.tmpl.ExecuteTemplate(w, "transaction_result", result)
}
