package main

import (
	"log"
	"net/http"

	"github.com/boltdb/bolt"
)

func main() {
	db, err := bolt.Open("files.db", 0644, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

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
