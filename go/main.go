package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
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

func setHeader(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Origin", "*")
}
func getChat(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	setHeader(w)
	lastId := ps.ByName("id")
	name := ps.ByName("name")
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

	rows, _ := stmtOut.Query(name, lastId)
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

	router := httprouter.New()
	router.GET("/chat/:name/:id", getChat)

	log.Fatal(http.ListenAndServe(":8080", router))
}
