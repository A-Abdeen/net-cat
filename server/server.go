package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	penguin "penguin/pkg"
)

func main() {
	// check if correct number of arguments is used
	if len(os.Args) != 2 && len(os.Args) != 1 {
		fmt.Println("[USAGE]: ./TCPChat $port")
		return
	}

	port := "8989" // Default port
	if len(os.Args) == 2 {
		port = os.Args[1]
	}

	// A channel to handle the shutdown signal
	shutdownSignal := make(chan os.Signal, 1)
	signal.Notify(shutdownSignal, syscall.SIGINT, syscall.SIGTERM)

	fmt.Printf("Listening on the IP: %s and port :%s", penguin.GetLOcalIp(), port)

	// Listen for the shutdown signal
	go func() {
		<-shutdownSignal
		// Perform any cleanup or shutdown tasks here
		fmt.Println("\nclosing server")
		for _, client := range penguin.Clients {
			client.Socket.Write([]byte("\nServer is closed"))
		}
		penguin.AllMessages = append(penguin.AllMessages, "Server is closed\n")
		// For example, close network connections, save data, etc.
		// Save all messages to a file before exiting
		penguin.SaveAllMessagesToFile("all_chat_messages.txt")
		// Then exit the program gracefully
		os.Exit(0)
	}()

	// Start TCP server
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Error starting the server: %s", err.Error())
	}
	defer listener.Close()

	// broadcast messages
	go penguin.BroadcastMessages()

	for {
		conn, err := listener.Accept()
		if penguin.UserCounter > penguin.MaxUsers {
			// Server is full, reject the connection
			
			if err != nil {
				log.Printf("Error accepting connection: %s", err.Error())
			} else {
				conn.Write([]byte("Server is full. Please try again later."))
				conn.Close()
			}
		} else {
			// Accept new connections until the user limit is reached
			if err != nil {
				log.Printf("Error accepting connection: %s", err.Error())
			} else {
				go penguin.HandleClient(conn)
				penguin.UserCounter++
			}
		}
	}
}
