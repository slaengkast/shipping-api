package billing

import (
	"context"

	"github.com/slaengkast/shipping-api/internal/errors"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type rateStore interface {
	GetRateByRegion(context.Context, string) (float32, error)
}

type priceStore interface {
	GetPriceByWeightClass(context.Context, string) (float32, error)
}

type locationStore interface {
	GetByCode(context.Context, string) (*location, error)
}

type Service struct {
	rateStore     rateStore
	priceStore    priceStore
	locationStore locationStore
	logger        zerolog.Logger
}

func NewService(ratestore rateStore, pricestore priceStore, locationstore locationStore) Service {
	return Service{
		rateStore:     ratestore,
		priceStore:    pricestore,
		locationStore: locationstore,
		logger:        log.With().Str("component", "booking").Logger(),
	}
}

func (s Service) CalculateShippingCost(ctx context.Context, origin, destination string, weight float32) (float32, error) {
	s.logger.Info().Str("origin", origin).Str("destination", destination).Float32("weight", weight)

	if origin == "" {
		return 0, errors.FromMessage("empty origin", errors.ErrorInput)
	}
	if destination == "" {
		return 0, errors.FromMessage("empty destination", errors.ErrorInput)
	}

	originLocation, err := s.locationStore.GetByCode(ctx, origin)
	if err != nil {
		return 0, err
	}

	destinationLocation, err := s.locationStore.GetByCode(ctx, destination)
	if err != nil {
		return 0, err
	}

	rate, err := s.rateStore.GetRateByRegion(ctx, getRegion(originLocation, destinationLocation))
	if err != nil {
		return 0, err
	}

	weightClass, err := calculateWeightClass(weight)
	if err != nil {
		return 0, err
	}

	price, err := s.priceStore.GetPriceByWeightClass(ctx, weightClass)
	if err != nil {
		return 0, err
	}

	totalPrice := price * rate
	return totalPrice, nil
}

var ErrorInvalidWeight = errors.FromMessage("invalid weight", errors.ErrorInput)

func calculateWeightClass(weight float32) (string, error) {
	switch {
	case weight < 0:
		return "", ErrorInvalidWeight
	case weight < 10:
		return "small", nil
	case weight < 25:
		return "medium", nil
	case weight < 50:
		return "large", nil
	case weight < 1000:
		return "huge", nil
	default:
		return "", ErrorInvalidWeight
	}
}
