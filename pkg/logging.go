package penguin

import (
	"log"
	"os"
)

// SaveAllMessagesToFile saves all messages to a file.
func SaveAllMessagesToFile(filename string) {
	file, err := os.Create(filename)
	if err != nil {
		log.Printf("Error creating chat log file %s: %s", filename, err.Error())
		return
	}
	defer file.Close()

	for index, chat := range AllMessages {
		_, err := file.WriteString(index + "\n")
		if err != nil {
			log.Printf("Error writing message to chat log file: %s", err.Error())
		}
		for _, msg := range chat {
			_, err := file.WriteString(msg + "\n")
			if err != nil {
				log.Printf("Error writing message to chat log file: %s", err.Error())
			}
		}
	}
}
