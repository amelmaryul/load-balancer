package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"sync"
	"time"
)

type ConnectionPool struct {
	mu       sync.Mutex
	servers  []item
	smallest int
}
type item struct {
	Port        string
	Connections int
	IsAlive     bool
}

func (c *ConnectionPool) updateSmallest() {
	for i, v := range c.servers {
		if !c.servers[c.smallest].IsAlive && v.IsAlive || (v.IsAlive && v.Connections < c.servers[c.smallest].Connections) {
			c.smallest = i
		}
	}
}

func sendHeartBeats(c *ConnectionPool) {
	for {
		c.mu.Lock()
		for i, v := range c.servers {
			isAlive := checkIfAlive(v.Port)
			c.servers[i].IsAlive = isAlive

		}
		c.updateSmallest()
		c.mu.Unlock()
		time.Sleep(500 * time.Millisecond)
	}
}

func checkIfAlive(port string) bool {
	conn, err := net.Dial("tcp", port)
	if err != nil {
		return false
	}

	_, err = conn.Write([]byte("Are you alive"))
	if err != nil {
		return false
	}

	return true

}

func balance(clientConn net.Conn, c *ConnectionPool) {
	c.mu.Lock()
	port := c.servers[c.smallest].Port
	if !c.servers[c.smallest].IsAlive {
		clientConn.Write([]byte("All servers are down"))
		c.mu.Unlock()
		return
	}
	c.servers[c.smallest].Connections++
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
		c.servers[idx].Connections--
		c.updateSmallest()
		c.mu.Unlock()
	}()

	for {
		n, err := clientConn.Read(clientBuffer)
		if err != nil {
			fmt.Printf("Error reading from client connection: %s\n", err.Error())
			return
		}
		_, err = serverConn.Write(clientBuffer[:n])
		if err != nil {
			fmt.Println("Error writing to connection")
			return
		}

		n, err = serverConn.Read(serverBuffer)
		if err != nil {
			fmt.Printf("Error reading from server connection: %s\n", err.Error())
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
	b, err := os.ReadFile("./config.json")
	if err != nil {
		fmt.Printf("Error reading config.json: %s\n", err.Error())
		return
	}

	var s []item
	err = json.Unmarshal(b, &s)
	if err != nil {
		fmt.Printf("Error unmarshalling: %s\n", err.Error())
		return
	}

	c := ConnectionPool{
		servers:  s,
		smallest: 0,
	}

	count := 0
	listener, err := net.Listen("tcp", ":60")
	if err != nil {
		fmt.Println("Error while trying ot create server socket")
		return
	}
	defer listener.Close()

	go sendHeartBeats(&c)

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
