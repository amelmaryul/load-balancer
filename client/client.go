package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", ":50")
	if err != nil {
		fmt.Println("Yeah mate this dont be working so gg")
		return
	}
	defer conn.Close()
	scanner := bufio.NewScanner(os.Stdin)

	buffer := make([]byte, 1024)

	for {
		fmt.Println("Type in message to send!")
		scanner.Scan()
		line := scanner.Text()

		_, err := conn.Write([]byte(line))
		if err != nil {
			fmt.Println("Failed sendig message to server")
			return
		}
		if line == "turn off" {
			break
		}

		n, err := conn.Read(buffer)
		fmt.Printf("Server says: %s\n", string(buffer[:n]))

	}

}
