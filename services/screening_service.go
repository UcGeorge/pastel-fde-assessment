package services

import (
	"context"

	"github.com/UcGeorge/pastel-fde-assessment/pkg/sigma"
)

// ScreeningService wraps the Sigma SDK for PEP and Sanctions screening operations.
type ScreeningService struct {
	client sigma.SigmaClient
}

// NewScreeningService creates a new ScreeningService.
func NewScreeningService(client sigma.SigmaClient) *ScreeningService {
	return &ScreeningService{client: client}
}

// ScreeningInput holds the form data for PEP/Sanctions screening.
type ScreeningInput struct {
	Name      string
	Threshold float64
	Country   string
}

// ScreeningResult bundles the request and response for template rendering.
type ScreeningResult struct {
	Input    ScreeningInput
	Type     string // "PEP" or "Sanction"
	Response *sigma.ScreeningResponse
	Error    string
}

// CheckPEP screens an individual against PEP lists.
func (s *ScreeningService) CheckPEP(ctx context.Context, input ScreeningInput) *ScreeningResult {
	req := &sigma.ScreeningRequest{
		Name:      input.Name,
		Threshold: input.Threshold,
	}
	if input.Country != "" {
		req.Country = &input.Country
	}

	resp, err := s.client.CheckPEP(ctx, req)
	if err != nil {
		return &ScreeningResult{Input: input, Type: "PEP", Error: err.Error()}
	}
	return &ScreeningResult{Input: input, Type: "PEP", Response: resp}
}

// CheckSanction screens an individual against international sanctions lists.
func (s *ScreeningService) CheckSanction(ctx context.Context, input ScreeningInput) *ScreeningResult {
	req := &sigma.ScreeningRequest{
		Name:      input.Name,
		Threshold: input.Threshold,
	}
	if input.Country != "" {
		req.Country = &input.Country
	}

	resp, err := s.client.CheckSanction(ctx, req)
	if err != nil {
		return &ScreeningResult{Input: input, Type: "Sanction", Error: err.Error()}
	}
	return &ScreeningResult{Input: input, Type: "Sanction", Response: resp}
}
