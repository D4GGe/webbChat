package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"log"
	"math"
	"net/http"
	"regexp"
	"strconv"
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

func getCircCord(x0, y0, r, x float64) (float64, float64) {
	y1 := y0 + math.Sqrt(r*r-x*x+2*x*x0-x0*x0)
	y2 := y0 - math.Sqrt(r*r-x*x+2*x*x0-x0*x0)
	return y1, y2

}
func genImg(w http.ResponseWriter, r *http.Request) {
	m := image.NewRGBA(image.Rect(0, 0, 500, 500))
	blue := color.RGBA{40, 40, 40, 255}
	draw.Draw(m, m.Bounds(), &image.Uniform{blue}, image.ZP, draw.Src)
	var y int
	for x := m.Rect.Min.X; x < m.Rect.Max.X; x++ {
		f := float64(x)
		y1, y2 := getCircCord(250.0, 250.0, 250.0, f)
		y12, y22 := getCircCord(250.0, 250.0, 230.0, f)
		//fmt.Println("y1: " + strconv.FormatFloat(y1, 'f', 6, 64) + " y12: " + strconv.FormatFloat(y12, 'f', 6, 64))
		//fmt.Println("y2: " + strconv.FormatFloat(y2, 'f', 6, 64) + " y22: " + strconv.FormatFloat(y22, 'f', 6, 64))

		if int(y12) > 0 && int(y1) > 0 {

			for y = int(y12); y < int(y1); y++ {

				m.Set(x, int(y), color.RGBA{uint8(math.Abs(math.Sin(float64(x)/100)*100) + 120), uint8(math.Abs(math.Sin((float64(x)+100)/100)*100) + 120), uint8(math.Abs(math.Sin((float64(x)+200)/100)*100) + 120), 255})
			}

		} else {
			for y = 250; y < int(y1); y++ {
				m.Set(x, int(y), color.RGBA{uint8(math.Abs(math.Sin(float64(x)/100)*100) + 120), uint8(math.Abs(math.Sin((float64(x)+100)/100)*100) + 120), uint8(math.Abs(math.Sin((float64(x)+200)/100)*100) + 120), 255})
			}

		}

		if int(y22) >= 0 && int(y2) >= 0 {
			for y = int(y22); y > int(y2); y-- {

				m.Set(x, int(y), color.RGBA{uint8(math.Abs(math.Sin(float64(x)/100)*100) + 120), uint8(math.Abs(math.Sin((float64(x)+100)/100)*100) + 120), uint8(math.Abs(math.Sin((float64(x)+200)/100)*100) + 120), 255})
			}
		} else {
			for y = 250; y > int(y2); y-- {

				m.Set(x, int(y), color.RGBA{uint8(math.Abs(math.Sin(float64(x)/100)*100) + 120), uint8(math.Abs(math.Sin((float64(x)+100)/100)*100) + 120), uint8(math.Abs(math.Sin((float64(x)+200)/100)*100) + 120), 255})
			}

		}

	}
	var img image.Image = m

	writeImage(w, &img)
}

// writeImage encodes an image 'img' in jpeg format and writes it into ResponseWriter.
func writeImage(w http.ResponseWriter, img *image.Image) {

	buffer := new(bytes.Buffer)
	if err := jpeg.Encode(buffer, *img, nil); err != nil {
		log.Println("unable to encode image.")
	}

	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))
	if _, err := w.Write(buffer.Bytes()); err != nil {
		log.Println("unable to write image.")
	}
}

func main() {

	router := mux.NewRouter()
	http.Handle("/", router)
	router.HandleFunc("/img", genImg).Methods("GET")
	router.HandleFunc("/login", login).Methods("POST")
	router.HandleFunc("/chat/{name}/{id}", getChat).Methods("GET")
	router.HandleFunc("/chat/{name}", postMsg).Methods("POST")
	router.HandleFunc("/chat/{name}", setHeader).Methods("OPTIONS")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
