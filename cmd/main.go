package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/slaengkast/shipping-api/internal/billing"
	"github.com/slaengkast/shipping-api/internal/booking"
	"github.com/slaengkast/shipping-api/internal/server"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

func main() {
	var (
		logLevel string
		port     int
	)

	app := &cli.App{
		Name:  "shipping-api",
		Usage: "Start a Sendify server",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "logLevel",
				Value:       "warn",
				Usage:       "Set the log level, valid values: debug, info, warn, error",
				Destination: &logLevel,
			},
			&cli.IntFlag{
				Name:        "port",
				Value:       80,
				Usage:       "Set the server port",
				Destination: &port,
			},
		},
		Action: func(ctx *cli.Context) error {
			configureLogging(logLevel)
			return run(port)
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Error().Err(err).Msg("")
	}
}

func configureLogging(logLevel string) {
	switch logLevel {
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
		log.Warn().Msg("Unknown logLevel " + logLevel)
	}
	log.Logger = log.With().Caller().Logger()
}

func run(port int) error {
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
			return err
		}

		if err := locationStore.AddLocation(context.Background(), location); err != nil {
			return err
		}
	}
	billingService := billing.NewService(rateStore, priceStore, locationStore)

	bookingStore := booking.NewInMemoryStore()
	bookingService := booking.NewService(bookingStore, billingService)

	bookingHandler := booking.NewHandler(bookingService)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	s := server.New(bookingHandler, port)
	go func() {
		if err := s.Run(); err != nil {
			log.Fatal().Err(err).Msg("")
		}
	}()
	<-ctx.Done()

	log.Info().Msg("Shutting down gracefully...")
	return nil
}
