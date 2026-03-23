package handlers

import (
	"html/template"
	"net/http"
	"strconv"
	"time"

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

	parseFloat := func(k string) float64 {
		v, _ := strconv.ParseFloat(r.FormValue(k), 64)
		return v
	}
	parseInt := func(k string) int {
		v, _ := strconv.Atoi(r.FormValue(k))
		return v
	}
	parseTime := func(k string) time.Time {
		// Attempt to parse datetime-local HTML input. Example: "2023-11-20T14:30"
		if r.FormValue(k) == "" {
			return time.Time{}
		}
		t, err := time.Parse("2006-01-02T15:04", r.FormValue(k))
		if err != nil {
		    t2, err2 := time.Parse("2006-01-02", r.FormValue(k))
		    if err2 != nil {
			    return time.Now()
			}
			return t2
		}
		return t
	}

	input := services.TransactionInput{
		Reference:                  r.FormValue("reference"),
		Amount:                     parseFloat("amount"),
		IsExternalPayment:          r.FormValue("is_external_payment") == "true",
		Type:                       r.FormValue("type"),
		Channel:                    r.FormValue("channel"),
		TransactionDate:            parseTime("transaction_date"),
		Status:                     r.FormValue("status") == "true",
		Currency:                   r.FormValue("currency"),
		SenderAccount:              r.FormValue("sender_account"),
		ReceiverAccount:            r.FormValue("receiver_account"),
		BalanceBefore:              parseFloat("balance_before"),
		Email:                      r.FormValue("email"),
		Narration:                  r.FormValue("narration"),
		Refund:                     r.FormValue("refund") == "true",
		IsCheque:                   r.FormValue("is_cheque") == "true",
		VasReceiver:                r.FormValue("vas_receiver"),
		IsInternalAccount:          r.FormValue("is_internal_account") == "true",
		IsStaffAccount:             r.FormValue("is_staff_account") == "true",
		SessionID:                  r.FormValue("session_id"),
		IsDormantAccount:           r.FormValue("is_dormant_account") == "true",

		UniqueID:                   r.FormValue("unique_id"),
		IsBanned:                   r.FormValue("is_banned") == "true",
		IsIdentityVerified:         r.FormValue("is_identity_verified") == "true",
		UserEmail:                  r.FormValue("user_email"),
		AccountType:                r.FormValue("account_type"),
		BusinessCategory:           r.FormValue("business_category"),
		IsPhoneNumberVerified:      r.FormValue("is_phone_number_verified") == "true",
		DateJoined:                 parseTime("date_joined"),
		Age:                        parseInt("age"),
		State:                      r.FormValue("state"),
		City:                       r.FormValue("city"),
		UserCountry:                r.FormValue("user_country"),

		DeviceID:                   r.FormValue("device_id"),
		Manufacturer:               r.FormValue("manufacturer"),
		DeviceName:                 r.FormValue("device_name"),
		OSName:                     r.FormValue("os_name"),
		OSVersion:                  r.FormValue("os_version"),

		Latitude:                   r.FormValue("latitude"),
		Longitude:                  r.FormValue("longitude"),
		LocCountry:                 r.FormValue("loc_country"),

		ThirdPartyUniqueID:         r.FormValue("third_party_unique_id"),
		CardPan:                    r.FormValue("card_pan"),
		ThirdPartyEmail:            r.FormValue("third_party_email"),

		DailyLimit:                 parseFloat("daily_limit"),
		OverdraftLimit:             parseFloat("overdraft_limit"),
		IndividualTransactionLimit: parseFloat("individual_transaction_limit"),

		ScreeningSenderName:        r.FormValue("screening_sender_name"),
		ScreeningReceiverName:      r.FormValue("screening_receiver_name"),

		HasBeneficiary:             r.FormValue("has_beneficiary") == "true",
		IsRegisteredBeneficiary:    r.FormValue("is_registered_beneficiary") == "true",
		IsNewBeneficiary:           r.FormValue("is_new_beneficiary") == "true",

		IPAddress:                  r.FormValue("ip_address"),
	}

	result := h.service.SubmitTransaction(r.Context(), input)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	h.tmpl.ExecuteTemplate(w, "transaction_result", result)
}
