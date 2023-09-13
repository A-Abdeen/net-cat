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

	for _, msg := range AllMessages {
		_, err := file.WriteString(msg + "\n")
		if err != nil {
			log.Printf("Error writing message to chat log file: %s", err.Error())
		}
	}
}
