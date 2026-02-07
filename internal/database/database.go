package database

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	bolt "go.etcd.io/bbolt"
)

type Database struct {
	db *bolt.DB
}

type User struct {
	ID           int64     `json:"id"`
	Username     string    `json:"username"`
	FirstName    string    `json:"first_name"`
	Limit        int       `json:"limit"`
	Exp          int       `json:"exp"`
	Level        int       `json:"level"`
	Premium      bool      `json:"premium"`
	PremiumUntil time.Time `json:"premium_until"`
	Registered   bool      `json:"registered"`
	RegisteredAt time.Time `json:"registered_at"`
	LastSeen     time.Time `json:"last_seen"`
}

type Chat struct {
	ID      int64  `json:"id"`
	Type    string `json:"type"`
	Title   string `json:"title"`
	Welcome bool   `json:"welcome"`
	Muted   bool   `json:"muted"`
}

func New(path string) (*Database, error) {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	db, err := bolt.Open(path, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, err
	}

	err = db.Update(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists([]byte("users")); err != nil {
			return err
		}
		if _, err := tx.CreateBucketIfNotExists([]byte("chats")); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		db.Close()
		return nil, err
	}

	return &Database{db: db}, nil
}

func (d *Database) Close() error {
	return d.db.Close()
}

func (d *Database) GetUser(userID int64) (*User, error) {
	var user User
	err := d.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("users"))
		data := b.Get(itob(userID))
		if data == nil {
			return nil
		}
		return json.Unmarshal(data, &user)
	})
	return &user, err
}

func (d *Database) SaveUser(user *User) error {
	return d.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("users"))
		data, err := json.Marshal(user)
		if err != nil {
			return err
		}
		return b.Put(itob(user.ID), data)
	})
}

func (d *Database) GetOrCreateUser(userID int64, username, firstName string) (*User, error) {
	user, err := d.GetUser(userID)
	if err != nil {
		return nil, err
	}

	if user.ID == 0 {
		user = &User{
			ID:           userID,
			Username:     username,
			FirstName:    firstName,
			Limit:        30,
			Premium:      false,
			Registered:   false,
			RegisteredAt: time.Now(),
			LastSeen:     time.Now(),
		}
		if err := d.SaveUser(user); err != nil {
			return nil, err
		}
	} else {
		user.LastSeen = time.Now()
		d.SaveUser(user)
	}

	return user, nil
}

func (d *Database) GetChat(chatID int64) (*Chat, error) {
	var chat Chat
	err := d.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("chats"))
		data := b.Get(itob(chatID))
		if data == nil {
			return nil
		}
		return json.Unmarshal(data, &chat)
	})
	return &chat, err
}

func (d *Database) SaveChat(chat *Chat) error {
	return d.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("chats"))
		data, err := json.Marshal(chat)
		if err != nil {
			return err
		}
		return b.Put(itob(chat.ID), data)
	})
}

func itob(v int64) []byte {
	return []byte(fmt.Sprintf("%d", v))
}

func (d *Database) GetAllUsers() []*User {
	var users []*User
	d.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("users"))
		b.ForEach(func(k, v []byte) error {
			var user User
			if err := json.Unmarshal(v, &user); err == nil {
				users = append(users, &user)
			}
			return nil
		})
		return nil
	})
	return users
}

func (d *Database) GetTotalUsers() int {
	count := 0
	d.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("users"))
		if b != nil {
			count = b.Stats().KeyN
		}
		return nil
	})
	return count
}
