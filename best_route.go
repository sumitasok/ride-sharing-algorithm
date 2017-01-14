package ride

import (
	"fmt"
	"time"
	"math"
	"github.com/kr/pretty"
)

type DeviationResult struct {
	V                vehicle
	Route            pinList
	Deviation        time.Duration
	VehicleDeviation time.Duration
	ExpectedLastTime time.Time
}

func calculateDeviation(v vehicle, reqPickUpPin , reqDropPin pin, now time.Time, out chan DeviationResult) ( []pinList, time.Time, time.Duration ){

	riderPins := v.GetRiderReqPins()
	pins := pinList{}

	vehiclePin := NewPinFromVehicle(v)
	pins = riderPins.append(*vehiclePin)

	pins = append(pins, reqPickUpPin)
	pins = append(pins, reqDropPin)


	pinsWithMetrics := GDistanceMatrix(pins)

	vehiclePinMatrix, _ := pinsWithMetrics.findByLatLongString(vehiclePin.latLongString())	// vehiclePinMatrix gives metrics of vehicle that is all the distances covered having vehicle the start point

	//now := time.Now()

	requestPickUpPinMatrix, _ := pinsWithMetrics.findByLatLongString(reqPickUpPin.latLongString())

	// from vehiclePinMatrix, calculating the direct time to pick up the upcoming rider
	reqPickUpTime := now.Add(vehiclePinMatrix.Distance[reqPickUpPin.latLongString()].Time)

	// direct time from upcoming riders pickup to his drop
	reqBestDropTime := reqPickUpTime.Add(requestPickUpPinMatrix.Distance[reqDropPin.latLongString()].Time)

	// requestorTravel time if directly from pickup to drop
	reqBestTravelTime := reqBestDropTime.Sub(reqPickUpTime)

	// adding reqDropTime to the upcoming rider's pickup pin
	reqPickUpPin.Rider.DirectDropTime = reqBestDropTime

	// adding reqDropTime to the upcoming rider's drop pin
	reqDropPin.Rider.DirectDropTime = reqBestDropTime

	riderPins = append(riderPins, reqPickUpPin, reqDropPin)

	routes_calculated := []pinList{}

	for combination := range generateCombinations(riderPins, riderPins.count()) {
		fmt.Println(combination.toString()) // This is instead of process(combination)
		routes_calculated = append(routes_calculated, processEachPinWithMatrix(*vehiclePin, combination, pinsWithMetrics))
	}

	bestRouteDeviation := time.Duration(math.MaxInt64)
	bestRoute := pinList{}
	bestRouteIndex := 0
	vehicleDeviation := time.Duration(math.MaxInt64)
	var stepTime time.Time
	for pinID, pins := range routes_calculated {
		routeDeviation := time.Duration(0)
		stepTime = now
		for _, route := range pins {
			stepTime = stepTime.Add(route.TimeToCover)
			if route.RextState == drop {
				// Deviation for driver who is going to be dropped at this step
				dev := stepTime.Sub(route.Rider.DirectDropTime)
				routeDeviation += dev
			}
		}
		//pretty.Println("Route No :: ", pinID, "Route::", pins, "Deviation::", routeDeviation.Minutes())

		if routeDeviation < bestRouteDeviation {
			bestRouteDeviation = routeDeviation
			bestRoute = pins
			bestRouteIndex = pinID
			vehicleDeviation = stepTime.Sub(v.ExpectedLastDropTime)
			fmt.Println("stepTime",stepTime,"deltaDeviation", vehicleDeviation,"reqBestDropTime",reqBestDropTime,"v.expectedLastDropTime",v.ExpectedLastDropTime,"bestRouteDeviation",bestRouteDeviation,"reqBestTravelTime",reqBestTravelTime)
		}
	}

	pretty.Println("Best route index", bestRouteIndex, "Deviation::: ", bestRouteDeviation.Minutes(),"Delta Deviation:: ", vehicleDeviation.Minutes())

	out <- DeviationResult{
		v,
		bestRoute,
		bestRouteDeviation,
		vehicleDeviation,
		stepTime,
	}

	return routes_calculated, reqBestDropTime, vehicleDeviation
}
