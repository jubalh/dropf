package main

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
)

var SessionStore map[string]string

const CookieName string = "dropf"

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

func DestroySession(id string) {
	delete(SessionStore, id)
}

func GetUsername(id string) (string, error) {
	if SessionStore[id] == "" {
		return "", errors.New("No user for ID")
	} else {
		return SessionStore[id], nil
	}
}

func GetSessionId(r *http.Request) (id string, err error) {
	cookie, err := r.Cookie(CookieName)
	if err != nil {
		return "", fmt.Errorf("No cookie named: %s", CookieName)
	}

	if SessionStore[cookie.Value] == "" {
		return "", fmt.Errorf("No session for cookie: %s", cookie.Value)
	}
	return cookie.Value, nil
}
