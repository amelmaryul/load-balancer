package main

import (
	"fmt"
	"net"
)

var servers = map[int]string{
	1: ":51",
	2: ":52",
	3: ":53",
}

func balance(clientConn net.Conn, port string) {
	serverConn, err := net.Dial("tcp", port)
	if err != nil {
		fmt.Println("Error creating connection object")
		return
	}

	clientBuffer := make([]byte, 1024)
	serverBuffer := make([]byte, 1024)

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
		go balance(conn, servers[count])
	}

}
