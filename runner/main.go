package main

import(
	"bufio"
	"fmt"
	"os"
	ride "bitbucket.org/z_team_gojek/ride-fair"
	"strings"
)

func chomp(command string) string {
	return strings.Trim(command, "\n")
}

func main() {
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter text: ")
		text, _ := reader.ReadString('\n')

		switch(chomp(text)) {
		case "ADD_VEH":
			ride.AddVeh()
		case "RM_VEH":
			ride.RemoveVeh()
		case "UPDT_VEH":
		case "REQ_RIDE":
		case "CNCL_RIDE":
		case "PICKUP":
			// "vehicle_name", "rider_id"
			err := ride.PickupRider()
			fmt.Println("ERR: ", err.Error())
		case "EXIT":
			os.Exit(0)
		default:
			fmt.Println(text)
		}
    }
}