package penguin

import "os"

// Read welcome message from file
func readWelcomeMsg() (string, error) {
	data, err := os.ReadFile("../message/welcome.txt")
	if err != nil {
		return "", err
	}
	return string(data), nil
}
