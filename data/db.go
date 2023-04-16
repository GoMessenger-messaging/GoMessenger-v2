package data

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/Jero075/GoMessenger-V2/encryption"
	"os"
	"time"
)

type User struct {
	Username string        `json:"username"`      //Display name
	ID       string        `json:"id"`            //Unique identifier
	PWHash   string        `json:"pw_hash"`       //Hashed password
	Photo    string        `json:"photo"`         //Profile photo
	Status   string        `json:"status"`        //Status message
	PUBKey   rsa.PublicKey `json:"pub_key"`       //Public Key
	PRIKey   []byte        `json:"pri_key"`       //Private Key
	Access   []string      `json:"access"`        //What private channels the user has access to
	Request  bool          `json:"join_request"`  //If Private channels have to send a request to add this user
	Requests []string      `json:"join_requests"` //Current open requests by private channels
	Premium  bool          `json:"premium"`       //If the user has a premium account
	Sessions []Session     `json:"sessions"`      //The user's sessions
	Recovery string        `json:"recovery"`      //Recovery code
}

type Session struct {
	ID      string    `json:"id"`      //Unique identifier
	Expires time.Time `json:"expires"` //Expiration date
}

type PublicChannel struct {
	Name        string    `json:"name"`        //Display name
	ID          string    `json:"id"`          //Unique identifier
	Photo       string    `json:"photo"`       //Channel photo
	Description string    `json:"description"` //Channel description
	BlockedIDs  []string  `json:"blocked_ids"` //IDs of users that are blocked on this channel
	Admins      []string  `json:"admins"`      //IDs of users that are admins of this channel
	Messages    []Message `json:"messages"`    //Messages in this channel
}

type PrivateChannel struct {
	Name        string    `json:"name"`        //Display name
	ID          string    `json:"id"`          //Unique identifier
	Photo       string    `json:"photo"`       //Channel photo
	Description string    `json:"description"` //Channel description
	AccessIDs   []string  `json:"access_ids"`  //IDs of users that have access to this channel
	Admins      []string  `json:"admins"`      //IDs of users that are admins of this channel
	MaxUsers    int       `json:"max_users"`   //Maximum amount of users that can be in this channel; 0 = unlimited
	Messages    []Message `json:"messages"`    //Messages in this channel
}

type Message struct {
	ID          string    `json:"id"`           //Unique identifier
	SenderID    string    `json:"sender_id"`    //ID of the user that sent this message
	Time        time.Time `json:"time"`         //Time the message was sent
	Content     string    `json:"content"`      //Content of the message
	ReplyTo     string    `json:"reply_to"`     //ID of the message this message is a reply to
	RenderMedia []string  `json:"render_media"` //Media that should be rendered in the message
	Attachments []string  `json:"attachments"`  //Attachments to the message
}

type DB struct {
	Users           []User           `json:"users"`
	PublicChannels  []PublicChannel  `json:"public_channels"`
	PrivateChannels []PrivateChannel `json:"private_channels"`
}

func Idgen(n int) string {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		fmt.Println(time.Now().UTC().String() + " | Error generating random ID: " + err.Error())
	}
	id := hex.EncodeToString(b)
	return id[:n]
}

func OpenDB() DB {
	f, _ := os.ReadFile("data/db.json")
	db := DB{}
	err := json.Unmarshal(f, &db)
	if err != nil {
		fmt.Println(time.Now().UTC().String() + " | Error reading db: " + err.Error())
	}
	return db
}
func SaveDB(db DB) {
	f, _ := json.MarshalIndent(db, "", "	")
	err := os.WriteFile("data/db.json", f, 0644)
	if err != nil {
		fmt.Println(time.Now().UTC().String() + " | Error writing to db: " + err.Error())
	}
}

func AddUser(username string, password string) (userID string) {
	db := OpenDB()

	id := Idgen(8)

	pub, pri := encryption.GenerateKeys()
	priKeyJson, err := json.Marshal(pri)
	if err != nil {
		fmt.Println(time.Now().UTC().String() + " | Error marshalling privateKey: " + err.Error())
	}
	priKeyEnc := encryption.GenerateCiphertext(id, password, priKeyJson)
	db.Users = append(db.Users, User{username, id, encryption.GenerateHash512(password, id), "defaults/user.jpg", "", pub, priKeyEnc, []string{}, false, []string{}, false, []Session{}, ""})

	SaveDB(db)

	return id
}
func RemoveUser(id string) {
	db := OpenDB()

	for i, user := range db.Users {
		if user.ID == id {
			db.Users[i].Username = "[deleted]"
			db.Users[i].PWHash = ""
			db.Users[i].Photo = "defaults/user.jpg"
			db.Users[i].PUBKey = rsa.PublicKey{}
			db.Users[i].PRIKey = nil
			db.Users[i].Premium = false
			db.Users[i].Sessions = []Session{}
			break
		}
	}
	for i, channel := range db.PublicChannels {
		for j, bid := range channel.BlockedIDs {
			if bid == id {
				db.PublicChannels[i].BlockedIDs[j] = db.PublicChannels[i].BlockedIDs[len(db.PublicChannels[i].BlockedIDs)-1]
				db.PublicChannels[i].BlockedIDs = db.PublicChannels[i].BlockedIDs[:len(db.PublicChannels[i].BlockedIDs)-1]
				break
			}
		}
	}
	for i, channel := range db.PrivateChannels {
		for j, aid := range channel.AccessIDs {
			if aid == id {
				db.PrivateChannels[i].AccessIDs[j] = db.PrivateChannels[i].AccessIDs[len(db.PrivateChannels[i].AccessIDs)-1]
				db.PrivateChannels[i].AccessIDs = db.PrivateChannels[i].AccessIDs[:len(db.PrivateChannels[i].AccessIDs)-1]
				break
			}
		}
		for j, aid := range channel.Admins {
			if aid == id {
				db.PrivateChannels[i].Admins[j] = db.PrivateChannels[i].Admins[len(db.PrivateChannels[i].Admins)-1]
				db.PrivateChannels[i].Admins = db.PrivateChannels[i].Admins[:len(db.PrivateChannels[i].Admins)-1]
				break
			}
		}
	}

	SaveDB(db)
}
func ChangeUser(new User) {
	db := OpenDB()

	for i, user := range db.Users {
		if user.ID == new.ID {
			db.Users[i] = new
			break
		}
	}

	SaveDB(db)
}
func GetUser(id string) (userData User) {
	db := OpenDB()

	for _, user := range db.Users {
		if user.ID == id {
			return user
		}
	}
	return User{}
}
func AddPublicChannel(name string, creator string) {
	db := OpenDB()

	db.PublicChannels = append(db.PublicChannels, PublicChannel{name, Idgen(12), "defaults/channel.jpg", "", []string{}, []string{creator}, []Message{}})

	SaveDB(db)
}
func RemovePublicChannel(id string) {
	db := OpenDB()

	for i, channel := range db.PublicChannels {
		if channel.ID == id {
			db.PublicChannels[i] = db.PublicChannels[len(db.PublicChannels)-1]
			db.PublicChannels = db.PublicChannels[:len(db.PublicChannels)-1]
			break
		}
	}

	SaveDB(db)
}
func ChangePublicChannel(new PublicChannel) {
	db := OpenDB()

	for i, channel := range db.PublicChannels {
		if channel.ID == new.ID {
			db.PublicChannels[i] = new
			break
		}
	}

	SaveDB(db)
}
func GetPublicChannel(id string) (channelData PublicChannel) {
	db := OpenDB()

	for _, channel := range db.PublicChannels {
		if channel.ID == id {
			return channel
		}
	}
	return PublicChannel{}
}
func AddPrivateChannel() {
	//TODO
}
func RemovePrivateChannel() {
	//TODO
}
func ChangePrivateChannel() {
	//TODO
}
func GetPrivateChannel() {
	//TODO
}
func AddMessagePublic(channel string, sender string, content string, replyTo string, renderMedia []string, attachments []string) {
	db := OpenDB()

	for i, c := range db.PublicChannels {
		if c.ID == channel {
			db.PublicChannels[i].Messages = append(db.PublicChannels[i].Messages, Message{Idgen(16), sender, time.Now().UTC(), content, replyTo, renderMedia, attachments})
			break
		}
	}

	SaveDB(db)
}
func RemoveMessagePublic(channel string, id string) {
	db := OpenDB()

	for i, c := range db.PublicChannels {
		if c.ID == channel {
			for j, message := range c.Messages {
				if message.ID == id {
					db.PublicChannels[i].Messages[j] = Message{db.PublicChannels[i].Messages[j].ID, db.PublicChannels[i].Messages[j].SenderID, db.PublicChannels[i].Messages[j].Time, "[deleted]", db.PublicChannels[i].Messages[j].ReplyTo, []string{}, []string{}}
					break
				}
			}
		}
	}

	SaveDB(db)
}
func AddMessagePrivate(channel string, sender string, content string, replyTo string, renderMedia []string, attachments []string) {
	//TODO
}
func RemoveMessagePrivate(channel string, id string) {
	//TODO
}
