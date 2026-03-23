package sigma

import (
	"context"
	"time"
)

// ENUMS

// TransactionType represents the direction of the transaction
type TransactionType string

const (
	TxTypeDebit  TransactionType = "debit"
	TxTypeCredit TransactionType = "credit"
)

// TransactionChannel represents how the transaction was initiated
type TransactionChannel string

const (
	ChannelCardPayment    TransactionChannel = "card payment"
	ChannelWebTransfer    TransactionChannel = "web transfer"
	ChannelBankTransfer   TransactionChannel = "bank transfer"
	ChannelMobileTransfer TransactionChannel = "mobile transfer"
	ChannelCashPayment    TransactionChannel = "cash payment"
	ChannelCashWithdrawal TransactionChannel = "cash withdrawal"
	ChannelAirtime        TransactionChannel = "airtime"
	ChannelCable          TransactionChannel = "cable"
	ChannelSportsBetting  TransactionChannel = "sports betting"
	ChannelElectricity    TransactionChannel = "electricity"
	ChannelInternet       TransactionChannel = "internet"
	ChannelDataPurchase   TransactionChannel = "data purchase"
	ChannelATM            TransactionChannel = "atm"
	ChannelPOS            TransactionChannel = "pos"
)

// AccountType represents the type of user account
type AccountType string

const (
	AccountTypeIndividual AccountType = "individual"
	AccountTypeCorporate  AccountType = "corporate"
)

// ActionType represents the final decision on the transaction.
type ActionType string

const (
	ActionApproved ActionType = "approved"
	ActionDeclined ActionType = "rejected"
)

// Numerical code representation of the result. 1 = Approved, 0 = Rejected
type ActionCode int

const (
	ActionCodeApproved ActionCode = 1
	ActionCodeRejected ActionCode = 0
)

// SeverityLevel represents the severity of a flagged rule
type SeverityLevel string

const (
	SeverityLow    SeverityLevel = "low"
	SeverityMedium SeverityLevel = "medium"
	SeverityHigh   SeverityLevel = "high"
)

// REQUEST MODELS

// SubmitTransactionRequest represents the payload for transaction screening
type SubmitTransactionRequest struct {
	TransactionData    TransactionData     `json:"transactionData"`
	AnonymizedUserData AnonymizedUserData  `json:"anonymizedUserData"`
	Device             *Device             `json:"device,omitempty"`
	Location           *Location           `json:"location,omitempty"`
	ThirdPartyUserData *ThirdPartyUserData `json:"thirdPartyUserData,omitempty"`
	Limits             *Limits             `json:"limits,omitempty"`
	ScreeningData      *ScreeningData      `json:"screeningData,omitempty"`
	Beneficiary        *Beneficiary        `json:"beneficiary,omitempty"`
	IPAddress          *string             `json:"ipAddress,omitempty"`
}

type TransactionData struct {
	Reference         string             `json:"reference"`
	Amount            float64            `json:"amount"`
	IsExternalPayment bool               `json:"isExternalPayment"`
	Type              TransactionType    `json:"type"`
	Channel           TransactionChannel `json:"channel"`
	TransactionDate   time.Time          `json:"transactionDate"`
	Status            bool               `json:"status"`
	Currency          string             `json:"currency"`
	SenderAccount     *string            `json:"senderAccount,omitempty"`
	ReceiverAccount   *string            `json:"receiverAccount,omitempty"`
	BalanceBefore     *float64           `json:"balanceBefore,omitempty"`
	Email             *string            `json:"email,omitempty"`
	Narration         *string            `json:"narration,omitempty"`
	Refund            *bool              `json:"refund,omitempty"`
	IsCheque          *bool              `json:"isCheque,omitempty"`
	VasReceiver       *string            `json:"vasReceiver,omitempty"`
	IsInternalAccount *bool              `json:"isInternalAccount,omitempty"`
	IsStaffAccount    *bool              `json:"isStaffAccount,omitempty"`
	SessionID         *string            `json:"sessionId,omitempty"`
	IsDormantAccount  *bool              `json:"isDormantAccount,omitempty"`
}

type AnonymizedUserData struct {
	UniqueID              string       `json:"uniqueId"`
	IsBanned              bool         `json:"isBanned"`
	IsIdentityVerified    bool         `json:"isIdentityVerified"`
	Email                 *string      `json:"email,omitempty"`
	AccountType           *AccountType `json:"accountType,omitempty"`
	BusinessCategory      *string      `json:"businessCategory,omitempty"`
	IsPhoneNumberVerified *bool        `json:"isPhoneNumberVerified,omitempty"`
	DateJoined            *time.Time   `json:"dateJoined,omitempty"`
	Age                   *int         `json:"age,omitempty"`
	State                 *string      `json:"state,omitempty"`
	City                  *string      `json:"city,omitempty"`
	Country               *string      `json:"country,omitempty"`
}

type Device struct {
	DeviceID     string  `json:"deviceId"`
	Manufacturer *string `json:"manufacturer,omitempty"`
	Name         *string `json:"name,omitempty"`
	OSName       *string `json:"osName,omitempty"`
	OSVersion    *string `json:"osVersion,omitempty"`
}

type Location struct {
	Latitude  *string `json:"latitude,omitempty"`
	Longitude *string `json:"longitude,omitempty"`
	Country   *string `json:"country,omitempty"`
}

type ThirdPartyUserData struct {
	UniqueID string  `json:"uniqueId"`
	CardPan  string  `json:"cardPan"`
	Email    *string `json:"email,omitempty"`
}

type Limits struct {
	DailyLimit                 *float64 `json:"dailyLimit,omitempty"`
	OverdraftLimit             *float64 `json:"overdraftLimit,omitempty"`
	IndividualTransactionLimit *float64 `json:"individualTransactionLimit,omitempty"`
}

type ScreeningData struct {
	SenderName   *string `json:"senderName,omitempty"`
	ReceiverName *string `json:"receiverName,omitempty"`
}

type Beneficiary struct {
	IsRegisteredBeneficiary bool `json:"isRegisteredBeneficiary"`
	IsNewBeneficiary        bool `json:"isNewBeneficiary"`
}

// RESPONSE MODELS

// TransactionResponse represents the result of the monitoring API
type TransactionResponse struct {
	Message string                   `json:"message"`
	Data    *TransactionResponseData `json:"data"`
}

type TransactionResponseData struct {
	TransactionID string            `json:"transactionId"`
	RiskScore     string            `json:"riskScore"`
	Action        TransactionAction `json:"action"`
	RuleResult    string            `json:"ruleResult"`
	Reason        TransactionReason `json:"reason"`
	Screening     ScreeningResult   `json:"screening"`
}

type TransactionAction struct {
	Result ActionType `json:"result"`
	Code   ActionCode `json:"code"`
}

type TransactionReason struct {
	Code     string          `json:"code"`
	Message  string          `json:"message"`
	Severity SeverityLevel   `json:"severity"`
	Rule     TransactionRule `json:"rule"`
}

type TransactionRule struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ScreeningResult struct {
	Sender   ScreeningParty `json:"sender"`
	Receiver ScreeningParty `json:"receiver"`
}

type ScreeningParty struct {
	PEP      string `json:"pep"`
	Sanction string `json:"sanction"`
}

// CLIENT METHODS

// SubmitTransaction sends a transaction to the Sigma instant monitoring endpoint.
func (c *Client) SubmitTransaction(ctx context.Context, req *SubmitTransactionRequest) (*TransactionResponse, error) {
	var resp TransactionResponse

	err := c.doSigmaRequest(ctx, "POST", "api/v1/transaction-monitoring/instant", req, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}
