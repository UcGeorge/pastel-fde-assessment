package sigma

import "context"

// REQUEST MODELS

// AdverseMediaRequest represents the payload for the adverse media screening endpoint.
type AdverseMediaRequest struct {
	Query string `json:"q"`
	Limit int    `json:"limit"`
}

// RESPONSE MODELS

// AdverseMediaResponse is the top-level response from the adverse media endpoint.
type AdverseMediaResponse struct {
	Message string             `json:"message"`
	Data    *AdverseMediaData  `json:"data"`
}

// AdverseMediaData holds the details of the adverse media check request status.
type AdverseMediaData struct {
	ID               string   `json:"id"`
	Query            string   `json:"query"`
	BusinessProfile  string   `json:"businessProfile"`
	Status           string   `json:"status"`
	CreatedAt        string   `json:"createdAt"`
	UpdatedAt        string   `json:"updatedAt"`
	FindingsReturned *bool    `json:"findingsReturned,omitempty"`
	RiskCategory     *string  `json:"riskCategory,omitempty"`
	Sources          []string `json:"sources,omitempty"`
}

// CLIENT METHODS

// CheckAdverseMedia submits an adverse media screening request.
// NOTE: This endpoint is webhook-based, so the response contains a pending status
// and an ID. The full results would be delivered via a configured webhook.
func (c *Client) CheckAdverseMedia(ctx context.Context, req *AdverseMediaRequest) (*AdverseMediaResponse, error) {
	var resp AdverseMediaResponse

	err := c.doAMLRequest(ctx, "POST", "api/v1/adverse-media", req, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}
