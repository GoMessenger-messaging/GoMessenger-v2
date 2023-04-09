package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

func Signup(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(400)
		return
	}
	username = r.Form.Get("username")
	password = r.Form.Get("password")

	data.AddUser(username, password)
	w.WriteHeader(201)
}
func Login(w http.ResponseWriter, r *http.Request) {

}
func ViewInfo(w http.ResponseWriter, r *http.Request) {

}
func ChangeSettings(w http.ResponseWriter, r *http.Request) {

}
func UpgradeAccount(w http.ResponseWriter, r *http.Request) {

}
func DeleteAccount(w http.ResponseWriter, r *http.Request) {

}
func MakeChannel(w http.ResponseWriter, r *http.Request) {

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

	//API
	http.HandleFunc("/api/account/signup", Signup)
	http.HandleFunc("/api/account/login", Login)
	http.HandleFunc("/api/account/view", ViewInfo)
	http.HandleFunc("/api/account/settings", ChangeSettings)
	http.HandleFunc("/api/account/upgrade", UpgradeAccount)
	http.HandleFunc("/api/account/delete", DeleteAccount)
	http.HandleFunc("/api/account/channels", GetChannels)
	http.HandleFunc("/api/channel/create", MakeChannel)
	http.HandleFunc("/api/channel/delete", DeleteChannel)
	http.HandleFunc("/api/channel/messages", GetMessages)
	http.HandleFunc("/api/message/send", SendMessage)
	http.HandleFunc("/api/message/delete", DeleteMessage)

	//Pages
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, ".html") {
			http.Redirect(w, r, r.URL.Path[0:len(r.URL.Path)-5], http.StatusFound)
		}
		if r.URL.Path != "/" && r.URL.Path != "/index.html" {
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

	//images
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) { http.ServeFile(w, r, "pages/images/favicon.ico") })

	//Server
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println(time.Now().UTC().String() + " | Error starting server: " + err.Error())
	}
}
