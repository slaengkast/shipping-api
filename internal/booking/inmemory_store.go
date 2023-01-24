package booking

import (
	"context"
	"fmt"
	"sync"

	"github.com/slaengkast/shipping-api/internal/errors"
)

type bookingModel struct {
	id          string
	origin      string
	destination string
	weight      float32
	price       float32
	currency    string
}

type inMemoryStore struct {
	bookings map[string]bookingModel
	mtx      *sync.RWMutex
}

func NewInMemoryStore() inMemoryStore {
	return inMemoryStore{bookings: make(map[string]bookingModel, 0), mtx: &sync.RWMutex{}}
}

func (r inMemoryStore) GetBooking(_ context.Context, id string) (*booking, error) {
	r.mtx.RLock()
	defer r.mtx.RUnlock()

	bookingModel, ok := r.bookings[id]
	if !ok {
		return nil, errors.FromMessage(fmt.Sprintf("booking %s not found", id), errors.ErrorNotFound)
	}
	return unmarshalBooking(bookingModel)
}

func (r inMemoryStore) AddBooking(ctx context.Context, sh *booking) error {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	if _, ok := r.bookings[sh.Id()]; ok {
		return errors.FromMessage("booking already exists", errors.ErrorConflict)
	}

	r.bookings[sh.Id()] = marshalBooking(sh)
	return nil
}

func unmarshalBooking(bookingModel bookingModel) (*booking, error) {
	return NewBooking(
		bookingModel.id,
		bookingModel.origin,
		bookingModel.destination,
		bookingModel.weight,
		bookingModel.price,
	)
}

func marshalBooking(b *booking) bookingModel {
	return bookingModel{
		id:          b.id,
		origin:      b.origin,
		destination: b.destination,
		weight:      b.weight,
		price:       b.price,
	}
}
