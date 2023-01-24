package billing

import (
	"context"
	"fmt"
	"sync"

	"github.com/slaengkast/shipping-api/internal/errors"
)

type locationModel struct {
	code            string
	hasEUMembership bool
}

type inMemoryLocationStore struct {
	locations map[string]locationModel
	mtx       *sync.RWMutex
}

func NewInMemoryLocationStore() inMemoryLocationStore {
	return inMemoryLocationStore{locations: make(map[string]locationModel, 0), mtx: &sync.RWMutex{}}
}

func (r inMemoryLocationStore) GetByCode(_ context.Context, code string) (*location, error) {
	r.mtx.RLock()
	defer r.mtx.RUnlock()

	l, ok := r.locations[code]
	if !ok {
		return nil, errors.FromMessage(fmt.Sprintf("no location with code %s", code), errors.ErrorInput)
	}

	return unmarshalLocation(l)
}

func (r inMemoryLocationStore) AddLocation(_ context.Context, location *location) error {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	r.locations[location.GetCode()] = locationModel{location.code, location.hasEUMembership}

	return nil
}

func unmarshalLocation(l locationModel) (*location, error) {
	return NewLocation(l.code, l.hasEUMembership)
}
