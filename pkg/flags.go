package penguin

import ("net"
"strings"
"time"
)
func Flags(clientMessage string, connection net.Conn, currentClient Client){
	if len(clientMessage) > 6 && clientMessage[0:7] == "--name=" {
	previousName := currentClient.Name
	newName := clientMessage[7:]
	newName = strings.ReplaceAll(newName, "\n", "")
	newName = strings.TrimSpace(newName)
	currentClient = Client{Name: newName, Socket: connection}
	Clients[connection] = currentClient
	for _, client := range Clients { // broadcast message to all users that current client disconnected
		if currentClient.Socket != client.Socket { // send to all clients that someone changed his name except that person
			client.Socket.Write([]byte("\n" + previousName + " has changed his name to " + currentClient.Name + "\n"))
			client.Socket.Write([]byte("[" + time.Now().Format("2006-01-02 15:04:05") + "][" + client.Name + "]: "))
		} else {
			client.Socket.Write([]byte("[" + time.Now().Format("2006-01-02 15:04:05") + "][" + client.Name + "]: "))
		}
	}
	AllMessages = append(AllMessages, previousName + " has changed his name to " + currentClient.Name + "\n")
} else if len(clientMessage) > 6 && clientMessage[0:7] == "--users" { // flag for number of users
	var arrayForUser []byte
	arrayForUser = append(arrayForUser, byte(UserCounter+47))
	connection.Write([]byte("number of users in all group chats is " + string(arrayForUser) + "\n"))
	connection.Write([]byte("[" + time.Now().Format("2006-01-02 15:04:05") + "][" + currentClient.Name + "]: "))
} else { // if wrong flag used or only '--' present show all available flags
	connection.Write([]byte("available flags are:\n" + "'--users': shows number of users in group\n"+"'--name=': to change your name\n"))
	connection.Write([]byte("[" + time.Now().Format("2006-01-02 15:04:05") + "][" + currentClient.Name + "]: "))
}
}