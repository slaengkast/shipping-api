package billing

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewLocation(t *testing.T) {
	testCases := []struct {
		name            string
		code            string
		hasEUMembership bool
		shouldFail      bool
	}{
		{
			name:            "valid location",
			code:            "SE",
			hasEUMembership: true,
			shouldFail:      false,
		},
		{
			name:            "empty code",
			code:            "",
			hasEUMembership: true,
			shouldFail:      true,
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			_, err := NewLocation(tc.code, tc.hasEUMembership)

			if tc.shouldFail {
				require.NotNil(t, err, "expected an error, got nil")
				return
			}

			require.Nilf(t, err, "unexpected error")
		})
	}
}

func Test(t *testing.T) {
	testCases := []struct {
		name          string
		origin        location
		destination   location
		expectedClass string
	}{
		{
			name:          "domestic",
			origin:        location{code: "SE", hasEUMembership: true},
			destination:   location{code: "SE", hasEUMembership: true},
			expectedClass: "domestic",
		},
		{
			name:          "eu",
			origin:        location{code: "SE", hasEUMembership: true},
			destination:   location{code: "DK", hasEUMembership: true},
			expectedClass: "eu",
		},
		{
			name:          "international",
			origin:        location{code: "SE", hasEUMembership: true},
			destination:   location{code: "US", hasEUMembership: false},
			expectedClass: "international",
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tc.expectedClass, getRegion(&tc.origin, &tc.destination))
		})
	}
}
