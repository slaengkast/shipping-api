package booking

import (
	"errors"
)

type booking struct {
	id          string
	origin      string
	destination string
	weight      float32
	price       float32
}

func NewBooking(id string, origin, destination string, weight float32, price float32) (*booking, error) {
	if id == "" {
		return nil, errors.New("id is empty")
	}
	if origin == "" {
		return nil, errors.New("origin is empty")
	}
	if destination == "" {
		return nil, errors.New("destination is empty")
	}
	if weight <= 0 {
		return nil, errors.New("invalid weight")
	}
	if price < 0 {
		return nil, errors.New("invalid price")
	}

	return &booking{
		id:          id,
		origin:      origin,
		destination: destination,
		weight:      weight,
		price:       price,
	}, nil
}

func (s *booking) Id() string {
	return s.id
}

func (s *booking) Origin() string {
	return s.origin
}

func (s *booking) Destination() string {
	return s.destination
}

func (s *booking) Weight() float32 {
	return s.weight
}

func (s *booking) Price() float32 {
	return s.price
}
