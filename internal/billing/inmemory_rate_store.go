package billing

import (
	"context"
	"fmt"

	"github.com/slaengkast/shipping-api/internal/errors"
)

type inMemoryRateStore struct {
	rates map[string]float32
}

func NewInMemoryRateStore(rates map[string]float32) inMemoryRateStore {
	return inMemoryRateStore{
		rates: rates,
	}
}

func (r inMemoryRateStore) GetRateByRegion(ctx context.Context, region string) (float32, error) {
	if _, ok := r.rates[region]; !ok {
		return 0, errors.FromMessage(fmt.Sprintf("no rate found for region %s", region), errors.ErrorInternal)
	}

	return r.rates[region], nil
}
