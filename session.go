package main

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
)

// SessionStore stores a session, mapping ID and username.
var SessionStore map[string]string

// CookieName is the name our cookie will have
const CookieName string = "dropf"

// InitSessionStore initializes the store of the cookies.
func InitSessionStore() {
	SessionStore = make(map[string]string)
}

// CreateSession creates a new session for a given username.
func CreateSession(username string) (id string) {
	data := make([]byte, 32)
	rand.Read(data)
	id = base64.StdEncoding.EncodeToString(data)

	SessionStore[id] = username

	return id
}

// DestroySession destroys a session given a session ID.
func DestroySession(id string) {
	delete(SessionStore, id)
}

// GetUsername returns a username for a given session ID.
func GetUsername(id string) (string, error) {
	if SessionStore[id] == "" {
		return "", errors.New("No user for ID")
	}
	return SessionStore[id], nil
}

// GetSessionID returns a sessions ID.
func GetSessionID(r *http.Request) (id string, err error) {
	cookie, err := r.Cookie(CookieName)
	if err != nil {
		return "", fmt.Errorf("No cookie named: %s", CookieName)
	}

	if SessionStore[cookie.Value] == "" {
		return "", fmt.Errorf("No session for cookie: %s", cookie.Value)
	}
	return cookie.Value, nil
}
