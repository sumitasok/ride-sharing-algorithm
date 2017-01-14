package main

import(
	"bufio"
	"fmt"
	"os"
	ride "bitbucket.org/sumitasok/ride-fair"
	"strings"
)

var(
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
		case "RM_VEH":
		case "UPDT_VEH":
		case "REQ_RIDE":
		case "CNCL_RIDE":
		case "EXIT":
			os.Exit(0)
		default:
			fmt.Println(text)
		}
    }
	ridefair.Store()
}

func addVeh() (*ride.vehicle, error) {
	
}