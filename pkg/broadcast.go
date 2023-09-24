package penguin

import "time"

func BroadcastMessages() {
	for {
		msg := <-messages
		for _, client := range Clients {
		if msg.Group == client.Group {
			if msg.Socket != client.Socket {
				client.Socket.Write([]byte("\n[Group " + client.Group+ "][" + time.Now().Format("15:04:05") + "][" + msg.Name + "]:" + msg.Message + "\n"))
				client.Socket.Write([]byte("[Group " + client.Group+ "][" + time.Now().Format("15:04:05") + "][" + client.Name + "]:"))
			} else {
				client.Socket.Write([]byte("[Group " + msg.Group+ "][" + time.Now().Format("15:04:05") + "][" + msg.Name + "]:"))
			}
		}
	}
	}
}
