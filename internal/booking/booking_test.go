package booking

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewBooking(t *testing.T) {
	testCases := []struct {
		name        string
		id          string
		origin      string
		destination string
		weight      float32
		price       float32
		shouldFail  bool
	}{
		{
			name:        "valid booking",
			id:          "test-id",
			origin:      "SE",
			destination: "DK",
			weight:      300,
			price:       300,
		},
		{
			name:        "missing id",
			id:          "",
			origin:      "SE",
			destination: "DK",
			weight:      300,
			price:       300,
			shouldFail:  true,
		},
		{
			name:        "missing origin",
			id:          "test-id",
			origin:      "",
			destination: "DK",
			weight:      300,
			price:       300,
			shouldFail:  true,
		},
		{
			name:        "missing destination",
			id:          "test-id",
			origin:      "SE",
			destination: "",
			weight:      300,
			price:       300,
			shouldFail:  true,
		},
		{
			name:        "bad weight",
			id:          "test-id",
			origin:      "SE",
			destination: "DK",
			weight:      0,
			price:       300,
			shouldFail:  true,
		},
		{
			name:        "bad price",
			id:          "test-id",
			origin:      "SE",
			destination: "DK",
			weight:      300,
			price:       -10,
			shouldFail:  true,
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			_, err := NewBooking(tc.id, tc.origin, tc.destination, tc.weight, tc.price)

			if tc.shouldFail {
				require.NotNil(t, err, "expected an error, got nil")
				return
			}

			require.Nilf(t, err, "unexpected error")
		})
	}
}
