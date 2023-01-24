package billing

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

type locationReturn struct {
	location *location
	err      error
}

type priceReturn struct {
	price float32
	err   error
}

type rateReturn struct {
	rate float32
	err  error
}

var (
	successfulLocation = locationReturn{&location{}, nil}
	errorLocation      = locationReturn{nil, errors.New("error location")}
	successfulPrice    = priceReturn{100, nil}
	errorPrice         = priceReturn{0, errors.New("error price")}
	successfulRate     = rateReturn{2.0, nil}
	errorRate          = rateReturn{0, errors.New("error rate")}
)

func TestBookShipping(t *testing.T) {
	testCases := []struct {
		name           string
		origin         string
		destination    string
		weight         float32
		locationReturn locationReturn
		priceReturn    priceReturn
		rateReturn     rateReturn
		shouldFail     bool
	}{
		{
			name:           "good booking",
			origin:         "SE",
			destination:    "SE",
			weight:         100,
			rateReturn:     successfulRate,
			priceReturn:    successfulPrice,
			locationReturn: successfulLocation,
			shouldFail:     false,
		},
		{
			name:           "bad origin",
			origin:         "",
			destination:    "SE",
			weight:         100,
			rateReturn:     successfulRate,
			priceReturn:    successfulPrice,
			locationReturn: successfulLocation,
			shouldFail:     true,
		},
		{
			name:           "bad destination",
			origin:         "SE",
			destination:    "",
			weight:         100,
			rateReturn:     successfulRate,
			priceReturn:    successfulPrice,
			locationReturn: successfulLocation,
			shouldFail:     true,
		},
		{
			name:           "negative weight",
			origin:         "SE",
			destination:    "SE",
			weight:         -100,
			rateReturn:     successfulRate,
			priceReturn:    successfulPrice,
			locationReturn: successfulLocation,
			shouldFail:     true,
		},
		{
			name:           "invalid weight",
			origin:         "SE",
			destination:    "SE",
			weight:         1000000,
			rateReturn:     successfulRate,
			priceReturn:    successfulPrice,
			locationReturn: successfulLocation,
			shouldFail:     true,
		},
		{
			name:           "rate error",
			origin:         "SE",
			destination:    "SE",
			weight:         100,
			rateReturn:     errorRate,
			priceReturn:    successfulPrice,
			locationReturn: successfulLocation,
			shouldFail:     true,
		},
		{
			name:           "price error",
			origin:         "SE",
			destination:    "SE",
			weight:         100,
			rateReturn:     successfulRate,
			priceReturn:    errorPrice,
			locationReturn: successfulLocation,
			shouldFail:     true,
		},
		{
			name:           "location error",
			origin:         "SE",
			destination:    "SE",
			weight:         100,
			rateReturn:     successfulRate,
			priceReturn:    successfulPrice,
			locationReturn: errorLocation,
			shouldFail:     true,
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			bundle := newTestBundle()
			bundle.ratestore.rate = tc.rateReturn.rate
			bundle.ratestore.err = tc.rateReturn.err
			bundle.pricestore.price = tc.priceReturn.price
			bundle.pricestore.err = tc.priceReturn.err
			bundle.locationstore.location = tc.locationReturn.location
			bundle.locationstore.err = tc.locationReturn.err
			id, err := bundle.service.CalculateShippingCost(
				context.Background(),
				tc.origin,
				tc.destination,
				tc.weight,
			)
			if tc.shouldFail {
				require.NotNil(t, err, "expected an error, got nil")
				return
			}

			require.Nilf(t, err, "unexpected error")
			require.NotEqual(t, "", id, "expected id to not be empty")
		})
	}
}

type ratestoreMock struct {
	rate float32
	err  error
}

func (r ratestoreMock) GetRateByRegion(_ context.Context, region string) (float32, error) {
	return r.rate, r.err
}

type pricestoreMock struct {
	price float32
	err   error
}

func (r pricestoreMock) GetPriceByWeightClass(_ context.Context, class string) (float32, error) {
	return r.price, r.err
}

type locationstoreMock struct {
	location *location
	err      error
}

func (r locationstoreMock) GetByCode(_ context.Context, code string) (*location, error) {
	return r.location, r.err
}

type bundle struct {
	service       Service
	ratestore     *ratestoreMock
	pricestore    *pricestoreMock
	locationstore *locationstoreMock
}

func newTestBundle() bundle {
	locationstore := &locationstoreMock{}
	pricestore := &pricestoreMock{}
	ratestore := &ratestoreMock{}
	return bundle{
		service:       NewService(ratestore, pricestore, locationstore),
		ratestore:     ratestore,
		pricestore:    pricestore,
		locationstore: locationstore,
	}
}
