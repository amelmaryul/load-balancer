package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

var client = &http.Client{Timeout: 10 * time.Second}

type Route struct {
	PathPrefx string `json:"pathprefix"`
	Backend   string
}

type Config struct {
	LbPort string
	Routes []Route
}

func closure(config Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		for _, v := range config.Routes {
			if path == v.PathPrefx {
				resp, err := client.Get(v.Backend)
				if err != nil {
					fmt.Printf("Error at get request: %s\n", err.Error())
					return
				}
				defer resp.Body.Close()
				body, err := io.ReadAll(resp.Body)

				if err != nil {
					fmt.Printf("Error reading body: %s\n", err.Error())
					return
				}

				w.Write([]byte(body))
				fmt.Println(string(body))

			}
		}
	}
}
func main() {
	b, err := os.ReadFile("./config.json")
	if err != nil {
		fmt.Printf("Error reading config file: %s\n", err.Error())
		return
	}
	var config Config
	err = json.Unmarshal(b, &config)
	if err != nil {
		fmt.Printf("Error unmarshaling: %s\n", err.Error())
		return
	}
	handler := closure(config)

	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
