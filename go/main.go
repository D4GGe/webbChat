package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"log"
	"net/http"
	"regexp"
)
import _ "github.com/go-sql-driver/mysql"

var store = sessions.NewCookieStore([]byte("something-very-secret"))

type User struct {
	id      int
	name    string
	picture string
}

type Message struct {
	User string
	Text string
	Time string
	Id   int
}

func setHeader(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Origin", "*")
}
func getChat(w http.ResponseWriter, r *http.Request) {
	reg, _ := regexp.Compile("p([a-z]+)ch")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	vars := mux.Vars(r)
	lastId := vars["id"]
	name := vars["name"]
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
		msg.Text = reg.ReplaceAllString(msg.Text, "\"<img src='http://a.deviantart.net/avatars/n/u/number1peachfan.gif?3'>\"")
		msgs = append(msgs, msg)

	}

	t, err := json.MarshalIndent(msgs, "", "  ")
	if err == nil {

		fmt.Fprint(w, string(t))

	}
}
func postMsg(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	vars := mux.Vars(r)
	name := vars["name"]
	db, _ := sql.Open("mysql", "root:@/webchat")
	//msg, room , userid
	stmtOut, _ := db.Prepare("CALL send_msg('" + r.FormValue("msg") + "','" + name + "',1)")
	stmtOut.Query()
	//fmt.Println("room:" + name + "   msg:" + r.FormValue("msg"))

}

func login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	name := r.FormValue("name")
	password := r.FormValue("password")
	session, _ := store.Get(r, "user")
	user := User{}
	if len(name) > 3 && len(password) > 3 && session.Values["id"] == nil {
		db, _ := sql.Open("mysql", "root:@/webchat")
		stmtOut, _ := db.Prepare("SELECT `user`.id,`user`.`name` FROM `user` WHERE `user`.`name` =? and `user`.`password`=?")
		rows, _ := stmtOut.Query(name, password)
		if rows.Next() {
			rows.Scan(&user.id, &user.name)
			session.Values["name"] = user.id
			session.Values["id"] = user.name
			session.Save(r, w)
		}

	}
	t, err := json.MarshalIndent(user, "", "  ")
	if err == nil {
		fmt.Fprint(w, string(t))
	}

}

func main() {

	router := mux.NewRouter()
	http.Handle("/", router)
	router.HandleFunc("/login", login).Methods("POST")
	router.HandleFunc("/chat/{name}/{id}", getChat).Methods("GET")
	router.HandleFunc("/chat/{name}", postMsg).Methods("POST")
	router.HandleFunc("/chat/{name}", setHeader).Methods("OPTIONS")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
