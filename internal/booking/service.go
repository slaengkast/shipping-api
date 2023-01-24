package booking

import (
	"context"

	"github.com/slaengkast/shipping-api/internal/errors"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type store interface {
	GetBooking(context.Context, string) (*booking, error)
	AddBooking(context.Context, *booking) error
}

type billingService interface {
	CalculateShippingCost(context.Context, string, string, float32) (float32, error)
}

type Service struct {
	store          store
	billingService billingService
	logger         zerolog.Logger
}

func NewService(store store, billingService billingService) Service {
	return Service{
		store:          store,
		billingService: billingService,
		logger:         log.With().Str("component", "booking").Logger(),
	}
}

func (s *Service) GetBooking(ctx context.Context, id string) (*booking, error) {
	s.logger.Debug().Str("id", id).Msg("")
	return s.store.GetBooking(ctx, id)
}

func (s *Service) BookShipping(ctx context.Context, origin, destination string, weight float32) (string, error) {
	s.logger.Info().Str("origin", origin).Str("destination", destination).Float32("weight", weight).Msg("")

	if origin == "" {
		return "", errors.FromMessage("empty origin", errors.ErrorInput)
	}
	if destination == "" {
		return "", errors.FromMessage("empty destination", errors.ErrorInput)
	}

	price, err := s.billingService.CalculateShippingCost(ctx, origin, destination, weight)
	if err != nil {
		return "", err
	}

	id := uuid.New().String()
	sh, err := NewBooking(
		id,
		origin,
		destination,
		weight,
		price,
	)
	if err != nil {
		return "", err
	}

	return id, s.store.AddBooking(ctx, sh)
}
