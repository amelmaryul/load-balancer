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
	Backend   string `json:"backend"`
	Stats     string `json:"stats"`
}

type Config struct {
	LbPort string
	Routes []Route
}

type item struct {
	Port        string
	Connections int
	IsAlive     bool
}

type l4stats struct {
	Size    int    `json:"size"`
	Servers []item `json:"servers"`
}

type stats struct {
	Size     int       `json:"size"`
	Clusters []l4stats `json:"clusters"`
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

func closure2(config Config) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		l4 := make([]l4stats, 0)
		s := stats{}

		for _, v := range config.Routes {
			resp, err := client.Get(v.Stats)
			if err != nil {
				fmt.Printf("Error gettign stats from %s  %s\n", v.Stats, err.Error())
			}
			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				fmt.Printf("Error reading body from %s  %s\n", v.Stats, err.Error())
			}
			temp := l4stats{}
			json.Unmarshal(body, &temp)
			l4 = append(l4, temp)
		}

		s.Size = len(l4)
		s.Clusters = l4
		b, _ := json.Marshal(s)
		w.Write([]byte(b))

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

	getStats := closure2(config)

	http.HandleFunc("/", handler)
	http.HandleFunc("/api/stats", getStats)
	http.ListenAndServe(":8080", nil)
}
