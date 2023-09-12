package main

import (
	"fmt"
	"log"
	"net"
	"os"

	penguin "penguin/pkg"
)

func main() {
	if len(os.Args) != 2 && len(os.Args) != 1 {
		fmt.Println("[USAGE]: ./TCPChat $port")
		return
	}

	port := "8989" // Default port
	if len(os.Args) == 2 {
		port = os.Args[1]
	}

	fmt.Println("Server starting on port:", port)

	// Start TCP server
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Error starting the server: %s", err.Error())
	}
	defer listener.Close()

	// broadcast messages
	// go penguin.BroadcastMessages()

	for {
		connection, err := listener.Accept()
		if err != nil {
			log.Fatalf("Error accepting connection: %s", err.Error())
		}
		go penguin.HandleClient(connection)
	}
}
