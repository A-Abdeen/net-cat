package penguin

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"
)

func Flags(clientMessage string, connection net.Conn, currentClient Client) {
	if len(clientMessage) > 5 && clientMessage[0:6] == "--name" {
		previousName := currentClient.Name
		newName := clientMessage[6:]
		newName = strings.ReplaceAll(newName, "\n", "")
		newName = strings.TrimSpace(newName)
		if newName != "" {
			for i := 0; i < 1; i++ {
				for _, client := range Clients {
					if newName == client.Name {
						newName = newName + "0"
						i = -1
					}
				}
			}
			currentClient.Name = newName
			ClientsMutex.Lock()                 // lock global variable for other go routines
			Clients[connection] = currentClient // update the Clients map
			ClientsMutex.Unlock()               // unlock global variable for other go routines
			for _, client := range Clients {
				if currentClient.Socket != client.Socket && client.Group == currentClient.Group { // send to all clients in specific chat that the current user changed his name
					client.Socket.Write([]byte("\n" + previousName + " has changed his name to " + currentClient.Name + "\n"))
					client.Socket.Write([]byte("[Group " + client.Group + "][" + time.Now().Format("15:04:05") + "][" + client.Name + "]:"))
				} else if currentClient.Socket == client.Socket {
					client.Socket.Write([]byte("Your name has successfully changed to " + currentClient.Name + "\n"))
					client.Socket.Write([]byte("[Group " + client.Group + "][" + time.Now().Format("15:04:05") + "][" + client.Name + "]:"))
				} else {
					continue
				}
			}
			AllMessages[currentClient.Group] = append(AllMessages[currentClient.Group], previousName+" has changed his name to "+currentClient.Name)
		} else {
			connection.Write([]byte("name cannot be empty"))
			connection.Write([]byte("[Group " + currentClient.Group + "][" + time.Now().Format("15:04:05") + "][" + currentClient.Name + "]:"))
		}
	} else if len(clientMessage) > 6 && clientMessage[0:7] == "--users" { // flag to show the number of users
		var arrayForUser []byte
		arrayForUser = []byte(fmt.Sprint(UserCounter - 1))
		var max = []byte(fmt.Sprint(10))
		if len(arrayForUser) == len(max) {
			connection.Write([]byte("Total number of users in all groups is " + string(arrayForUser) + " (Server is full).\n"))
		} else {
			connection.Write([]byte("Total number of users in all groups is " + string(arrayForUser) + "\n"))
		}
		connection.Write([]byte("[Group " + currentClient.Group + "][" + time.Now().Format("15:04:05") + "][" + currentClient.Name + "]:"))
	} else if len(clientMessage) > 7 && clientMessage[0:8] == "--switch" { // flag for switching groups
		groupIn := currentClient.Group
		groupToSwitch := clientMessage[8:]
		groupToSwitch = strings.ReplaceAll(groupToSwitch, "\n", "")
		groupToSwitch = strings.TrimSpace(groupToSwitch)
		if groupToSwitch == "1" || groupToSwitch == "2" {
			if Clients[connection].Group == groupToSwitch { // if the user chooses the group he is already in
				connection.Write([]byte("You are already in group chat " + Clients[connection].Group + "\n"))
				connection.Write([]byte("[Group " + currentClient.Group + "][" + time.Now().Format("15:04:05") + "][" + currentClient.Name + "]:"))
			} else {
				currentClient.Group = groupToSwitch
				ClientsMutex.Lock()                 // lock global variable for other go routines
				Clients[connection] = currentClient // update the Clients map
				ClientsMutex.Unlock()               // unlock global variable for other go routines

				// will show chat history for users that join later
				if len(AllMessages[groupToSwitch]) != 0 {
					connection.Write([]byte("\n----------------------history----------------------\n"))
				}
				for _, pastMessage := range AllMessages[groupToSwitch] {
					connection.Write([]byte(pastMessage))
				}
				if len(AllMessages[groupToSwitch]) != 0 {
					connection.Write([]byte("----------------------history----------------------\n"))
				}

				for _, client := range Clients {
					if currentClient.Socket != client.Socket && client.Group == currentClient.Group { // send to all clients that the current user has switched groups
						client.Socket.Write([]byte("\n" + currentClient.Name + " has joined this chat " + "\n"))
						client.Socket.Write([]byte("[Group " + client.Group + "][" + time.Now().Format("15:04:05") + "][" + client.Name + "]:"))
					} else if currentClient.Socket == client.Socket {
						client.Socket.Write([]byte("You joined Group " + groupToSwitch + "\n"))
						client.Socket.Write([]byte("[Group " + currentClient.Group + "][" + time.Now().Format("15:04:05") + "][" + currentClient.Name + "]:"))
					} else if currentClient.Socket != client.Socket && client.Group == groupIn {
						client.Socket.Write([]byte("\n" + currentClient.Name + " has left this chat " + "\n"))
						client.Socket.Write([]byte("[Group " + client.Group + "][" + time.Now().Format("15:04:05") + "][" + client.Name + "]:"))
					} else {
						continue
					}
				}
				AllMessages[currentClient.Group] = append(AllMessages[currentClient.Group], "\n"+currentClient.Name+" has joined chat "+"\n")
				AllMessages[groupIn] = append(AllMessages[groupIn], "\n"+currentClient.Name+" has left this chat "+"\n")

			}
		} else if groupToSwitch == "private" {
			if Clients[connection].Group == groupToSwitch { // if the user chooses the group he is already in
				connection.Write([]byte("You are already in group chat " + Clients[connection].Group + "\n"))
				connection.Write([]byte("[Group " + currentClient.Group + "][" + time.Now().Format("15:04:05") + "][" + currentClient.Name + "]:"))
			} else {
				reader := bufio.NewReader(connection)
				connection.Write([]byte("This is a private chamber, type the secret password: "))

				for {
					password, err := reader.ReadString('\n')

					password = strings.ReplaceAll(password, "\n", "")
					password = strings.ReplaceAll(password, " ", "")
					if err != nil {
						// Lock the ClientsMutex before accessing the Clients map.
						ClientsMutex.Lock()

						// Remove the client from the Clients map
						delete(Clients, currentClient.Socket)

						// unlock the variable so other go routines can access the variable
						ClientsMutex.Unlock()
						connection.Close()
						UserCounter--
						return
					}
					if password == Secret {
						break
					}
					if password == "--quit" {
						connection.Write([]byte("You chose to quit(loser), returning to current group\n"))
						connection.Write([]byte("[Group " + currentClient.Group + "][" + time.Now().Format("15:04:05") + "][" + currentClient.Name + "]:"))
						return
					}
					connection.Write([]byte("Password is wrong. try again noob or do --quit: "))
				}
				currentClient.Group = groupToSwitch
				ClientsMutex.Lock()                 // lock global variable for other go routines
				Clients[connection] = currentClient // update the Clients map
				ClientsMutex.Unlock()               // unlock global variable for other go routines

				// will show chat history for users that join later
				if len(AllMessages[groupToSwitch]) != 0 {
					connection.Write([]byte("\n----------------------history----------------------\n"))
				}
				for _, pastMessage := range AllMessages[groupToSwitch] {
					connection.Write([]byte(pastMessage))
				}
				if len(AllMessages[groupToSwitch]) != 0 {
					connection.Write([]byte("----------------------history----------------------\n"))
				}

				for _, client := range Clients {
					if currentClient.Socket != client.Socket && client.Group == currentClient.Group { // send to all clients that the current user has switched groups
						client.Socket.Write([]byte("\n" + currentClient.Name + " has joined this chat " + "\n"))
						client.Socket.Write([]byte("[Group " + client.Group + "][" + time.Now().Format("15:04:05") + "][" + client.Name + "]:"))
					} else if currentClient.Socket == client.Socket {
						client.Socket.Write([]byte("You joined Group " + groupToSwitch + "\n"))
						client.Socket.Write([]byte("[Group " + currentClient.Group + "][" + time.Now().Format("15:04:05") + "][" + currentClient.Name + "]:"))
					} else {
						client.Socket.Write([]byte("\n" + currentClient.Name + " has left this chat " + "\n"))
						client.Socket.Write([]byte("[Group " + client.Group + "][" + time.Now().Format("15:04:05") + "][" + client.Name + "]:"))
					}
				}
				AllMessages[currentClient.Group] = append(AllMessages[currentClient.Group], "\n"+currentClient.Name+" has joined chat "+"\n")
				AllMessages[groupIn] = append(AllMessages[groupIn], "\n"+currentClient.Name+" has left this chat "+"\n")
			}
		} else { // if the user chooses group chat that is not available
			connection.Write([]byte("chat chosen is not available\n(avaiable group chats: 1, 2 and private)\n"))
			connection.Write([]byte("[Group " + currentClient.Group + "][" + time.Now().Format("15:04:05") + "][" + currentClient.Name + "]:"))
		}
	} else { // if wrong flag used or only '--' present show all available flags
		connection.Write([]byte("available flags are:\n" + "'--users': shows number of users in group\n" + "'--name': to change your name\n"))
		connection.Write([]byte("'--switch': to switch to another group chat \navailable group chats are: 1, 2 or private\n"))
		connection.Write([]byte("[Group " + currentClient.Group + "][" + time.Now().Format("15:04:05") + "][" + currentClient.Name + "]:"))
	}
}
