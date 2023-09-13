package penguin

import "time"

func BroadcastMessages() {
	for {
		msg := <-messages
		for _, client := range Clients {
			if msg.Socket != client.Socket {
				client.Socket.Write([]byte("\n" + "[" + time.Now().Format("2006-01-02 15:04:05") + "][" + msg.Name + "]: " + msg.Message))
				client.Socket.Write([]byte("[" + time.Now().Format("2006-01-02 15:04:05") + "][" + client.Name + "]: "))
			} else {
				client.Socket.Write([]byte("[" + time.Now().Format("2006-01-02 15:04:05") + "][" + client.Name + "]: "))
			}
		}

	}
}
