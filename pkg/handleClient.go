package penguin

import (
	"bufio"
	// "fmt"
	"log"
	"net"
	"strings"
	"time"
)

func HandleClient(connection net.Conn) {
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
		writer.WriteString("Name cannot be empty. Reconnect\n")
		writer.Flush()
		connection.Close()
		UserCounter--
		return
	}

	// Create a Client struct and add it to the clients map
	currentClient := Client{Name: clientName, Socket: connection}
	Clients[connection] = currentClient

	for _, client := range Clients {
		if currentClient.Socket != client.Socket {
			client.Socket.Write([]byte("\n" + currentClient.Name + " has joined our chat...\n"))
			client.Socket.Write([]byte("[" + time.Now().Format("2006-01-02 15:04:05") + "][" + client.Name + "]: "))
		}
	}
	AllMessages = append(AllMessages, currentClient.Name+" has joined our chat...\n")

	go func() {
		defer connection.Close()

		contreader := bufio.NewReader(connection)
		connection.Write([]byte("[" + time.Now().Format("2006-01-02 15:04:05") + "][" + currentClient.Name + "]: "))
		for {

			clientMesagge, err := contreader.ReadString('\n')
			if err != nil {
				if err.Error() == "EOF" {
					for _, client := range Clients {
						if currentClient.Socket != client.Socket {
							client.Socket.Write([]byte("\n" + currentClient.Name + " has left our chat...\n"))
							client.Socket.Write([]byte("[" + time.Now().Format("2006-01-02 15:04:05") + "][" + client.Name + "]: "))
						}
					}
					AllMessages = append(AllMessages, currentClient.Name+" has left our chat...\n")
					connection.Close()
					UserCounter--
				}
				return
			}
			// fmt.Print("[" + time.Now().Format("2006-01-02 15:04:05") + "][" + currentClient.Name + "]: " + clientMesagge) // XXX

			// append to all messages slice
			AllMessages = append(AllMessages, "["+time.Now().Format("2006-01-02 15:04:05")+"]["+currentClient.Name+"]: "+clientMesagge)

			currentClient.Message = clientMesagge
			messages <- currentClient

		}
	}()
}
