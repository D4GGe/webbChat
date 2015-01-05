package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Customer struct {
	Name    string `json:"name"`
	Company string `json:"company"`
}

type Message struct {
	Name string
	Body string
	Time int64
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	m := Customer{"Alice", "Hello"}
	t, err := json.Marshal(m)
	if err == nil {
		fmt.Fprint(w, string(t))
	}
}

func main() {
	http.HandleFunc("/test.json", handler)
	http.ListenAndServe(":8080", nil)
}
