package penguin

import (
	"net"
)

type Client struct {
	Name    string
	Socket  net.Conn
	Message string
}

var (
	Clients     = make(map[net.Conn]Client)
	messages    = make(chan Client)
	AllMessages []string
	// User limit and counter
	MaxUsers    = 10
	UserCounter = 1
)
