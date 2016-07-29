// handlers.go contains all the HTTP handlers for dropf
package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
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
	_, err := GetSessionId(request)
	if err == nil {
		http.Redirect(response, request, "/userspace", 302)
		return
	}
	executeTemplate("index.html", response)
}

// loginHandler is used to log the user in.
func loginHandler(response http.ResponseWriter, request *http.Request) {
	name := request.FormValue("name")
	password := request.FormValue("password")
	target := "/"

	for _, username := range Config.Users {
		if username.Name == name {
			if username.Password == password {
				id := CreateSession(username.Name)

				cookie := &http.Cookie{
					Name:  CookieName,
					Value: id,
				}

				http.SetCookie(response, cookie)
				target = "/userspace"
			} else {
				fmt.Println("Failed login for user: ", name)
			}
		} else {
			fmt.Println("Failed login: ", name)
		}
	}

	http.Redirect(response, request, target, 302)
}

// userspaceHandler shows the users private home area.
func userspaceHandler(response http.ResponseWriter, request *http.Request) {
	id, err := GetSessionId(request)
	if err != nil {
		fmt.Println("Debug: not logged in")
		http.Redirect(response, request, "/", 302)
		return
	}

	username, err := GetUsername(id)
	if err != nil {
		fmt.Println("Debug: ", err)
		return
	}

	fmt.Println("Debug: Logged in as:", username)

	executeTemplate("userspace.html", response)
}

// uploadHandler uploads the file.
func uploadHandler(response http.ResponseWriter, request *http.Request) {
	_, err := GetSessionId(request)
	if err != nil {
		http.Redirect(response, request, "/", 302)
		return
	}

	if request.Method != "POST" {
		http.Redirect(response, request, "/userspace", 302)
		return
	}

	file, header, err := request.FormFile("file")
	if err != nil {
		fmt.Fprintln(response, "Something went wrong!")
		fmt.Fprintln(os.Stderr, err)
		return
	}
	defer file.Close()

	output, err := os.Create(filepath.Join(Config.Path, header.Filename))
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

// staticHandler takes care of images and other static files
func staticHandler(response http.ResponseWriter, request *http.Request) {
	if strings.Contains(request.URL.Path, ".png") {
		http.ServeFile(response, request, request.URL.Path[1:])
	}
}
