package sigma

import "context"

// SigmaClient defines the interface for interacting with the Sigma API.
// Both the live Client and MockClient implement this interface, enabling
// dependency injection and clean testing without hitting the live API.
type SigmaClient interface {
	// SubmitTransaction sends a transaction to the Sigma instant monitoring endpoint.
	SubmitTransaction(ctx context.Context, req *SubmitTransactionRequest) (*TransactionResponse, error)

	// CheckPEP screens an individual against global PEP lists.
	CheckPEP(ctx context.Context, req *ScreeningRequest) (*ScreeningResponse, error)

	// CheckSanction screens an individual against international sanctions lists.
	CheckSanction(ctx context.Context, req *ScreeningRequest) (*ScreeningResponse, error)

	// CheckAdverseMedia submits an adverse media screening request.
	CheckAdverseMedia(ctx context.Context, req *AdverseMediaRequest) (*AdverseMediaResponse, error)
}

// Compile-time assertion that Client implements SigmaClient.
var _ SigmaClient = (*Client)(nil)
