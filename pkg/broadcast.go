package penguin

func BroadcastMessages() {
	for {
		msg := <-messages
		for _, client := range clients {
			client.Data <- msg
		}
	}
}
