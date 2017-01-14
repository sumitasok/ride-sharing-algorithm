package ride

import (
	"errors"
)

var (
	ErrRiderAlreadyExist = errors.New("Rider already exist")
)

func NewVehicle(capacity int64, location location) vehicle {
	return vehicle{
		Capacity: capacity,
		Location: location,
		Requestors: map[string]*requestor{},
		Riders: map[string]*requestor{},
	}
}

func NewVehicleWithName(name string, capacity int64, location location) vehicle {
	return vehicle{
		ID: name,
		Capacity: capacity,
		Location: location,
		Requestors: map[string]*requestor{},
		Riders: map[string]*requestor{},
	}
}

type state struct {
	Active   bool
	Location location
}

type vehicle struct {
	State                state
	Capacity             int64
	Location             location
	ID                   string
	// requestors always sorted by first one to drop.
	Requestors           map[string]*requestor
	Riders               map[string]*requestor
}

func (v *vehicle) addRequestor(r requestor) error {
	if _, ok := v.Riders[r.Identifier]; ok {
		return ErrRiderAlreadyExist
	}

	v.Riders[r.Identifier] = &r
	return nil
}

func (v vehicle) occupancyStatus() int64 {
	count := int64(0)
	for _, r := range v.Requestors {
		count += r.Quantity
	}

	for _, r := range v.Riders {
		count += r.Quantity
	}

	return count
}

func (v vehicle) remainingOccupancy() int64 {
	return v.Capacity - v.occupancyStatus()
}