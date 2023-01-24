package billing

import (
	"context"
	"fmt"

	"github.com/slaengkast/shipping-api/internal/errors"
)

type inMemoryPriceStore struct {
	prices map[string]float32
}

func NewInMemoryPriceStore(prices map[string]float32) inMemoryPriceStore {
	return inMemoryPriceStore{
		prices: prices,
	}
}

func (r inMemoryPriceStore) GetPriceByWeightClass(ctx context.Context, class string) (float32, error) {
	if _, ok := r.prices[class]; !ok {
		return 0, errors.FromMessage(fmt.Sprintf("no price found for class %s", class), errors.ErrorInternal)
	}

	return r.prices[class], nil
}
