package main

import (
	"fmt"
	"log"
	"net/http"
)

const indexPage = `
<h1>dropf<h2>
<p>Drop your files here</p>
<form method="post" action="/login">
    <label for="name">User name</label>
    <input type="text" id="name" name="name">
    <label for="password">Password</label>
    <input type="password" id="password" name="password">
    <button type="submit">Login</button>
</form>
`

func indexHandler(response http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(response, indexPage)
}

func loginHandler(response http.ResponseWriter, request *http.Request) {
	name := request.FormValue("name")
	password := request.FormValue("password")
	target := "/"
	if name == "gandalf" && password == "mellon" {
		target = "/userspace"
	}
	http.Redirect(response, request, target, 302)
}

func userspaceHandler(response http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(response, "good")
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/userspace", userspaceHandler)
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
