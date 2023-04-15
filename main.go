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
	Photo    string `json:"photo"`
	Status   string `json:"status"`
	Premium  bool   `json:"premium"`
}

func CheckSession(r *http.Request) (bool, int) {
	err := r.ParseForm()
	if err != nil {
		return false, 1
	}
	id := r.Form.Get("id")
	session := r.Form.Get("session")
	hash := encryption.GenerateHash256(id, session)
	user := data.GetUser(id)
	for _, s := range user.Sessions {
		if s.ID == hash {
			return true, 0
		}
	}
	return false, 2
}
func CleanSessions() {
	for {
		db := data.OpenDB()
		for i, user := range db.Users {
			for j := len(user.Sessions) - 1; j >= 0; j-- {
				if user.Sessions[j].Expires.Before(time.Now()) {
					user.Sessions[j] = user.Sessions[len(user.Sessions)-1]
					user.Sessions = user.Sessions[:len(user.Sessions)-1]
				}
			}
			db.Users[i] = user
		}
		data.SaveDB(db)
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
	if username != "" && password != "" {
		id := data.AddUser(username, password)
		_, wErr := w.Write([]byte(id))
		if wErr != nil {
			fmt.Println(time.Now().UTC().String() + " | Error responding to Signup request: " + err.Error())
		}
	} else {
		w.WriteHeader(400)
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

	if encryption.GenerateHash512(password, user.ID) == user.PWHash {
		session := data.Idgen(32)
		hash := encryption.GenerateHash256(id, session)
		user.Sessions = append(user.Sessions, data.Session{ID: hash, Expires: time.Now().Add(time.Hour)})
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
	resp := ViewAccount{Username: user.Username, ID: user.ID, Photo: user.Photo, Status: user.Status, Premium: user.Premium}
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
		if username != "" {
			user := data.GetUser(id)
			user.Username = username
			data.ChangeUser(user)
			w.WriteHeader(201)
		} else {
			w.WriteHeader(400)
		}
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
	if passwordNew != "" {
		user := data.GetUser(id)
		if encryption.GenerateHash512(passwordOld, user.ID) == user.PWHash {
			user.PWHash = encryption.GenerateHash512(passwordNew, user.ID)
			data.ChangeUser(user)
			w.WriteHeader(201)
		} else {
			w.WriteHeader(401)
		}
	} else {
		w.WriteHeader(400)
	}
}
func ChangePhoto(w http.ResponseWriter, r *http.Request) {
	//TODO
}
func ChangeStatus(w http.ResponseWriter, r *http.Request) {
	//TODO
}
func ChangeRecoveryCode(w http.ResponseWriter, r *http.Request) {
	//TODO
}
func JoinChannel(w http.ResponseWriter, r *http.Request) {
	valid, status := CheckSession(r)
	if valid {
		err := r.ParseForm()
		if err != nil {
			w.WriteHeader(400)
			return
		}
		id := r.Form.Get("id")
		channelId := r.Form.Get("channel")
		channel := data.GetPublicChannel(channelId)
		if channel.ID != "" {
			user := data.GetUser(id)
			user.Access = append(user.Access, channel.ID)
			data.ChangeUser(user)
			w.WriteHeader(201)
		} else {
			w.WriteHeader(400)
		}
	} else {
		if status == 1 {
			w.WriteHeader(400)
		} else if status == 2 {
			w.WriteHeader(401)
		}
	}
}
func LeavePublicChannel(w http.ResponseWriter, r *http.Request) {
	valid, status := CheckSession(r)
	if valid {
		err := r.ParseForm()
		if err != nil {
			w.WriteHeader(400)
			return
		}
		id := r.Form.Get("id")
		channelId := r.Form.Get("channel")
		channel := data.GetPublicChannel(channelId)
		if channel.ID != "" {
			user := data.GetUser(id)
			for i, id := range user.Access {
				if id == channel.ID {
					user.Access[i] = user.Access[len(user.Access)-1]
					user.Access = user.Access[:len(user.Access)-1]
					break
				}
			}
			data.ChangeUser(user)
			w.WriteHeader(201)
		} else {
			w.WriteHeader(400)
		}
	} else {
		if status == 1 {
			w.WriteHeader(400)
		} else if status == 2 {
			w.WriteHeader(401)
		}
	}
}
func LeavePrivateChannel(w http.ResponseWriter, r *http.Request) {
	//TODO
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
func RecoverAccount(w http.ResponseWriter, r *http.Request) {
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
	if encryption.GenerateHash512(password, user.ID) == user.PWHash {
		data.RemoveUser(id)
		w.WriteHeader(201)
	} else {
		w.WriteHeader(401)
	}
}
func GetChannels(w http.ResponseWriter, r *http.Request) {

}
func CreatePublicChannel(w http.ResponseWriter, r *http.Request) {
	valid, status := CheckSession(r)
	if valid {
		err := r.ParseForm()
		if err != nil {
			w.WriteHeader(400)
			return
		}
		id := r.Form.Get("id")
		name := r.Form.Get("name")
		if name != "" && data.GetUser(id).ID != "" {
			data.AddPublicChannel(name, id)
			w.WriteHeader(201)
		} else {
			w.WriteHeader(400)
		}
	} else {
		if status == 1 {
			w.WriteHeader(400)
		} else if status == 2 {
			w.WriteHeader(401)
		}
	}
}
func CreatePrivateChannel(w http.ResponseWriter, r *http.Request) {
	//TODO
}
func ChangePublicChannelName(w http.ResponseWriter, r *http.Request) {
	//TODO
}
func ChangePrivateChannelName(w http.ResponseWriter, r *http.Request) {
	//TODO
}
func ChangePublicChannelPhoto(w http.ResponseWriter, r *http.Request) {
	//TODO
}
func ChangePrivateChannelPhoto(w http.ResponseWriter, r *http.Request) {
	//TODO
}
func ChangePublicChannelDescription(w http.ResponseWriter, r *http.Request) {
	//TODO
}
func ChangePrivateChannelDescription(w http.ResponseWriter, r *http.Request) {
	//TODO
}
func ChangePublicChannelMembers(w http.ResponseWriter, r *http.Request) {
	//TODO
}
func ChangePrivateChannelMembers(w http.ResponseWriter, r *http.Request) {
	//TODO
}
func ChangePublicChannelAdmins(w http.ResponseWriter, r *http.Request) {
	//TODO
}
func ChangePrivateChannelAdmins(w http.ResponseWriter, r *http.Request) {
	//TODO
}
func DeletePublicChannel(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(400)
		return
	}
	id := r.Form.Get("id")
	password := r.Form.Get("password")
	channelId := r.Form.Get("channel")
	channel := data.GetPublicChannel(channelId)
	if channel.ID != "" {
		user := data.GetUser(id)
		if encryption.GenerateHash512(password, user.ID) == user.PWHash {
			data.RemovePublicChannel(channelId)
			w.WriteHeader(201)
		} else {
			w.WriteHeader(401)
		}
	} else {
		w.WriteHeader(400)
	}
}
func DeletePrivateChannel(w http.ResponseWriter, r *http.Request) {
	//TODO
}
func SendPublicMessage(w http.ResponseWriter, r *http.Request) {

}
func SendPrivateMessage(w http.ResponseWriter, r *http.Request) {
	//TODO
}
func DeletePublicMessage(w http.ResponseWriter, r *http.Request) {

}
func DeletePrivateMessage(w http.ResponseWriter, r *http.Request) {
	//TODO
}
func GetPublicMessages(w http.ResponseWriter, r *http.Request) {

}
func GetPrivateMessages(w http.ResponseWriter, r *http.Request) {
	//TODO
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
	http.HandleFunc("/api/account/edit/photo", ChangePhoto)
	http.HandleFunc("/api/account/edit/status", ChangeStatus)
	http.HandleFunc("/api/account/edit/recovery", ChangeRecoveryCode)
	http.HandleFunc("/api/account/edit/join", JoinChannel)
	http.HandleFunc("/api/account/edit/leave/public", LeavePublicChannel)
	http.HandleFunc("/api/account/edit/leave/private", LeavePrivateChannel)
	http.HandleFunc("/api/account/edit/sessions", RemoveSessions)
	http.HandleFunc("/api/account/edit/upgrade", UpgradeAccount)
	http.HandleFunc("/api/account/edit/recover", RecoverAccount)
	http.HandleFunc("/api/account/edit/delete", DeleteAccount)
	http.HandleFunc("/api/account/channels", GetChannels)
	http.HandleFunc("/api/channel/create/public", CreatePublicChannel)
	http.HandleFunc("/api/channel/create/private", CreatePrivateChannel)
	http.HandleFunc("/api/channel/edit/name/public", ChangePublicChannelName)
	http.HandleFunc("/api/channel/edit/name/private", ChangePrivateChannelName)
	http.HandleFunc("/api/channel/edit/photo/public", ChangePublicChannelPhoto)
	http.HandleFunc("/api/channel/edit/photo/private", ChangePrivateChannelPhoto)
	http.HandleFunc("/api/channel/edit/description/public", ChangePublicChannelDescription)
	http.HandleFunc("/api/channel/edit/description/private", ChangePrivateChannelDescription)
	http.HandleFunc("/api/channel/edit/members/public", ChangePublicChannelMembers)
	http.HandleFunc("/api/channel/edit/members/private", ChangePrivateChannelMembers)
	http.HandleFunc("/api/channel/edit/admins/public", ChangePublicChannelAdmins)
	http.HandleFunc("/api/channel/edit/admins/private", ChangePrivateChannelAdmins)
	http.HandleFunc("/api/channel/edit/delete/public", DeletePublicChannel)
	http.HandleFunc("/api/channel/edit/delete/private", DeletePrivateChannel)
	http.HandleFunc("/api/message/send/public", SendPublicMessage)
	http.HandleFunc("/api/message/send/private", SendPrivateMessage)
	http.HandleFunc("/api/message/delete/public", DeletePublicMessage)
	http.HandleFunc("/api/message/delete/private", DeletePrivateMessage)
	http.HandleFunc("/api/channel/messages/public", GetPublicMessages)
	http.HandleFunc("/api/channel/messages/private", GetPrivateMessages)

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
	http.Handle("/uploads/", http.FileServer(http.Dir("")))
	http.HandleFunc("/publickey-contact@gomessenger.link.asc", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/pgp-keys")
		http.ServeFile(w, r, "pages/publickey.contact@gomessenger.link-f1884ecfd460d5f66d9fbccd67366d95cfe8d84d.asc")
	})
	http.HandleFunc("/robots.txt", func(w http.ResponseWriter, r *http.Request) { http.ServeFile(w, r, "pages/robots.txt") })

	//images
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) { http.ServeFile(w, r, "pages/favicon.ico") })

	fmt.Println(time.Now().UTC().String() + " | Pages loaded.")

	//Server
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println(time.Now().UTC().String() + " | Error starting server: " + err.Error())
	}
}
