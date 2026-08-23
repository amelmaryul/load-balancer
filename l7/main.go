package main

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

var client = &http.Client{Timeout: 10 * time.Second}

// i should have these configuered
func getData(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	switch path {
	case "/players":
		resp, err := client.Get("http://localhost:61/api/players")
		defer resp.Body.Close()
		if err != nil {
			fmt.Printf("Error at get request: %s\n", err.Error())
			return
		}

		body, err := io.ReadAll(resp.Body)

		if err != nil {
			fmt.Printf("Error reading body: %s\n", err.Error())
			return
		}

		w.Write([]byte(body))
		fmt.Println(string(body))

	case "/teams":
		resp, err := client.Get("http://localhost:62/api/teams")
		defer resp.Body.Close()
		if err != nil {
			fmt.Printf("Error at get request: %s\n", err.Error())
			return
		}
		body, err := io.ReadAll(resp.Body)

		if err != nil {
			fmt.Printf("Error reading body: %s\n", err.Error())
			return
		}

		w.Write([]byte(body))
		fmt.Println(string(body))

	case "/leagues":
		resp, err := client.Get("http://localhost:63/api/leagues")
		defer resp.Body.Close()
		if err != nil {
			fmt.Printf("Error at get request: %s\n", err.Error())
			return
		}
		body, err := io.ReadAll(resp.Body)

		if err != nil {
			fmt.Printf("Error reading body: %s\n", err.Error())
			return
		}

		w.Write([]byte(body))
		fmt.Println(string(body))

	default:
		// return an error message
		w.WriteHeader(404)
	}
}

func main() {
	http.HandleFunc("/", getData)
	http.ListenAndServe(":8080", nil)
}
