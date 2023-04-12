package main

import (
	"encoding/json"
	"fmt"
	"github.com/Jero075/GoMessenger-V2/data"
	"github.com/Jero075/GoMessenger-V2/encryption"
	"net/http"
	"strings"
	"time"
)

type ViewAccount struct {
	Username string `json:"username"`
	ID       string `json:"id"`
	Premium  bool   `json:"premium"`
}

func CheckSession(r *http.Request) (bool, int) {
	err := r.ParseForm()
	if err != nil {
		return false, 1
	}
	id := r.Form.Get("id")
	session := r.Form.Get("session")
	user := data.GetUser(id)
	for _, s := range user.Sessions {
		if s.ID == session {
			return true, 0
		}
	}
	return false, 2
}
func CleanSessions() {
	for {
		db := data.OpenDB()
		for _, user := range db.Users {
			for i, session := range user.Sessions {
				if session.Expires.Before(time.Now()) {
					user.Sessions[i] = user.Sessions[len(user.Sessions)-1]
					user.Sessions = user.Sessions[:len(user.Sessions)-1]
				}
			}
		}
		time.Sleep(time.Minute * 10)
	}
}

func Signup(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(400)
		return
	}
	username := r.Form.Get("username")
	password := r.Form.Get("password")

	id := data.AddUser(username, password)
	_, wErr := w.Write([]byte(id))
	if wErr != nil {
		fmt.Println(time.Now().UTC().String() + " | Error responding to Signup request: " + err.Error())
	}
}
func Login(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(400)
		return
	}
	id := r.Form.Get("id")
	password := r.Form.Get("password")

	user := data.GetUser(id)

	if encryption.GenerateHash512(password, user.Username) == user.PWHash {
		session := data.Idgen(32)
		user.Sessions = append(user.Sessions, data.Session{ID: session, Expires: time.Now().Add(time.Hour)})
		data.ChangeUser(user)
		_, wErr := w.Write([]byte(session))
		if wErr != nil {
			fmt.Println(time.Now().UTC().String() + " | Error responding to Login request: " + wErr.Error())
		}
	} else {
		_, wErr := w.Write([]byte(""))
		if wErr != nil {
			fmt.Println(time.Now().UTC().String() + " | Error responding to Login request: " + wErr.Error())
		}
	}
}
func ViewInfo(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(400)
		return
	}
	id := r.Form.Get("id")

	user := data.GetUser(id)
	resp := ViewAccount{Username: user.Username, ID: user.ID, Premium: user.Premium}
	jResp, jErr := json.Marshal(resp)
	if jErr != nil {
		fmt.Println(time.Now().UTC().String() + " | Error marshalling ViewInfo response: " + jErr.Error())
	}
	_, wErr := w.Write(jResp)
	if wErr != nil {
		fmt.Println(time.Now().UTC().String() + " | Error responding to ViewInfo request: " + wErr.Error())
	}
}
func ChangeUsername(w http.ResponseWriter, r *http.Request) {
	valid, status := CheckSession(r)
	if valid {
		err := r.ParseForm()
		if err != nil {
			w.WriteHeader(400)
			return
		}
		id := r.Form.Get("id")
		username := r.Form.Get("username")

		user := data.GetUser(id)
		user.Username = username
		data.ChangeUser(user)
		w.WriteHeader(201)
	} else {
		if status == 1 {
			w.WriteHeader(400)
		} else if status == 2 {
			w.WriteHeader(401)
		}
	}
}
func ChangePassword(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(400)
		return
	}
	id := r.Form.Get("id")
	passwordNew := r.Form.Get("new_password")
	passwordOld := r.Form.Get("old_password")

	user := data.GetUser(id)

	if encryption.GenerateHash512(passwordOld, user.Username) == user.PWHash {
		user.PWHash = encryption.GenerateHash512(passwordNew, user.Username)
		data.ChangeUser(user)
		w.WriteHeader(201)
	} else {
		w.WriteHeader(401)
	}
}
func RemoveSessions(w http.ResponseWriter, r *http.Request) {
	valid, status := CheckSession(r)
	if valid {
		err := r.ParseForm()
		if err != nil {
			w.WriteHeader(400)
			return
		}
		id := r.Form.Get("id")

		user := data.GetUser(id)
		user.Sessions = []data.Session{}
		data.ChangeUser(user)
	} else {
		if status == 1 {
			w.WriteHeader(400)
		} else if status == 2 {
			w.WriteHeader(401)
		}
	}
}
func UpgradeAccount(w http.ResponseWriter, r *http.Request) {
	//TODO
}
func DeleteAccount(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(400)
		return
	}
	id := r.Form.Get("id")
	password := r.Form.Get("password")

	user := data.GetUser(id)
	if encryption.GenerateHash512(password, user.Username) == user.PWHash {
		data.RemoveUser(id)
		w.WriteHeader(201)
	} else {
		w.WriteHeader(401)
	}
}
func CreateChannel(w http.ResponseWriter, r *http.Request) {

}
func DeleteChannel(w http.ResponseWriter, r *http.Request) {

}
func SendMessage(w http.ResponseWriter, r *http.Request) {

}
func DeleteMessage(w http.ResponseWriter, r *http.Request) {

}
func GetChannels(w http.ResponseWriter, r *http.Request) {

}
func GetMessages(w http.ResponseWriter, r *http.Request) {

}

func main() {
	fmt.Println(time.Now().UTC().String() + " | Starting server...")

	go CleanSessions()
	fmt.Println(time.Now().UTC().String() + " | Session cleaner started.")

	//API
	http.HandleFunc("/api/account/signup", Signup)
	http.HandleFunc("/api/account/login", Login)
	http.HandleFunc("/api/account/info", ViewInfo)
	http.HandleFunc("/api/account/edit/username", ChangeUsername)
	http.HandleFunc("/api/account/edit/password", ChangePassword)
	http.HandleFunc("/api/account/edit/sessions", RemoveSessions)
	http.HandleFunc("/api/account/edit/upgrade", UpgradeAccount)
	http.HandleFunc("/api/account/edit/delete", DeleteAccount)
	http.HandleFunc("/api/account/channels", GetChannels)
	http.HandleFunc("/api/channel/create", CreateChannel)
	http.HandleFunc("/api/channel/delete", DeleteChannel)
	http.HandleFunc("/api/channel/messages", GetMessages)
	http.HandleFunc("/api/message/send", SendMessage)
	http.HandleFunc("/api/message/delete", DeleteMessage)

	//Pages
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, ".html") {
			http.Redirect(w, r, r.URL.Path[0:len(r.URL.Path)-5], http.StatusFound)
		}
		if r.URL.Path != "/" && r.URL.Path != "/index" && r.URL.Path != "/home" {
			http.ServeFile(w, r, "pages/404.html")
			return
		}
		http.ServeFile(w, r, "pages/index.html")
	})
	http.HandleFunc("/index.css", func(w http.ResponseWriter, r *http.Request) { http.ServeFile(w, r, "pages/index.css") })
	http.HandleFunc("/pricing", func(w http.ResponseWriter, r *http.Request) { http.ServeFile(w, r, "pages/pricing.html") })
	http.HandleFunc("/pricing.css", func(w http.ResponseWriter, r *http.Request) { http.ServeFile(w, r, "pages/pricing.css") })
	http.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) { http.ServeFile(w, r, "pages/about.html") })
	http.HandleFunc("/about.css", func(w http.ResponseWriter, r *http.Request) { http.ServeFile(w, r, "pages/about.css") })
	http.HandleFunc("/webclient", func(w http.ResponseWriter, r *http.Request) { http.ServeFile(w, r, "pages/webclient.html") })
	http.HandleFunc("/webclient.css", func(w http.ResponseWriter, r *http.Request) { http.ServeFile(w, r, "pages/webclient.css") })
	http.HandleFunc("/webclient.js", func(w http.ResponseWriter, r *http.Request) { http.ServeFile(w, r, "pages/webclient.js") })
	http.HandleFunc("/publickey", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "pages/publickey.contact@gomessenger.link-f1884ecfd460d5f66d9fbccd67366d95cfe8d84d.asc")
	})

	//images
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) { http.ServeFile(w, r, "pages/favicon.ico") })

	fmt.Println(time.Now().UTC().String() + " | Pages loaded.")

	//Server
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println(time.Now().UTC().String() + " | Error starting server: " + err.Error())
	}
}
