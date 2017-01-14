package ride

func NewVehicle(capacity int64, location location) vehicle {
	return vehicle{
		Capacity: capacity,
		Location: location,
	}
}

func NewVehicleWithName(name string, capacity int64, location location) vehicle {
	return vehicle{
		ID: name,
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