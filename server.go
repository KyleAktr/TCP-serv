package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
)

func handleServerError(err error) {
	if err != nil {
		panic(err)
	}
}

const SERVER_PORT = "3569"

type Client struct {
	name string
	conn net.Conn
}

var clientNameMap = make(map[net.Conn]*Client) // Global gestion for clients
var clientMutex sync.Mutex                     // Mutex to synchronize access to the counter

// StartServer continuously listens for connections and launches a goroutine for each client that connects, allowing the server to handle multiple clients simultaneously
func StartServer() {
	ln, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", SERVER_PORT))
	handleServerError(err)
	defer ln.Close()

	fmt.Println("Server is listening on port", SERVER_PORT)

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		// Handle each client in a separate goroutine
		go AssignClientName(conn)
	}
}

// assignClientName increments the counter and returns a unique name for the client
func AssignClientName(conn net.Conn) {

	fmt.Fprint(conn, "Enter your name") //formate to be able to link the adress of the client logged with the message
	name, _ := bufio.NewReader(conn).ReadString('\n')
	name = strings.TrimSpace(name)
	clients := &Client{name: name, conn: conn}

	clientMutex.Lock()                                                        //lock the access to the function the time that it is used
	clientNameMap[conn] = clients                                             // set the Client struct value
	clientMutex.Unlock()                                                      // make the function available again
	fmt.Printf("%s has connected from %s\n", clients.name, conn.RemoteAddr()) //remoteAddr allow to write the located data and not the location
	Broadcast(fmt.Sprintf("%s: has joined the chat\n", clients.name), conn)

	HandleClient(conn, clients.name)

}

func Broadcast(message string, sender net.Conn) {
	fmt.Print(message) //print the message on the server to help with gestion

	clientMutex.Lock()
	defer clientMutex.Unlock()                //avoir race condition during the send of the message
	for conn, client := range clientNameMap { //check for every client around
		if conn != sender { // check if the client is the sender himself
			_, err := fmt.Fprintln(client.conn, message) //print the message to the client
			if err != nil {
				fmt.Println("Sending error:", client.name, err) //we're avoiding panic because it would shut down the server which is something we definitly won't
				conn.Close()
				delete(clientNameMap, conn)
			}

		}
	}
}

func HandleClient(conn net.Conn, clientName string) {
	defer func() { //we're using defer function to clear all the information to not miss any and avoid extra line, it's as any function in itself
		clientMutex.Lock()
		delete(clientNameMap, conn) //withdrawing client's name from NameMap
		clientMutex.Unlock()
		conn.Close() //assure that the connexion is closed
		fmt.Printf("%s has left the server\n", clientName)
	}()

	scanner := bufio.NewScanner(conn) //checking what the client is writing
	for scanner.Scan() {
		message := scanner.Text()                                      //reading each client's line
		Broadcast(fmt.Sprintf("%s : %s\n", clientName, message), conn) //print formated message
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Can't read the message\n", clientName)

	}
}
