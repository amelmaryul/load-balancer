package main

import (
	"fmt"
	"net"
	"sync"
)

type ConnectionPool struct {
	mu       sync.Mutex
	servers  []item
	smallest int
}
type item struct {
	port        string
	connections int
}

func (c *ConnectionPool) updateSmallest() {
	for i, v := range c.servers {
		if c.servers[c.smallest].connections > v.connections {
			c.smallest = i
		}
	}
}

func balance(clientConn net.Conn, c *ConnectionPool) {
	c.mu.Lock()
	port := c.servers[c.smallest].port
	c.servers[c.smallest].connections++
	idx := c.smallest
	c.updateSmallest()
	c.mu.Unlock()

	serverConn, err := net.Dial("tcp", port)
	if err != nil {
		fmt.Println("Error creating connection object")
		return
	}

	clientBuffer := make([]byte, 1024)
	serverBuffer := make([]byte, 1024)

	defer func() {
		c.mu.Lock()
		c.servers[idx].connections--
		c.updateSmallest()
		c.mu.Unlock()
	}()

	for {
		n, err := clientConn.Read(clientBuffer)
		if err != nil {
			fmt.Println("Error reading from connection")
			return
		}
		_, err = serverConn.Write(clientBuffer[:n])
		if err != nil {
			fmt.Println("Error writing to connection")
			return
		}

		n, err = serverConn.Read(serverBuffer)
		if err != nil {
			fmt.Println("Error reading from connection")
			return
		}

		_, err = clientConn.Write(serverBuffer[:n])
		if err != nil {
			fmt.Println("Error writing to connection")
			return
		}

	}
}

func main() {
	s := make([]item, 3) // 3 here is the amount of servers we have
	s[0] = item{
		port:        ":51",
		connections: 0,
	}

	s[1] = item{
		port:        ":52",
		connections: 0,
	}

	s[2] = item{
		port:        ":53",
		connections: 0,
	}

	c := ConnectionPool{
		servers:  s,
		smallest: 0,
	}

	go runServerA()
	go runServerB()
	go runServerC()

	count := 0
	listener, err := net.Listen("tcp", ":50")
	if err != nil {
		fmt.Println("Error while trying ot create server socket")
		return
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error creating connection object")
			return
		}
		count++
		if count > 3 {
			count = 1
		}
		go balance(conn, &c)
	}

}
