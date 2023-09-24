package penguin

import (
	"net"
	"strings"
	"time"
)
func Flags(clientMessage string, connection net.Conn, currentClient Client){
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
	Clients[connection] = currentClient
	for _, client := range Clients { 
		if currentClient.Socket != client.Socket { // send to all clients that the current user changed his name
			client.Socket.Write([]byte("\n" + previousName + " has changed his name to " + currentClient.Name + "\n"))
			client.Socket.Write([]byte("[Group " + client.Group+ "][" + time.Now().Format("15:04:05") + "][" + client.Name + "]:"))
		} else {
			client.Socket.Write([]byte("[Group " + currentClient.Group+ "][" + time.Now().Format("15:04:05") + "][" + currentClient.Name + "]:"))
		}
	}
	AllMessages = append(AllMessages, previousName + " has changed his name to " + currentClient.Name)
} else {
	connection.Write([]byte("name cannot be empty"))
	connection.Write([]byte("[Group " + currentClient.Group+ "][" + time.Now().Format("15:04:05") + "][" + currentClient.Name + "]:"))
}
} else if len(clientMessage) > 6 && clientMessage[0:7] == "--users" { // flag to show the number of users
	var arrayForUser []byte
	arrayForUser = append(arrayForUser, byte(UserCounter+47))
	connection.Write([]byte("number of users in all group chats is " + string(arrayForUser) + "\n"))
	connection.Write([]byte("[Group " + currentClient.Group+ "][" + time.Now().Format("15:04:05") + "][" + currentClient.Name + "]:"))
}  else if len(clientMessage) > 7 && clientMessage[0:8] == "--switch" { // flag for switching groups
	groupIn := currentClient.Group
	groupToSwitch := clientMessage[8:]
	groupToSwitch = strings.ReplaceAll(groupToSwitch, "\n", "")
	groupToSwitch = strings.TrimSpace(groupToSwitch)
	if groupToSwitch == "1" || groupToSwitch == "2" || groupToSwitch == "3"{
	
	if Clients[connection].Group == groupToSwitch { // if the user chooses the group he is already in
		connection.Write([]byte("You are already in group chat " + Clients[connection].Group + "\n"))
	connection.Write([]byte("[Group " + currentClient.Group+ "][" + time.Now().Format("15:04:05") + "][" + currentClient.Name + "]:"))
	} else {
		currentClient.Group = groupToSwitch
		Clients[connection] = currentClient
	for _, client := range Clients { 
		if currentClient.Socket != client.Socket { // send to all clients that the current user has switched groups
			client.Socket.Write([]byte("\n" + currentClient.Name + " has switched from group chat " + groupIn + " to " + groupToSwitch + "\n"))
			client.Socket.Write([]byte("[Group " + client.Group+ "][" + time.Now().Format("15:04:05") + "][" + client.Name + "]:"))
		} else {
			client.Socket.Write([]byte("You switched from group chat " + groupIn + " to " + groupToSwitch + "\n"))
			client.Socket.Write([]byte("[Group " + currentClient.Group+ "][" + time.Now().Format("15:04:05") + "][" + currentClient.Name + "]:"))
		}
	} 
	AllMessages = append(AllMessages, "\n" + currentClient.Name + " has switched from group chat " + groupIn + " to " + groupToSwitch + "\n")
	}} else { // if the user chooses group chat that is not available
		connection.Write([]byte("Group chat choosen is not available\n(avaiable group chats: 1, 2 and 3)\n"))
		connection.Write([]byte("[Group " + currentClient.Group+ "][" + time.Now().Format("15:04:05") + "][" + currentClient.Name + "]:"))
	}
	 } else { // if wrong flag used or only '--' present show all available flags
	connection.Write([]byte("available flags are:\n" + "'--users': shows number of users in group\n"+"'--name': to change your name\n"))
	connection.Write([]byte("'--switch': to switch to another group chat \navailable group chats are: 1, 2 or 3\n"))
	connection.Write([]byte("[Group " + currentClient.Group+ "][" + time.Now().Format("15:04:05") + "][" + currentClient.Name + "]:"))
}
}