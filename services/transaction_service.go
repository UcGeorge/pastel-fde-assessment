package services

import (
	"context"
	"time"

	"github.com/UcGeorge/pastel-fde-assessment/pkg/sigma"
)

// TransactionService wraps the Sigma SDK for transaction monitoring operations.
type TransactionService struct {
	client sigma.SigmaClient
}

// NewTransactionService creates a new TransactionService with the given SDK client.
func NewTransactionService(client sigma.SigmaClient) *TransactionService {
	return &TransactionService{client: client}
}

// TransactionInput holds the form data submitted by the user.
type TransactionInput struct {
	Reference       string
	Amount          float64
	Currency        string
	SenderAccount   string
	ReceiverAccount string
	Channel         string
	Type            string
	Narration       string
	UniqueID        string
	Country         string
	DeviceID        string
}

// TransactionResult bundles the request and response for template rendering.
type TransactionResult struct {
	Input    TransactionInput
	Response *sigma.TransactionResponse
	Error    string
}

func ptr[T any](v T) *T {
	return &v
}

// SubmitTransaction builds a realistic request from form input and calls the API.
func (s *TransactionService) SubmitTransaction(ctx context.Context, input TransactionInput) *TransactionResult {
	req := &sigma.SubmitTransactionRequest{
		TransactionData: sigma.TransactionData{
			Reference:         input.Reference,
			Amount:            input.Amount,
			Currency:          input.Currency,
			Status:            true,
			IsExternalPayment: false,
			SenderAccount:     ptr(input.SenderAccount),
			ReceiverAccount:   ptr(input.ReceiverAccount),
			Channel:           sigma.TransactionChannel(input.Channel),
			Type:              sigma.TransactionType(input.Type),
			TransactionDate:   time.Now().UTC(),
			Narration:         ptr(input.Narration),
			BalanceBefore:     ptr(1000.0),
		},
		Device: &sigma.Device{
			DeviceID:     input.DeviceID,
			Manufacturer: ptr("Apple"),
			Name:         ptr("iPhone 14 Pro"),
			OSName:       ptr("iOS"),
			OSVersion:    ptr("17.1.0"),
		},
		AnonymizedUserData: sigma.AnonymizedUserData{
			UniqueID:              input.UniqueID,
			AccountType:           ptr(sigma.AccountTypeIndividual),
			BusinessCategory:      ptr("retail"),
			IsPhoneNumberVerified: ptr(true),
			IsBanned:              false,
			DateJoined:            ptr(time.Date(2022, 1, 1, 23, 58, 0, 0, time.UTC)),
			Age:                   ptr(29),
			IsIdentityVerified:    true,
			State:                 ptr("lagos"),
			City:                  ptr("Ikeja"),
			Country:               ptr(input.Country),
		},
	}

	resp, err := s.client.SubmitTransaction(ctx, req)
	if err != nil {
		return &TransactionResult{Input: input, Error: err.Error()}
	}
	return &TransactionResult{Input: input, Response: resp}
}
