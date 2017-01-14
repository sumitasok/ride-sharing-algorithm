package ride

import (
	//"github.com/kr/pretty"
	"github.com/stretchr/testify/assert"
	"testing"
	//"time"
	"time"
	"github.com/kr/pretty"
)

/*func TestGDistanceMatrix(t *testing.T) {
	assert := assert.New(t)

	t_requestor := *NewRequestor("rider-2", 1, time.Now(),
		*NewLocationFromLatLong(12.975928, 77.638986),
		*NewLocationFromLatLong(12.959969, 77.641068))

	t2_requestor := *NewRequestor("rider-1", 1, time.Now(),
		*NewLocationFromLatLong(12.975201, 77.638986),
		*NewLocationFromLatLong(12.959401, 77.641068))

	t_vehicle := NewVehicle(4, *NewLocationFromLatLong(12.961543, 77.644396))
	t_rider_pins := makePinList(
		// first pin is always the cars current position
		*NewPinFromRequestor(t_requestor, pickup),
		*NewPinFromRequestor(t_requestor, drop),
		*NewPinFromRequestor(t2_requestor, drop),
	)

	pins := GDistanceMatrix(
		t_vehicle,
		t_rider_pins,
	)

	pretty.Println(pins)

	t_vehicle_pin := NewPinFromVehicle(t_vehicle)

	t_vehicle_rider_pins := addVehiclePinWithRider(*t_vehicle_pin, t_rider_pins)
	assert.Equal(t_vehicle_rider_pins.count(), 4)

	x, err := t_vehicle_rider_pins.findByLatLongString(t_vehicle_pin.latLongString())
	assert.NotNil(x)
	assert.NoError(err)

	assert.True(true)

}*/

/*func TestDistanceMatRoute(t *testing.T) {
	v := vehicle{
		capacity: 4,
		location: location{
			lat: 12.983710,		//swami vivekananda metro
			long: 77.640724,
		},
		riders: map[string]*requestor{
			"rider-1": &requestor{
				identifier: "rider-1",
				state:      pickedUp,
				quantity:   1,
				dropLocation: location{
					lat : 12.958343,	//Murgeshpalya
					long: 77.666473,
				},
				pickupTime: time.Now().Add(-time.Minute*30),
				directDropTime:time.Now().Add(time.Minute*17),

			},
			"rider-2": &requestor{
				identifier: "rider-2",
				state:      pickedUp,
				quantity:   1,
				dropLocation: location{
					lat : 12.956837,
					long: 77.701149,
				},
				pickupTime: time.Now().Add(-time.Minute*30),
				directDropTime:time.Now().Add(time.Minute*30),
			},
		},
	}

	req := requestor{
		identifier: "rider-2",
		state: rideRequested,
		quantity: 1,
		pickupLocation: location{
			lat : 12.967663, 	//pickup jevan bheemanagar
			long: 77.656775,
		},
		dropLocation: location{
			lat : 12.956503,	//drop marathalli
			long: 77.700634,
		},
	}

	_, directDropTime := calculateDeviation(v, req)

	//pretty.Println("TEST PRINT::::", routesCalculated)


	v = vehicle{
		capacity: 4,
		location: location{
			lat: 12.962693,		//swami vivekananda metro
			long: 77.664413,
		},
		riders: map[string]*requestor{
			"rider-1": &requestor{
				identifier: "rider-1",
				state:      pickedUp,
				quantity:   1,
				dropLocation: location{
					lat : 12.958343,	//Murgeshpalya
					long: 77.666473,
				},
				pickupTime: time.Now().Add(-time.Minute*30),
				directDropTime:time.Now().Add(time.Minute*17),

			},
			"rider-2": &requestor{
				identifier: "rider-2",
				state: rideRequested,
				quantity: 1,
				dropLocation: location{
					lat : 12.956503,	//drop marathalli
					long: 77.700634,
				},
				directDropTime: directDropTime,
			},
		},
	}

	req = requestor{
		identifier: "rider-3",
		state: rideRequested,
		quantity: 1,
		pickupLocation: location{
			lat : 12.959012, 	//pickup jevan bheemanagar
			long: 77.691193,
		},
		dropLocation: location{
			lat : 12.957005,	//drop marathalli
			long: 77.744923,
		},
	}

	 calculateDeviation(v, req)


	//pretty.Println("TEST PRINT::::", routesCalculated2)
}*/


func TestDistanceMatRouteKFC (t *testing.T) {
	assert := assert.New(t)

	carCurrLoc, _ := NewLocationFromAddress("Ganapathi Temple 18th Rd, 6th Block, Koramangala, Bengaluru, Karnataka 560030")
	Rider1Drop, _ := NewLocationFromAddress("2037, 1st Cross Rd, Kodihalli, Bengaluru, Karnataka 560008")
	Rider2PickUP, _ := NewLocationFromAddress("41, Srinivagilu Main Rd, Ejipura, Bengaluru, Karnataka 560007")
	Rider2Drop, _ := NewLocationFromAddress("IBM Ln, Embassy Golf Links Business Park, Challaghatta, Bengaluru, Karnataka 560071")
	Rider3Pickup, err := NewLocationFromAddress("1st Block Koramangala, Koramangala, Bengaluru, Karnataka 560034")
	Rider3Drop, _ := NewLocationFromAddress("Konen Agrahara, Konena Agrahara, Murgesh Pallya, Bengaluru, Karnataka 560017")

	assert.NoError(err)
	pretty.Println(Rider3Pickup)

	carCurrLoc2, _ := NewLocationFromAddress("No A, Floor,, 445, 18th Main Rd, 1st Stage, Koramangala, Bengaluru, Karnataka 560095")

	v := vehicle{
		Capacity: 4,
		Location: *carCurrLoc,
		Riders: map[string]*requestor{
			"rider-1": &requestor{
				Identifier: "rider-1",
				State:      pickedUp,
				Quantity:   1,
				DropLocation: *Rider1Drop,
				PickupTime: time.Now().Add(-time.Minute*7),
				DirectDropTime:time.Now().Add(time.Minute*16),

			},
		},
		ExpectedLastDropTime: time.Now().Add(time.Minute*25),
	}

	req := requestor{
		Identifier: "rider-2",
		State: rideRequested,
		Quantity: 1,
		PickupLocation: *Rider2PickUP,
		DropLocation: *Rider2Drop,
	}

	reqPickUpPin :=  *NewPinFromRequestor(req, pickup) 	// New pin for upcoming rider's pickup
	reqDropPin :=  *NewPinFromRequestor(req, drop)		// New pin for upcoming rider's drop

	out := make(chan DeviationResult)

	var directDropTime time.Time
	go func() {
		_, directDropTime, _ = calculateDeviation(v, reqPickUpPin, reqDropPin, time.Now(), out)
	}()
	//pretty.Println("TEST PRINT::::", routesCalculated)
	d := <- out
	pretty.Println("Go-JEK:::::",d.Route)


	v = vehicle{
		Capacity: 4,
		Location: *carCurrLoc2,
		Riders: map[string]*requestor {
			"rider-1": &requestor{
				Identifier: "rider-1",
				State:      pickedUp,
				Quantity:   1,
				DropLocation: *Rider1Drop,
				PickupTime: time.Now().Add(-time.Minute*30),
				DirectDropTime:time.Now().Add(time.Minute*25),

			},
		},
		Requestors: map[string]*requestor{
			"rider-2": &requestor{
				Identifier: "rider-2",
				State: rideRequested,
				Quantity: 1,
				PickupLocation: *Rider2PickUP,
				DropLocation: *Rider2Drop,
				DirectDropTime: directDropTime,
			},
		},
		ExpectedLastDropTime: d.ExpectedLastTime,
	}

	req = requestor{
		Identifier: "rider-3",
		State: rideRequested,
		Quantity: 1,
		PickupLocation: *Rider3Pickup,
		DropLocation: *Rider3Drop,
	}

	reqPickUpPin =  *NewPinFromRequestor(req, pickup) 	// New pin for upcoming rider's pickup
	reqDropPin =  *NewPinFromRequestor(req, drop)		// New pin for upcoming rider's drop
	go func() {
		calculateDeviation(v, reqPickUpPin, reqDropPin, time.Now().Add(time.Minute * 2), out) // time.Now() * 4 is the time between cars current location and Rider-2's pickup
	}()
	route := <- out
	// pretty.Println("Out Channel", <-out)

	assert.Equal(" -> START () -> PICK_UP (rider-2) -> DROP (rider-1) -> DROP (rider-2) -> PICK_UP (rider-3) -> DROP (rider-3)", route.Route.toString())

	close(out)
	//pretty.Println("TEST PRINT::::", routesCalculated2)
}

func PendingDistanceMatRouteKFCMultipleVehicle(t *testing.T) {

	carCurrLoc, _ := NewLocationFromAddress("Ganapathi Temple 18th Rd, 6th Block, Koramangala, Bengaluru, Karnataka 560030")
	Rider1Drop, _ := NewLocationFromAddress("2037, 1st Cross Rd, Kodihalli, Bengaluru, Karnataka 560008")
	Rider2PickUP, _ := NewLocationFromAddress("41, Srinivagilu Main Rd, Ejipura, Bengaluru, Karnataka 560007")
	Rider2Drop, _ := NewLocationFromAddress("IBM Ln, Embassy Golf Links Business Park, Challaghatta, Bengaluru, Karnataka 560071")

	v := vehicle{
		Capacity: 4,
		Location: *carCurrLoc,
		Riders: map[string]*requestor{
			"rider-1": &requestor{
				Identifier: "rider-1",
				State:      pickedUp,
				Quantity:   1,
				DropLocation: *Rider1Drop,
				PickupTime: time.Now().Add(-time.Minute*7),
				DirectDropTime:time.Now().Add(time.Minute*16),

			},
		},
		ExpectedLastDropTime: time.Now().Add(time.Minute*25),
	}

	// v2 := vehicle{
	// 	Capacity: 4,
	// 	Location: *carCurrLoc,
	// 	Riders: map[string]*requestor{
	// 		"rider-1": &requestor{
	// 			Identifier: "rider-1",
	// 			State:      pickedUp,
	// 			Quantity:   1,
	// 			DropLocation: *Rider1Drop,
	// 			PickupTime: time.Now().Add(-time.Minute*7),
	// 			DirectDropTime:time.Now().Add(time.Minute*16),

	// 		},
	// 	},
	// 	ExpectedLastDropTime: time.Now().Add(time.Minute*25),
	// }

	req := requestor{
		Identifier: "rider-2",
		State: rideRequested,
		Quantity: 1,
		PickupLocation: *Rider2PickUP,
		DropLocation: *Rider2Drop,
	}

	reqPickUpPin :=  *NewPinFromRequestor(req, pickup) 	// New pin for upcoming rider's pickup
	reqDropPin :=  *NewPinFromRequestor(req, drop)		// New pin for upcoming rider's drop

	out := make(chan DeviationResult)

	var directDropTime time.Time
	go func() {
		_, directDropTime, _ = calculateDeviation(v, reqPickUpPin, reqDropPin, time.Now(), out)
	}()
	//pretty.Println("TEST PRINT::::", routesCalculated)
	d := <- out
	pretty.Println("Go-JEK:::::",d.Route)
}
