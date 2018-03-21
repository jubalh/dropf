package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
)

// User holds information of the username and the password.
// Both are stored in plain text in the config.json file.
type User struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

// Config holds information about the content of the config.json file.
// Path is the path to where the files are stored.
var Config struct {
	Path  string `json:"path"`
	Users []User `json:"users"`
}

// readConfig reads the config.json file into the structs.
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

	port := flag.String("port", "9090", "Port to start service on")
	flag.Parse()

	// Create directory where files will be saved
	if Config.Path == "" {
		fmt.Println("'Path' is not defined in configuration. Fallback to 'files'")
		Config.Path = "files"
	}

	err = os.Mkdir(Config.Path, 0750)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Can notnot create directory: ", err)
		os.Exit(1)
	}

	InitSessionStore()

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/logout", logoutHandler)
	http.HandleFunc("/userspace", userspaceHandler)
	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/static/", staticHandler)
	http.HandleFunc("/file/", fileHandler)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	go func() {
		log.Println("Starting to listen on port", *port)
		err = http.ListenAndServe(":"+*port, nil)
		if err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
	}()

	<-stop
	fmt.Println()
	log.Println("dropf stopped")
}
