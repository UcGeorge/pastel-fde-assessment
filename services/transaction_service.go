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
	// TransactionData
	Reference         string
	Amount            float64
	IsExternalPayment bool
	Type              string
	Channel           string
	TransactionDate   time.Time
	Status            bool
	Currency          string
	SenderAccount     string
	ReceiverAccount   string
	BalanceBefore     float64
	Email             string
	Narration         string
	Refund            bool
	IsCheque          bool
	VasReceiver       string
	IsInternalAccount bool
	IsStaffAccount    bool
	SessionID         string
	IsDormantAccount  bool

	// AnonymizedUserData
	UniqueID              string
	IsBanned              bool
	IsIdentityVerified    bool
	UserEmail             string
	AccountType           string
	BusinessCategory      string
	IsPhoneNumberVerified bool
	DateJoined            time.Time // if zero, not passed
	Age                   int
	State                 string
	City                  string
	UserCountry           string

	// Device
	DeviceID     string
	Manufacturer string
	DeviceName   string
	OSName       string
	OSVersion    string

	// Location
	Latitude   string
	Longitude  string
	LocCountry string

	// ThirdPartyUserData
	ThirdPartyUniqueID string
	CardPan            string
	ThirdPartyEmail    string

	// Limits
	DailyLimit                 float64
	OverdraftLimit             float64
	IndividualTransactionLimit float64

	// ScreeningData
	ScreeningSenderName   string
	ScreeningReceiverName string

	// Beneficiary
	HasBeneficiary          bool
	IsRegisteredBeneficiary bool
	IsNewBeneficiary        bool

	// IPAddress
	IPAddress string
}

// TransactionResult bundles the request and response for template rendering.
type TransactionResult struct {
	Input    TransactionInput
	Response *sigma.TransactionResponse
	Error    string
}

func ptrStr(v string) *string {
	if v == "" {
		return nil
	}
	return &v
}

func ptrFloat(v float64) *float64 {
	if v == 0 {
		return nil
	}
	return &v
}

func ptrInt(v int) *int {
	if v == 0 {
		return nil
	}
	return &v
}

func ptrBoolForm(v bool, passed bool) *bool {
	if !passed {
		return nil
	}
	return &v
}

func ptrTime(v time.Time) *time.Time {
	if v.IsZero() {
		return nil
	}
	return &v
}


// SubmitTransaction builds a realistic request from form input and calls the API.
func (s *TransactionService) SubmitTransaction(ctx context.Context, input TransactionInput) *TransactionResult {
	req := &sigma.SubmitTransactionRequest{
		TransactionData: sigma.TransactionData{
			Reference:         input.Reference,
			Amount:            input.Amount,
			IsExternalPayment: input.IsExternalPayment,
			Type:              sigma.TransactionType(input.Type),
			Channel:           sigma.TransactionChannel(input.Channel),
			TransactionDate:   input.TransactionDate,
			Status:            input.Status,
			Currency:          input.Currency,
			SenderAccount:     ptrStr(input.SenderAccount),
			ReceiverAccount:   ptrStr(input.ReceiverAccount),
			BalanceBefore:     ptrFloat(input.BalanceBefore),
			Email:             ptrStr(input.Email),
			Narration:         ptrStr(input.Narration),
			Refund:            ptrBoolForm(input.Refund, true),
			IsCheque:          ptrBoolForm(input.IsCheque, true),
			VasReceiver:       ptrStr(input.VasReceiver),
			IsInternalAccount: ptrBoolForm(input.IsInternalAccount, true),
			IsStaffAccount:    ptrBoolForm(input.IsStaffAccount, true),
			SessionID:         ptrStr(input.SessionID),
			IsDormantAccount:  ptrBoolForm(input.IsDormantAccount, true),
		},
		AnonymizedUserData: sigma.AnonymizedUserData{
			UniqueID:              input.UniqueID,
			IsBanned:              input.IsBanned,
			IsIdentityVerified:    input.IsIdentityVerified,
			Email:                 ptrStr(input.UserEmail),
			AccountType:           (*sigma.AccountType)(ptrStr(input.AccountType)),
			BusinessCategory:      ptrStr(input.BusinessCategory),
			IsPhoneNumberVerified: ptrBoolForm(input.IsPhoneNumberVerified, true),
			DateJoined:            ptrTime(input.DateJoined),
			Age:                   ptrInt(input.Age),
			State:                 ptrStr(input.State),
			City:                  ptrStr(input.City),
			Country:               ptrStr(input.UserCountry),
		},
	}

	if input.DeviceID != "" {
		req.Device = &sigma.Device{
			DeviceID:     input.DeviceID,
			Manufacturer: ptrStr(input.Manufacturer),
			Name:         ptrStr(input.DeviceName),
			OSName:       ptrStr(input.OSName),
			OSVersion:    ptrStr(input.OSVersion),
		}
	}

	if input.Latitude != "" || input.Longitude != "" || input.LocCountry != "" {
		req.Location = &sigma.Location{
			Latitude:  ptrStr(input.Latitude),
			Longitude: ptrStr(input.Longitude),
			Country:   ptrStr(input.LocCountry),
		}
	}

	if input.ThirdPartyUniqueID != "" || input.CardPan != "" {
		req.ThirdPartyUserData = &sigma.ThirdPartyUserData{
			UniqueID: input.ThirdPartyUniqueID,
			CardPan:  input.CardPan,
			Email:    ptrStr(input.ThirdPartyEmail),
		}
	}

	if input.DailyLimit > 0 || input.OverdraftLimit > 0 || input.IndividualTransactionLimit > 0 {
		req.Limits = &sigma.Limits{
			DailyLimit:                 ptrFloat(input.DailyLimit),
			OverdraftLimit:             ptrFloat(input.OverdraftLimit),
			IndividualTransactionLimit: ptrFloat(input.IndividualTransactionLimit),
		}
	}

	if input.ScreeningSenderName != "" || input.ScreeningReceiverName != "" {
		req.ScreeningData = &sigma.ScreeningData{
			SenderName:   ptrStr(input.ScreeningSenderName),
			ReceiverName: ptrStr(input.ScreeningReceiverName),
		}
	}

	if input.HasBeneficiary {
		req.Beneficiary = &sigma.Beneficiary{
			IsRegisteredBeneficiary: input.IsRegisteredBeneficiary,
			IsNewBeneficiary:        input.IsNewBeneficiary,
		}
	}
	
	if input.IPAddress != "" {
		req.IPAddress = ptrStr(input.IPAddress)
	}

	resp, err := s.client.SubmitTransaction(ctx, req)
	if err != nil {
		return &TransactionResult{Input: input, Error: err.Error()}
	}
	return &TransactionResult{Input: input, Response: resp}
}
