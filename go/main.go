package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	_ "strings"
)
import _ "github.com/go-sql-driver/mysql"

type Customer struct {
	Name    string `json:"name"`
	Company string `json:"company"`
}

type Message struct {
	User int
	Text string
	Time string
}

func getChat(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	//url := strings.Split(r.URL.Path[1:], "/")
	//args := url[len(url)-1]
	db, err := sql.Open("mysql", "root:@/webchat")
	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
	defer db.Close()

	stmtOut, err := db.Prepare("SELECT user, text, time FROM messenges")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer stmtOut.Close()

	rows, _ := stmtOut.Query()
	msgs := []Message{}
	for rows.Next() {
		msg := Message{}
		rows.Scan(&msg.User, &msg.Text, &msg.Time)
		msgs = append(msgs, msg)
		fmt.Println(msg)
	}

	t, err := json.MarshalIndent(msgs, "", "  ")
	if err == nil {
		fmt.Fprint(w, string(t))

	}
}

func main() {

	http.HandleFunc("/chat/", getChat)
	http.ListenAndServe(":8080", nil)
}
