package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
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
type config struct {
	TcpPort  string `json:"tcp_port"`
	HttpPort string `json:"http_port"`
	Servers  []item `json:"servers"`
}

type stats struct {
	Size    int    `json:"size"`
	Servers []item `json:"servers"`
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

func startHttpServer(handler http.HandlerFunc, port string) {
	http.HandleFunc("/api/stats", handler)
	http.ListenAndServe(port, nil)

}

func closure(c *ConnectionPool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c.mu.Lock()
		defer c.mu.Unlock()
		s := stats{
			Size:    len(c.servers),
			Servers: c.servers,
		}

		b, err := json.Marshal(&s)
		if err != nil {
			fmt.Printf("Error creating json stats object: %s\n", err.Error())
			return
		}

		w.Write(b)

	}
}

func main() {
	b, err := os.ReadFile("./config.json")
	if err != nil {
		fmt.Printf("Error reading config.json: %s\n", err.Error())
		return
	}

	var co config
	err = json.Unmarshal(b, &co)
	if err != nil {
		fmt.Printf("Error unmarshalling: %s\n", err.Error())
		return
	}

	c := ConnectionPool{
		servers:  co.Servers,
		smallest: 0,
	}

	count := 0
	listener, err := net.Listen("tcp", co.TcpPort)
	if err != nil {
		fmt.Println("Error while trying ot create server socket")
		return
	}
	defer listener.Close()

	go sendHeartBeats(&c)

	handler := closure(&c)
	go startHttpServer(handler, co.HttpPort)

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
