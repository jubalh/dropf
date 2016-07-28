package main

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"net/http"
)

type Session struct {
	Id       string
	Username string
}

var SessionStore map[string]string

func InitSessionStore() {
	SessionStore = make(map[string]string)
}

func CreateSession(username string) (id string) {
	data := make([]byte, 32)
	rand.Read(data)
	id = base64.StdEncoding.EncodeToString(data)

	SessionStore[id] = username

	return id
}

func GetUsername(id string) (string, error) {
	if SessionStore[id] == "" {
		return "", errors.New("No user for ID")
	} else {
		return SessionStore[id], nil
	}
}

func SessionExistsForRequest(r *http.Request) bool {
	cookie, err := r.Cookie("dropf")
	if err == nil {
		if SessionStore[cookie.Value] == "" {
			return true
		}
	}
	return false
}
