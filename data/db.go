package data

import (
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/Jero075/GoMessenger-V2/encryption"
	"os"
	"time"
)

type User struct {
	Username string           `json:"username"`
	ID       string           `json:"id"`
	PWHash   string           `json:"pw_hash"`
	Photo    string           `json:"photo"`
	PUBKey   ecdsa.PublicKey  `json:"pub_key"`
	PRIKey   ecdsa.PrivateKey `json:"pri_key"`
	Premium  bool             `json:"premium"`
	Sessions []Session        `json:"sessions"`
}

type Session struct {
	ID      string    `json:"id"`
	Expires time.Time `json:"expires"`
}

type PublicChannel struct {
	Name       string    `json:"name"`
	ID         string    `json:"id"`
	Photo      string    `json:"photo"`
	BlockedIDs []string  `json:"blocked_ids"`
	Admins     []string  `json:"admins"`
	Messages   []Message `json:"messages"`
}

type PrivateChannel struct {
	Name      string    `json:"name"`
	ID        string    `json:"id"`
	Photo     string    `json:"photo"`
	AccessIDs []string  `json:"access_ids"`
	Admins    []string  `json:"admins"`
	Messages  []Message `json:"messages"`
}

type Message struct {
	ID          string    `json:"id"`
	SenderID    string    `json:"sender_id"`
	Time        time.Time `json:"time"`
	Content     string    `json:"content"`
	ReplyTo     string    `json:"reply_to"`
	RenderMedia []string  `json:"render_media"`
	Attachments []string  `json:"attachments"`
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

	pub, pri := encryption.GenerateKeys(username, password)
	id := Idgen(8)
	db.Users = append(db.Users, User{username, id, encryption.GenerateHash512(password, username), "", pub, pri, false, []Session{}})

	SaveDB(db)

	return id
}
func RemoveUser(id string) {
	db := OpenDB()

	for i, user := range db.Users {
		if user.ID == id {
			db.Users[i] = db.Users[len(db.Users)-1]
			db.Users = db.Users[:len(db.Users)-1]
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

	db.PublicChannels = append(db.PublicChannels, PublicChannel{name, Idgen(12), "", []string{}, []string{creator}, []Message{}})

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
func AddPrivateChannel(name string, creator string) {
	db := OpenDB()

	db.PrivateChannels = append(db.PrivateChannels, PrivateChannel{name, Idgen(12), "", []string{creator}, []string{creator}, []Message{}})

	SaveDB(db)
}
func RemovePrivateChannel(id string) {
	db := OpenDB()

	for i, channel := range db.PrivateChannels {
		if channel.ID == id {
			db.PrivateChannels[i] = db.PrivateChannels[len(db.PrivateChannels)-1]
			db.PrivateChannels = db.PrivateChannels[:len(db.PrivateChannels)-1]
			break
		}
	}

	SaveDB(db)
}
func ChangePrivateChannel(new PrivateChannel) {
	db := OpenDB()

	for i, channel := range db.PrivateChannels {
		if channel.ID == new.ID {
			db.PrivateChannels[i] = new
			break
		}
	}

	SaveDB(db)
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
