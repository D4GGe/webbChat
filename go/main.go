package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)
import _ "github.com/go-sql-driver/mysql"

type Customer struct {
	Name    string `json:"name"`
	Company string `json:"company"`
}

type Message struct {
	User string
	Text string
	Time string
	Id   int
}

func getChat(w http.ResponseWriter, r *http.Request) {
	lastId := r.FormValue("id")
	fmt.Printf(lastId)
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	url := strings.Split(r.URL.Path[1:], "/")
	args := url[len(url)-1]
	//fmt.Printf("%s", args)
	db, err := sql.Open("mysql", "root:@/webchat")
	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
	defer db.Close()

	stmtOut, err := db.Prepare("SELECT messenges.id, user.name, messenges.text, messenges.time FROM chatrooms INNER JOIN messenges ON messenges.room = chatrooms.id INNER JOIN user ON user.id = messenges.user WHERE chatrooms.name=? and messenges.id >? ORDER BY messenges.id")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer stmtOut.Close()

	rows, _ := stmtOut.Query(args, lastId)
	msgs := []Message{}
	for rows.Next() {
		msg := Message{}
		rows.Scan(&msg.Id, &msg.User, &msg.Text, &msg.Time)
		msgs = append(msgs, msg)

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
