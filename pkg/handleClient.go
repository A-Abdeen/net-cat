package penguin

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	//"time"
)

func HandleClient(connection net.Conn) {
	fmt.Println("New client connected:", connection.RemoteAddr().String())

	// Read welcome message from welcome.txt from the message folder
	welcomeMsg, err := readWelcomeMsg()
	if err != nil {
		log.Printf("Error reading welcome.txt: %s", err.Error())
		return
	}

	// Write the welcome message to the new client
	writer := bufio.NewWriter(connection)
	_, err = writer.WriteString(welcomeMsg)
	if err != nil {
		log.Printf("Error sending welcome message to %s: %s", connection.RemoteAddr().String(), err.Error())
		return
	}
	writer.Flush()

	// Receive client's name
	reader := bufio.NewReader(connection)
	clientName, err := reader.ReadString('\n')
	if err != nil {
		log.Printf("Error reading client name from %s: %s", connection.RemoteAddr().String(), err.Error())
		return
	}

	clientName = strings.TrimSpace(clientName)
	if clientName == "" {
		writer.WriteString("Name cannot be empty.\n")
		writer.Flush()
		return
	}
	/*
		// Create a Client struct and add it to the clients map
		currentClient := Client{Name: clientName, Socket: connection, Data: make(chan string)}
		clients[connection] = currentClient

		// Announce the new client to others
		msg := "[" + time.Now().Format("2006-01-02 15:04:05") + "]: " + currentClient.Name + " has joined our chat...\n"
		messages <- msg

		// Start a goroutine to keep reading client's messages
		go func() {
			defer connection.Close()
			for {
				buffer := make([]byte, 1024)
				bytesRead, err := reader.Read(buffer)
				if err != nil {
					// Client has disconnected
					delete(clients, connection)
					msg := "[" + time.Now().Format("2006-01-02 15:04:05") + "]: " + currentClient.Name + " has left the chat...\n"
					messages <- msg
					return
				}
				if bytesRead > 0 {
					// Sending received message to the channel
					msg := "[" + time.Now().Format("2006-01-02 15:04:05") + "][" + currentClient.Name + "]: " + string(buffer[:bytesRead])
					messages <- msg
				}
			}
		}()*/
}
