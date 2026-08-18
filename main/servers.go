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

		go func() {
			for {

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

				_, err = conn.Write([]byte("[Server A] Message received"))
				if err != nil {
					fmt.Println("Error sending message to client")
					conn.Close()
					return
				}

			}

		}()

	}

}

func runServerB() {
	listener, err := net.Listen("tcp", ":52")
	if err != nil {
		fmt.Println("Failed to create serverSocket")
		return
	}

	defer listener.Close()

	for {
		conn, err := listener.Accept()

		buffer := make([]byte, 1024)

		go func() {
			for {

				_, err = conn.Read(buffer)
				if err != nil {
					fmt.Println("Error reading from socket")
					return
				}
				fmt.Println("[Server B] Msg received from client, dont tell the silly load balancer.......")
				fmt.Printf("Msg: %s\n", string(buffer))
				if string(buffer) == "turn off" {
					conn.Close()
					break
				}

				_, err = conn.Write([]byte("[Server B] Message received"))
				if err != nil {
					fmt.Println("Error sending message to client")
					conn.Close()
					return
				}

			}

		}()

	}

}

func runServerC() {
	listener, err := net.Listen("tcp", ":53")
	if err != nil {
		fmt.Println("Failed to create serverSocket")
		return
	}

	defer listener.Close()

	for {
		conn, err := listener.Accept()

		buffer := make([]byte, 1024)

		go func() {
			for {

				_, err = conn.Read(buffer)
				if err != nil {
					fmt.Println("Error reading from socket")
					return
				}
				fmt.Println("[Server C] Msg received from client, dont tell the silly load balancer.......")
				fmt.Printf("Msg: %s\n", string(buffer))
				if string(buffer) == "turn off" {
					conn.Close()
					break
				}

				_, err = conn.Write([]byte("[Server C] Message received"))
				if err != nil {
					fmt.Println("Error sending message to client")
					conn.Close()
					return
				}

			}

		}()

	}

}
