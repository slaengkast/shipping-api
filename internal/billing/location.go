package billing

import (
	"github.com/slaengkast/shipping-api/internal/errors"
)

type location struct {
	code            string
	hasEUMembership bool
}

func NewLocation(code string, hasEUMembership bool) (*location, error) {
	if code == "" {
		return nil, errors.FromMessage("empty code", errors.ErrorInput)
	}

	return &location{code: code, hasEUMembership: hasEUMembership}, nil
}

func (l location) IsMemberOfEU() bool {
	return l.hasEUMembership
}

func (l location) GetCode() string {
	return l.code
}

func getRegion(origin, destination *location) string {
	switch {
	case origin.GetCode() == destination.GetCode():
		return "domestic"
	case origin.IsMemberOfEU() && destination.IsMemberOfEU():
		return "eu"
	default:
		return "international"
	}
}
