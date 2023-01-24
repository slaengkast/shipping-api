package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/slaengkast/shipping-api/internal/billing"
	"github.com/slaengkast/shipping-api/internal/booking"

	"github.com/stretchr/testify/require"
)

const (
	port    = 8080
	address = "localhost"
)

func TestBookShipping(t *testing.T) {
	t.Parallel()

	client := NewClient(fmt.Sprintf("http://%s:%d", address, port))

	id, err := client.BookShipping("SE", "DK", 400)
	require.Nil(t, err)
	require.NotEqual(t, "", id)
}

func TestGetBooking(t *testing.T) {
	t.Parallel()

	client := NewClient(fmt.Sprintf("http://%s:%d", address, port))

	id, err := client.BookShipping("SE", "DK", 400)
	require.Nil(t, err)
	require.NotEqual(t, "", id)

	booking, err := client.GetBooking(id)
	require.InDelta(t, 400, booking["weight"], 1e-9)
	require.Equal(t, "SE", booking["origin"])
	require.Equal(t, "DK", booking["destination"])
	require.InDelta(t, 3000, booking["price"], 1e-9)
}

func TestHealth(t *testing.T) {
	t.Parallel()

	client := NewClient(fmt.Sprintf("http://%s:%d", address, port))

	status, err := client.Health()
	require.Nil(t, err)
	require.Equal(t, "healthy", status)
}

func startHttp() error {
	rateStore := billing.NewInMemoryRateStore(
		map[string]float32{
			"domestic":      1.0,
			"eu":            1.5,
			"international": 2.5,
		},
	)
	priceStore := billing.NewInMemoryPriceStore(
		map[string]float32{
			"small":  100,
			"medium": 300,
			"large":  500,
			"huge":   2000,
		},
	)
	locations := []struct {
		code            string
		hasEuMembership bool
	}{
		{"SE", true},
		{"DK", true},
		{"DE", true},
		{"US", false},
		{"UG", false},
	}

	locationStore := billing.NewInMemoryLocationStore()
	for _, l := range locations {
		location, err := billing.NewLocation(l.code, l.hasEuMembership)
		if err != nil {
			panic(err)
		}

		if err := locationStore.AddLocation(context.Background(), location); err != nil {
			panic(err)
		}
	}
	billingService := billing.NewService(rateStore, priceStore, locationStore)

	bookingStore := booking.NewInMemoryStore()
	bookingService := booking.NewService(bookingStore, billingService)

	bookingHandler := booking.NewHandler(bookingService)
	s := New(bookingHandler, port)
	go func() {
		if err := s.Run(); err != nil {
			panic(err.Error())
		}
	}()
	return waitUntilReady()
}

func waitUntilReady() error {
	waitChan := make(chan int, 0)

	go func() {
		for {
			res, err := http.Get(fmt.Sprintf("http://%s:%d/health", address, port))
			if err != nil {
				time.Sleep(time.Second)
				continue
			}

			if res.StatusCode == http.StatusOK {
				waitChan <- 1
				return
			}
		}
	}()
	deadline := time.After(10 * time.Second)
	select {
	case <-waitChan:
		return nil
	case <-deadline:
		return errors.New("server not ready")
	}
}

func TestMain(m *testing.M) {
	if err := startHttp(); err != nil {
		fmt.Print(err.Error())
	}
	os.Exit(m.Run())
}
