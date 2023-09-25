package penguin

import (
	"net"
	"sync"
)

type Client struct {
	Name    string
	Socket  net.Conn
	Message string
	Group   string
}

var (
	ClientsMutex sync.Mutex
	Clients      = make(map[net.Conn]Client)
	messages     = make(chan Client)
	AllMessages  = make(map[string][]string)
	// User limit and counter
	MaxUsers    = 10
	UserCounter = 1
	Secret      = "8008"
)
