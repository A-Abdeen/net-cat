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
		connection.Close()
		UserCounter--
		return
	}

	// Write the welcome message to the new client
	writer := bufio.NewWriter(connection)
	_, err = writer.WriteString(welcomeMsg)
	if err != nil {
		log.Printf("Error sending welcome message to %s: %s", connection.RemoteAddr().String(), err.Error())
		connection.Close()
		UserCounter--
		return
	}
	writer.Flush()

	// Receive client's name
	reader := bufio.NewReader(connection)
	clientName, err := reader.ReadString('\n')
	if err != nil {
		connection.Close()
		UserCounter--
		return
	}
	// trim spaces from client's name
	clientName = strings.TrimSpace(clientName)
	if clientName == "" {
		writer.WriteString("Name cannot be empty. Reconnect\n")
		writer.Flush()
		connection.Close()
		UserCounter--
		return
	}
	for i := 0; i < 1; i++ {
		for _, client := range Clients {
			if clientName == client.Name {
				clientName = clientName + "0"
				i = -1
			}
		}
	}
	connection.Write([]byte("Choose Chat(1, 2 or private):"))
	choosen, err := reader.ReadString('\n')
	choosen = strings.ReplaceAll(choosen, "\n", "")
	if err != nil {
		connection.Close()
		UserCounter--
		return
	}
	var choosengroup string
	switch {
	case choosen == "1":
		choosengroup = "1"
	case choosen == "2":
		choosengroup = "2"
	case choosen == "private":
		choosengroup = "private"
		connection.Write([]byte("This is a private chamber, type the secret password: "))

		for {
			password, err := reader.ReadString('\n')

			password = strings.ReplaceAll(password, "\n", "")
			password = strings.ReplaceAll(password, " ", "")
			if err != nil {
				connection.Close()
				UserCounter--
				return
			}
			if password == Secret {
				break
			}
			if password == "--quit" {
				connection.Close()
				UserCounter--
				return
			}
			connection.Write([]byte("Password is wrong. try again noob or do --quit: "))
		}
	default:
		choosengroup = "1"
		connection.Write([]byte("default group chat 1 chosen\n"))
	}
	// will show chat history for users that join later
	if len(AllMessages[choosengroup]) != 0 {
		connection.Write([]byte("\n----------------------history----------------------\n"))
	}
	for _, pastMessage := range AllMessages[choosengroup] {
		connection.Write([]byte(pastMessage))
	}
	if len(AllMessages[choosengroup]) != 0 {
		connection.Write([]byte("----------------------history----------------------\n"))
	}
	// Create a Client struct and add it to the clients map

	currentClient := Client{Name: clientName, Socket: connection, Group: choosengroup}
	// lock the variable so no other go routine tries accessing it before it's added
	ClientsMutex.Lock()
	// add the client
	Clients[connection] = currentClient
	// unlock the variable so other go routines can access it after it has been added to the Clients map
	ClientsMutex.Unlock()
	// announce to all clients, the name of who joined our chat
	for _, client := range Clients {
		if currentClient.Socket != client.Socket && currentClient.Group == client.Group {
			client.Socket.Write([]byte("\n" + currentClient.Name + " has joined the chat: " + currentClient.Group + "\n"))
			client.Socket.Write([]byte("[Group " + client.Group + "][" + time.Now().Format("15:04:05") + "][" + client.Name + "]:"))
		}
	}
	AllMessages[currentClient.Group] = append(AllMessages[currentClient.Group], currentClient.Name+" has joined the chat "+"\n")

	// go routine that will keep reading each clients input
	go func() {
		defer connection.Close() // after programming is done running, it will make sure to close connection

		contreader := bufio.NewReader(connection) // variable of type reader(has capability to read)
		connection.Write([]byte("[Group " + currentClient.Group + "][" + time.Now().Format("15:04:05") + "][" + currentClient.Name + "]:"))
		for {
			clientMessage, err := contreader.ReadString('\n') // reads everything until first occurence of new line
			if err != nil {                                   // anytime an error happens, assume user has disconnected. errors could be EOF which means they did a signal interrupt
				for _, client := range Clients { // broadcast message to all users that current client disconnected
					if currentClient.Socket != client.Socket && client.Group == currentClient.Group { // send to all clients that someone left, except that person
						client.Socket.Write([]byte("\n" + currentClient.Name + " has left the chat " + "\n"))
						client.Socket.Write([]byte("[Group " + client.Group + "][" + time.Now().Format("15:04:05") + "][" + client.Name + "]:"))
					}
				}
				AllMessages[currentClient.Group] = append(AllMessages[currentClient.Group], currentClient.Name+" has left this chat \n")
				connection.Close()
				UserCounter--
				// Lock the ClientsMutex before accessing the Clients map.
				ClientsMutex.Lock()

				// Remove the client from the Clients map
				delete(Clients, currentClient.Socket)

				// unlock the variable so other go routines can access the variable
				ClientsMutex.Unlock()

				return
			}
			if len(clientMessage) > 1 && clientMessage[0:2] == "--" { // check for flag
				Flags(clientMessage, connection, currentClient)
				currentClient = Clients[connection]
			} else {
				// will check if client tries sending an empty message, if so it won't broadcast it
				clientMessage = strings.TrimSpace(clientMessage)
				if clientMessage == "" {
					connection.Write([]byte("[Group " + currentClient.Group + "][" + time.Now().Format("15:04:05") + "][" + currentClient.Name + "]:"))
					continue
				}

				// append to all messages slice, which stores all messages
				AllMessages[currentClient.Group] = append(AllMessages[currentClient.Group], "[Group "+currentClient.Group+"]["+time.Now().Format("15:04:05")+"]["+currentClient.Name+"]:"+clientMessage+"\n")

				// where messages are sent to be printed
				currentClient.Message = clientMessage
				messages <- currentClient // channel to communicate with broadcast message go routine, sends data of type client, along with his socket, message and name
			}
		}
	}()
}
