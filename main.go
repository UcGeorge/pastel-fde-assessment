package main

import (
	"context"
	"os"
	"time"

	"github.com/UcGeorge/pastel-fde-assessment/pkg/sigma"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// Setup structured, human-readable logging for development
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	// Load configuration from environment variables (Never hardcode secrets!)
	apiKey := "[[REDACTED]]"
	apiSecret := "[[REDACTED]]"
	apiBaseURL := "https://sigmaprod.sabipay.com/"

	// Initialize the SDK client
	client := sigma.New(
		apiKey, apiSecret,
		sigma.WithBaseURL(apiBaseURL),
	)
	log.Info().Msg("Sigma client initialized successfully.")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	submitTransaction(ctx, client)
	checkSanction(ctx, client)
	checkPEP(ctx, client)
}

func submitTransaction(ctx context.Context, client *sigma.Client) {
	// Create a realistic request payload as required by the test
	// This should contain valid data to get a meaningful response from the API
	req := sigma.SubmitTransactionRequest{
		TransactionData: sigma.TransactionData{
			Reference:         "2303pee2fc",
			BalanceBefore:     ptr(1000.0),
			Amount:            100.5,
			Currency:          "USD",
			Status:            true,
			IsExternalPayment: false,
			SenderAccount:     ptr("01929393923"),
			ReceiverAccount:   ptr("93829392233"),
			Channel:           sigma.ChannelPOS,
			Type:              sigma.TxTypeDebit,
			TransactionDate:   time.Now().UTC(),
			Narration:         ptr("Payment for services"),
		},
		Device: &sigma.Device{
			DeviceID:     "J020D23020D03300303203D3232DDD",
			Manufacturer: ptr("Apple"),
			Name:         ptr("iPhone 14 Pro"),
			OSName:       ptr("iOS"),
			OSVersion:    ptr("17.1.0"),
		},
		AnonymizedUserData: sigma.AnonymizedUserData{
			UniqueID:              "e8baeb9c-e563-11ed-b5ea-0242ac120002",
			AccountType:           ptr(sigma.AccountTypeIndividual),
			BusinessCategory:      ptr("retail"),
			IsPhoneNumberVerified: ptr(true),
			IsBanned:              false,
			DateJoined:            ptr(time.Date(2022, 1, 1, 23, 58, 0, 0, time.UTC)),
			Age:                   ptr(29),
			IsIdentityVerified:    true,
			State:                 ptr("lagos"),
			City:                  ptr("Ikeja"),
			Country:               ptr("Nigeria"),
		},
	}

	// Execute the SDK method
	log.Info().Msg("Submitting transaction for monitoring...")
	res, err := client.SubmitTransaction(ctx, &req)
	if err != nil {
		// Use the structured error logger for rich context
		log.Error().Err(err).Msg("Failed to submit transaction to Sigma API")
	}

	// Log the successful response in a structured way
	log.Info().Interface("response", res).Msg("Successfully received transaction response from Sigma")
}

func checkPEP(ctx context.Context, client *sigma.Client) {
	// Create a realistic request payload as required by the test
	// This should contain valid data to get a meaningful response from the API
	req := sigma.ScreeningRequest{
		Name:      "George Uche-Umeh",
		Threshold: 0.5,
		Country:   ptr("Nigeria"),
	}

	// Execute the SDK method
	log.Info().Msg("Checking PEP...")
	res, err := client.CheckPEP(ctx, &req)
	if err != nil {
		// Use the structured error logger for rich context
		log.Error().Err(err).Msg("Failed to check PEP")
	}

	// Log the successful response in a structured way
	log.Info().Interface("response", res).Msg("Successfully received response from Sigma")
}

func checkSanction(ctx context.Context, client *sigma.Client) {
	// Create a realistic request payload as required by the test
	// This should contain valid data to get a meaningful response from the API
	req := sigma.ScreeningRequest{
		Name:      "George Uche-Umeh",
		Threshold: 0.5,
		Country:   ptr("Nigeria"),
	}

	// Execute the SDK method
	log.Info().Msg("Checking sanctions...")
	res, err := client.CheckSanction(ctx, &req)
	if err != nil {
		// Use the structured error logger for rich context
		log.Error().Err(err).Msg("Failed to check sanctions")
	}

	// Log the successful response in a structured way
	log.Info().Interface("response", res).Msg("Successfully received response from Sigma")
}

func ptr[T any](v T) *T {
	return &v
}
