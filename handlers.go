// handlers.go contains all the HTTP handlers for dropf
package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"strings"
)

// executeTemplate loads a basic HTML file and writes it to a writer.
// Which usually should be a http.ResponseWriter.
func executeTemplate(templateName string, writer http.ResponseWriter) {
	t, err := template.ParseFiles("templates/" + templateName)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	t.Execute(writer, nil)
}

// indexHandler serves the index Page where a user can login.
func indexHandler(response http.ResponseWriter, request *http.Request) {
	executeTemplate("index.html", response)
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
	executeTemplate("userspace.html", response)
}

// uploadHandler uploads the file.
func uploadHandler(response http.ResponseWriter, request *http.Request) {
	if request.Method == "POST" {
		file, header, err := request.FormFile("file")
		if err != nil {
			fmt.Fprintln(response, "Something went wrong!")
			fmt.Fprintln(os.Stderr, err)
			return
		}
		defer file.Close()

		os.Mkdir("files", 0660)
		output, err := os.Create("files/" + header.Filename) //TODO: read config where to save
		if err != nil {
			fmt.Fprintln(response, "Something went wrong!")
			fmt.Fprintln(os.Stderr, err)
			return
		}
		defer output.Close()

		_, err = io.Copy(output, file)
		if err != nil {
			fmt.Fprintln(response, err)
			fmt.Fprintln(os.Stderr, err)
		}

		fmt.Fprintf(response, "File uploaded successfully: ")
		fmt.Fprintf(response, header.Filename)
		fmt.Println("Uploaded file:", header.Filename)
	}
}

// staticHandler takes care of images and other static files
func staticHandler(response http.ResponseWriter, request *http.Request) {
	if strings.Contains(request.URL.Path, ".png") {
		http.ServeFile(response, request, request.URL.Path[1:])
	}
}
