package ride

func NewVehicle(capacity int64, location location) vehicle {
	return vehicle{
		Capacity: capacity,
		Location: location,
	}
}

func NewVehicleWithName(name string, capacity int64, location location) vehicle {
	return vehicle{
		Capacity: capacity,
		Location: location,
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
}

type location struct {
	Address string  `json:"address"`
	Long    float64 `json:"long"`
	Lat     float64 `json:"lat"`
}