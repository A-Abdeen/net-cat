package penguin

import "net"

type Client struct {
	Name   string
	Socket net.Conn
	Data   chan string
}

var (
	clients  = make(map[net.Conn]Client)
	messages = make(chan string)
)
