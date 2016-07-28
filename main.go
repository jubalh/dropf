package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/boltdb/bolt"
)

type User struct {
	Name, Password string
}

var Config struct {
	Path  string
	Users []User
}

func readConfig() error {
	f, err := os.Open("config.json")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Can not read configuration file: ", err)
		return err
	}
	defer f.Close()

	decoder := json.NewDecoder(f)
	err = decoder.Decode(&Config)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Can not decode json configuration file: ", err)
		return err
	}
	return nil
}

func main() {
	err := readConfig()
	if err != nil {
		os.Exit(1)
	}

	// Create directory where files will be saved
	if Config.Path == "" {
		fmt.Println("'Path' is not defined in configuration. Fallback to 'files'")
		Config.Path = "files"
	}
	os.Mkdir(Config.Path, 0750)

	db, err := bolt.Open("files.db", 0644, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	InitSessionStore()

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/userspace", userspaceHandler)
	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/static/", staticHandler)

	err = http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
