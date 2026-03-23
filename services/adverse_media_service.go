package services

import (
	"context"

	"github.com/UcGeorge/pastel-fde-assessment/pkg/sigma"
)

// AdverseMediaService wraps the Sigma SDK for adverse media screening operations.
type AdverseMediaService struct {
	client *sigma.Client
}

// NewAdverseMediaService creates a new AdverseMediaService.
func NewAdverseMediaService(client *sigma.Client) *AdverseMediaService {
	return &AdverseMediaService{client: client}
}

// AdverseMediaInput holds the form data for adverse media screening.
type AdverseMediaInput struct {
	Query string
	Limit int
}

// AdverseMediaResult bundles the request and response for template rendering.
type AdverseMediaResult struct {
	Input    AdverseMediaInput
	Response *sigma.AdverseMediaResponse
	Error    string
}

// Check submits an adverse media screening request.
func (s *AdverseMediaService) Check(ctx context.Context, input AdverseMediaInput) *AdverseMediaResult {
	req := &sigma.AdverseMediaRequest{
		Query: input.Query,
		Limit: input.Limit,
	}

	resp, err := s.client.CheckAdverseMedia(ctx, req)
	if err != nil {
		return &AdverseMediaResult{Input: input, Error: err.Error()}
	}
	return &AdverseMediaResult{Input: input, Response: resp}
}
