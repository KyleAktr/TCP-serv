package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go [server|client <server_ip>]")
		return
	}

	if os.Args[1] == "server" {
		StartServer()
	} else if os.Args[1] == "client" && len(os.Args) == 3 {
		StartClient()
	} else {
		fmt.Println("Invalid arguments. Usage: go run main.go [server|client <server_ip>]")
	}
}
