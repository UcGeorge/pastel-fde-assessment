package di

import (
	"github.com/UcGeorge/pastel-fde-assessment/internal/config"
	"github.com/UcGeorge/pastel-fde-assessment/pkg/sigma"
	"github.com/UcGeorge/pastel-fde-assessment/services"
	"github.com/rs/zerolog/log"
	"github.com/samber/do/v2"
)

// NewContainer creates and wires the dependency injection container.
func NewContainer() *do.RootScope {
	injector := do.New()

	// Configuration
	do.Provide(injector, func(i do.Injector) (*config.Config, error) {
		return config.Load(), nil
	})

	// Sigma SDK Client — either mock or live depending on config
	do.Provide(injector, func(i do.Injector) (sigma.SigmaClient, error) {
		cfg := do.MustInvoke[*config.Config](i)

		if cfg.UseMock {
			log.Info().Msg("Using MOCK Sigma client (set USE_MOCK=false for live API)")
			return sigma.NewMockClient(), nil
		}

		log.Info().Msg("Using LIVE Sigma client")
		client := sigma.New(
			cfg.APIKey,
			cfg.APISecret,
			sigma.WithBaseURL(cfg.BaseURL),
			sigma.WithAMLBaseURL(cfg.AMLBaseURL),
		)
		return client, nil
	})

	// Services
	do.Provide(injector, func(i do.Injector) (*services.TransactionService, error) {
		client := do.MustInvoke[sigma.SigmaClient](i)
		return services.NewTransactionService(client), nil
	})

	do.Provide(injector, func(i do.Injector) (*services.ScreeningService, error) {
		client := do.MustInvoke[sigma.SigmaClient](i)
		return services.NewScreeningService(client), nil
	})

	do.Provide(injector, func(i do.Injector) (*services.AdverseMediaService, error) {
		client := do.MustInvoke[sigma.SigmaClient](i)
		return services.NewAdverseMediaService(client), nil
	})

	return injector
}
