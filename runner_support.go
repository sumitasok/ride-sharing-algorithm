package ride

import(
	"bufio"
	"fmt"
	"os"
	"strings"
	"strconv"
	"github.com/kr/pretty"
	"errors"
)

func AddVeh() (*vehicle, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter Vehicle name: ")
	text, _ := reader.ReadString('\n')
	name := chomp(text)

	reader = bufio.NewReader(os.Stdin)
	fmt.Print("Enter Vehicle capacity: ")
	text, _ = reader.ReadString('\n')

	capacity, _ := strconv.ParseInt(chomp(text), 10, 64)

	reader = bufio.NewReader(os.Stdin)
	fmt.Print("Enter Vehicle current location: ")
	text, _ = reader.ReadString('\n')

	address := chomp(text)

	loc := &location{}
	var err error

	for {
		loc, err = NewLocationFromAddress(address)

		if err == nil {
			break
		}
		println(err.Error())
	}

	v := NewVehicleWithName(name, capacity, *loc)


	redisST.AddVehicle("blr", name, v.Location.Long, v.Location.Lat)

	s, e := redisST.InsertVehicles(v)
	pretty.Print(s, e)

	return &v, nil
}

func PickupRider() error {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter Vehicle name: ")
	text, _ := reader.ReadString('\n')
	vehicle_name := chomp(text)

	reader = bufio.NewReader(os.Stdin)
	fmt.Print("Enter Passenger name: ")
	text, _ = reader.ReadString('\n')
	rider_name := chomp(text)

	vehicles, err := redisST.FetchVehicleDetail(vehicle_name)
	if err != nil {
		return err
	}

	if len(vehicles) == 0 {
		return errors.New("Vehicle not found")
	}

	v := vehicles[0]
	err = v.Pickup(rider_name)
	if err != nil {
		return err
	}

	_, err = redisST.InsertVehicles(v)
	if err != nil {
		return err
	}

	return nil
}

func RemoveVeh() (int64, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter Vehicle name: ")
	text, _ := reader.ReadString('\n')
	name := chomp(text)

	return redisST.RemoveVehicle("blr", name)
}

func chomp(command string) string {
	return strings.Trim(command, "\n")
}