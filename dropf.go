package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
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

const userspacePage = `
<h1>upload a file</h1>
<form enctype="multipart/form-data" action="/upload" method="post">
      <input type="file" name="file" id="file"/>
      <input type="submit" value="upload" />
</form>
`

// indexHandler serves the index Page where a user can login.
func indexHandler(response http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(response, indexPage)
}

// loginHandler is used to log the user in.
func loginHandler(response http.ResponseWriter, request *http.Request) {
	name := request.FormValue("name")
	password := request.FormValue("password")
	target := "/"
	if name == "gandalf" && password == "mellon" {
		target = "/userspace"
	}
	http.Redirect(response, request, target, 302)
}

// userspaceHandler shows the users private home area.
func userspaceHandler(response http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(response, userspacePage)
}

// uploadHandler uploads the file.
func uploadHandler(response http.ResponseWriter, request *http.Request) {
	if request.Method == "POST" {
		file, header, err := request.FormFile("file")
		if err != nil {
			fmt.Fprintln(response, "Something went wrong!")
			fmt.Println(err)
			return
		}
		defer file.Close()

		os.Mkdir("files", 0660)
		output, err := os.Create("files/" + header.Filename) //TODO: read config where to save
		if err != nil {
			fmt.Fprintln(response, "Something went wrong!")
			fmt.Println(err)
			return
		}
		defer output.Close()

		_, err = io.Copy(output, file)
		if err != nil {
			fmt.Fprintln(response, err)
			fmt.Println(err)
		}

		fmt.Fprintf(response, "File uploaded successfully: ")
		fmt.Fprintf(response, header.Filename)
		fmt.Println("Uploaded file:", header.Filename)
	}
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/userspace", userspaceHandler)
	http.HandleFunc("/upload", uploadHandler)

	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
