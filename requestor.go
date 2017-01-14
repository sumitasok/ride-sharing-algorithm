package ride

import(
	"time"
)

type travelState string

const (
	pickedUp         travelState = "PICKED_UP"
	rideRequested    travelState = "RIDE_REQUESTED"
	dropped          travelState = "DROPPED"
	requestCancelled travelState = "REQUEST_CANCELLED"
)

func NewRequestor(id string, quantity int64, pickupTime time.Time, pickupLocation, dropLocation location) *requestor {
	return &requestor{
		Identifier:     id,
		State:          rideRequested,
		Quantity:       quantity,
		PickupTime:     pickupTime,
		PickupLocation: pickupLocation,
		DropLocation:   dropLocation,
		RequestTime:    time.Now(),
	}
}

type requestor struct {
	Identifier       string
	State            travelState
	Quantity         int64
	PickupTime       time.Time
	PickupLocation   location
	DropLocation     location
	ExpectedDropTime time.Time
	RequestTime      time.Time
}