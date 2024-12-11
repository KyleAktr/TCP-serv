package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func StartClient() {
	conn, err := net.Dial("tcp", "127.0.0.1:3569")
	if err != nil {
		fmt.Println("connexion failed")
		return
	}
	defer conn.Close()

	//name of client as the server do need it
	fmt.Print("Enter your name: ")
	name, _ := bufio.NewReader(os.Stdin).ReadString('\n')

	fmt.Fprint(conn, name) //the configuration from server

	go ReadMessage(conn)
	writeMessage(conn)

}

func ReadMessage(conn net.Conn) {

	for {
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			(fmt.Println("Not on the server"))
			return //leaving immediately the loop in case of error
		}
		fmt.Print(message)
	}
}
func writeMessage(conn net.Conn) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() { //allow to read continusely the terminal
		message := scanner.Text()
		fmt.Fprintln(conn, message)
	}
}
