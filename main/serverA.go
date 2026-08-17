package main

import (
	"fmt"
	"net"
)

func runServerA() {
	listener, err := net.Listen("tcp", ":51")
	if err != nil {
		fmt.Println("Failed to create serverSocket")
		return
	}

	defer listener.Close()

	for {
		conn, err := listener.Accept()

		buffer := make([]byte, 1024)

		_, err = conn.Read(buffer)
		if err != nil {
			fmt.Println("Error reading from socket")
			return
		}
		fmt.Println("[Server A] Msg received from client, dont tell the silly load balancer.......")
		fmt.Printf("Msg: %s\n", string(buffer))
		if string(buffer) == "turn off" {
			conn.Close()
			break
		}

	}

}
