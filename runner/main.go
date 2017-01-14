package main

import(
	"bufio"
	"fmt"
	"os"
	"bitbucket.org/sumitasok/ride-fair"
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
		fmt.Println(text)

		if chomp(text) == "EXIT" {
			break
		}
    }
	ridefair.Store()
}