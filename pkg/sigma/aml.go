package sigma

import "context"

// REQUEST MODELS

// ScreeningRequest represents the common payload for both PEP and Sanctions screening.
type ScreeningRequest struct {
	Name      string  `json:"name"`
	Threshold float64 `json:"threshold"`
	Country   *string `json:"country,omitempty"`
}

// RESPONSE MODELS

// ScreeningResponse represents the common response structure from PEP and Sanctions screening.
type ScreeningResponse struct {
	Message    string          `json:"message"`
	Count      int             `json:"count"`
	PageNumber int             `json:"pageNumber"`
	Data       []MatchedEntity `json:"data"`
}

// MatchedEntity holds the detailed information for a single matched individual or organization.
type MatchedEntity struct {
	ID              string      `json:"_id"`
	Type            string      `json:"type"`
	ExternalID      string      `json:"externalId"`
	Version         int         `json:"__v"`
	Addresses       []string    `json:"addresses"`
	Aliases         []string    `json:"aliases"`
	BirthDate       string      `json:"birth_date"` // Using string for flexibility with date formats
	Countries       []string    `json:"countries"`
	CreatedAt       string      `json:"createdAt"`
	UpdatedAt       string      `json:"updatedAt"`
	Dataset         []Dataset   `json:"dataset"`
	EntityType      string      `json:"entityType"`
	Name            string      `json:"name"`
	Sanctions       []string    `json:"sanctions"`
	Photo           string      `json:"photo"`
	PoliticalParty  []string    `json:"politicalParty"`
	Gender          string      `json:"gender"`
	Positions       []Position  `json:"positions"`
	Education       []Education `json:"education"`
	Tags            []string    `json:"tags"`
	SearchScore     float64     `json:"searchScore"`
	ConfidenceScore float64     `json:"confidenceScore"`
	IsDeleted       bool        `json:"isDeleted"`
}

// Dataset represents a source where a match was found.
type Dataset struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// Position represents a role or position held by the matched entity.
type Position struct {
	ID        string  `json:"_id"`
	Title     string  `json:"title"`
	StartDate *string `json:"startDate"`
	EndDate   *string `json:"endDate"`
}

// Education represents an educational institution attended by the matched entity.
type Education struct {
	ID        string  `json:"_id"`
	Title     string  `json:"title"`
	StartDate *string `json:"startDate"`
	EndDate   *string `json:"endDate"`
}

// CLIENT METHODS

// CheckPEP screens an individual against global Politically Exposed Persons (PEP) lists.
func (c *Client) CheckPEP(ctx context.Context, req *ScreeningRequest) (*ScreeningResponse, error) {
	var resp ScreeningResponse

	err := c.doSigmaRequest(ctx, "POST", "api/v1/aml/pep/instant", req, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// CheckSanction screens an individual against international sanctions lists.
func (c *Client) CheckSanction(ctx context.Context, req *ScreeningRequest) (*ScreeningResponse, error) {
	var resp ScreeningResponse

	// NOTE: The provided documentation has a typo "sacntion".
	err := c.doSigmaRequest(ctx, "POST", "api/v1/aml/sacntion/instant", req, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}
