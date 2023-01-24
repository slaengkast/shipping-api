package booking

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

type billingReturn struct {
	price float32
	err   error
}

type storeReturn struct {
	sh  *booking
	err error
}

var (
	successfulBilling = billingReturn{50, nil}
	errorBilling      = billingReturn{0, errors.New("billing error")}
	successfulStore   = storeReturn{&booking{id: "test-id"}, nil}
	errorStore        = storeReturn{nil, errors.New("store error")}
)

func TestBookShipping(t *testing.T) {
	testCases := []struct {
		name          string
		origin        string
		destination   string
		weight        float32
		billingReturn billingReturn
		storeReturn   storeReturn
		shouldFail    bool
	}{
		{
			name:          "good booking",
			origin:        "SE",
			destination:   "SE",
			billingReturn: successfulBilling,
			storeReturn:   successfulStore,
			shouldFail:    false,
		},
		{
			name:          "bad origin",
			origin:        "",
			destination:   "DK",
			billingReturn: successfulBilling,
			storeReturn:   successfulStore,
			shouldFail:    true,
		},
		{
			name:          "bad destination",
			origin:        "SE",
			destination:   "",
			billingReturn: successfulBilling,
			storeReturn:   successfulStore,
			shouldFail:    true,
		},
		{
			name:          "billing error",
			origin:        "SE",
			destination:   "SE",
			billingReturn: errorBilling,
			storeReturn:   successfulStore,
			shouldFail:    true,
		},
		{
			name:          "store error",
			origin:        "SE",
			destination:   "SE",
			billingReturn: successfulBilling,
			storeReturn:   errorStore,
			shouldFail:    true,
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			bundle := newTestBundle()
			bundle.billingService.price = tc.billingReturn.price
			bundle.billingService.err = tc.billingReturn.err
			bundle.store.sh = tc.storeReturn.sh
			bundle.store.err = tc.storeReturn.err
			id, err := bundle.service.BookShipping(
				context.Background(),
				tc.origin,
				tc.destination,
				10,
			)
			if tc.shouldFail {
				require.NotNil(t, err, "expected an error, got nil")
				return
			}

			require.Nilf(t, err, "unexpected error")
			require.NotEqual(t, "", id, "expected bookingId to not be empty")
		})
	}
}

func TestGetBooking(t *testing.T) {
	testCases := []struct {
		name        string
		storeReturn storeReturn
		shouldFail  bool
	}{
		{
			name:        "good booking",
			storeReturn: successfulStore,
			shouldFail:  false,
		},
		{
			name:        "store error",
			storeReturn: errorStore,
			shouldFail:  true,
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			bundle := newTestBundle()
			bundle.store.sh = tc.storeReturn.sh
			bundle.store.err = tc.storeReturn.err
			actual, err := bundle.service.GetBooking(
				context.Background(),
				"mock-id",
			)
			if tc.shouldFail {
				require.NotNil(t, err, "expected an error, got nil")
				return
			}

			require.Nilf(t, err, "unexpected error")
			require.Equal(t, tc.storeReturn.sh.id, actual.id)
		})
	}
}

type billingServiceMock struct {
	price float32
	err   error
}

func (s billingServiceMock) CalculateShippingCost(_ context.Context, origin, destination string, weight float32) (float32, error) {
	return s.price, s.err
}

type storeMock struct {
	sh  *booking
	err error
}

func (r storeMock) AddBooking(_ context.Context, sh *booking) error {
	return r.err
}

func (r storeMock) GetBooking(_ context.Context, id string) (*booking, error) {
	return r.sh, r.err
}

type bundle struct {
	service        Service
	store          *storeMock
	billingService *billingServiceMock
}

func newTestBundle() bundle {
	store := &storeMock{}
	billingService := &billingServiceMock{}
	return bundle{
		service:        NewService(store, billingService),
		store:          store,
		billingService: billingService,
	}
}
