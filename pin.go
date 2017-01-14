package ride

import (
	"errors"
	"fmt"
	"strconv"
	"time"
)

type tdMetrixToLatLongMap map[string]tdMetrix

type tdMetrix struct {
	Distance int           // meters
	Time     time.Duration // seconds
}

func NewPin(_nextState nextState, _rider requestor, _location location) *pin {
	return &pin{RextState: _nextState, Rider: _rider, Location: _location}
}

func NewPinFromVehicle(v_vehicle vehicle) *pin {
	return NewPin(start, requestor{}, v_vehicle.Location)
}

func NewPinFromRequestor(_requestor requestor, _nextState nextState) *pin {
	if _nextState == pickup {
		return NewPin(_nextState, _requestor, _requestor.PickupLocation)
	} else {
		return NewPin(_nextState, _requestor, _requestor.DropLocation)
	}
}

type pin struct {
	RextState     nextState
	Rider         requestor
	Location      location
	MetersToCover int64
	TimeToCover   time.Duration
	// metrixs after fetching from google
	Distance    tdMetrixToLatLongMap
	TimeToReach time.Time
}

func (p pin) isVehicle() bool {
	if p.RextState == start {
		return true
	}

	return false
}

func (p pin) latLongString() string {
	return fmt.Sprint(strconv.FormatFloat(p.Location.Lat, 'f', 6, 64) + "," + strconv.FormatFloat(p.Location.Long, 'f', 6, 64))
}

type nextState string

const (
	start  nextState = "START" // for vehicle
	drop   nextState = "DROP"
	pickup nextState = "PICK_UP"
)

func makePinList(pins ...pin) pinList {
	return pinList(pins)
}

func NewPinList() pinList {
	return pinList{}
}

type pinList []pin

func (p pinList) findByLatLongString(latlong string) (*pin, error) {
	for _, p_pin := range p {
		if p_pin.latLongString() == latlong {
			return &p_pin, nil
		}
	}

	return nil, errors.New("Pin not found")
}

func (p pinList) latLongList() []string {
	list := []string{}

	for _, _p := range p {
		list = append(list, _p.latLongString())
	}

	return list
}

func (p *pinList) count() int {
	return len(*p)
}

func (p *pinList) nextStateCount(state nextState) int {
	var count int = 0

	for _, _p := range *p {
		if _p.RextState == state {
			count += 1
		}
	}

	return count
}

func (p pinList) append(_pin pin) pinList {
	pins := []pin(p)
	pins = append(pins, _pin)

	x := makePinList(pins...)
	return x
}

func (p pinList) remove(_pin pin) pinList {
	pins := []pin{}

	for _, _p := range p {
		if !(_p.Rider.Identifier == _pin.Rider.Identifier && _p.RextState == _pin.RextState) {
			pins = append(pins, _p)
		}
	}

	return makePinList(pins...)
}

func (p pinList) valid() bool {
	state := map[string]nextState{}

	for _, pin := range p {
		if prevState, ok := state[pin.Rider.Identifier]; ok {
			if prevState == drop && pin.RextState == pickup {
				return false
			}
		}

		state[pin.Rider.Identifier] = pin.RextState
	}

	return true
}

func (p pinList) toString() string {
	text := ""

	for _, pin := range p {
		text = fmt.Sprintf("%s -> %s (%s)", text, string(pin.RextState), string(pin.Rider.Identifier))
	}

	return text
}

func (p pinList) toMapAPI() (string, error) {
	if len(p) < 2 {
		return "", errors.New("toMapAPI require two pins")
	}
	text := "https://www.google.co.in/maps/dir"
	for _, pin := range p {
		//text = fmt.Sprintf("%s/%s", text, url.QueryEscape(pin.Location.Address) )
		text = fmt.Sprintf("%s/%s", text, strconv.FormatFloat(pin.Location.Lat, 'f', 6, 64)+","+strconv.FormatFloat(pin.Location.Long, 'f', 6, 64))
	}
	return text, nil

}

// generate the individual pins for a vehicle

func generatePins(v vehicle, requestors ...requestor) []pin {
	pins := []pin{}

	// added riders who are on board to pins as to be dropped
	for _, r := range v.Riders {
		pins = append(pins, pin{
			RextState: drop,
			Rider:     *r,
		})
	}

	// listing pickups and drops of yet to be picked up riders
	for _, r := range requestors {
		pins = append(pins, pin{
			RextState: pickup,
			Rider:     r,
		}, pin{
			RextState: drop,
			Rider:     r,
		})

	}

	return pins
}
