// handlers.go contains all the HTTP handlers for dropf
package main

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Hols information about the uploaded files
type File struct {
	Name    string
	Size    int64
	ModTime string
}

// A Filler type to fill in the template.
// It holds the infos that one wants in the userspace
type Filler struct {
	Username string
	Files    []File
}

// executeTemplate loads a basic HTML file and writes it to a writer.
// Which usually should be a http.ResponseWriter.
func executeTemplate(templateName string, writer http.ResponseWriter, filler *Filler) {
	t, err := template.ParseFiles("templates/" + templateName)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	t.Execute(writer, filler)
}

// indexHandler serves the index Page where a user can login.
func indexHandler(response http.ResponseWriter, request *http.Request) {
	_, err := GetSessionId(request)
	if err == nil {
		http.Redirect(response, request, "/userspace", 302)
		return
	}
	executeTemplate("index.html", response, nil)
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

func logoutHandler(response http.ResponseWriter, request *http.Request) {
	id, err := GetSessionId(request)
	if err == nil {
		DestroySession(id)
	}
	http.Redirect(response, request, "/", 302)
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

	f := []File{}

	files, err := ioutil.ReadDir(filepath.Join(Config.Path, username))
	if err != nil {
		fmt.Fprintln(os.Stderr, "User does not have uploaded any files yet")
	}

	for _, file := range files {
		newf := File{Name: file.Name(), Size: file.Size()}
		newf.ModTime = file.ModTime().Format("2006-01-02 10:10")
		f = append(f, newf)
	}
	filler := &Filler{Files: f, Username: username}

	executeTemplate("userspace.html", response, filler)
}

// uploadHandler uploads the file.
func uploadHandler(response http.ResponseWriter, request *http.Request) {
	id, err := GetSessionId(request)
	if err != nil {
		http.Redirect(response, request, "/", 302)
		return
	}

	if request.Method != "POST" {
		http.Redirect(response, request, "/userspace", 302)
		return
	}

	username, err := GetUsername(id)
	// shouldnt happen but always check errors
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	os.Mkdir(filepath.Join(Config.Path, username), 0750)

	request.ParseMultipartForm(32 << 20)

	for _, fh := range request.MultipartForm.File["ufiles"] {
		f, err := fh.Open()
		if err != nil {
			fmt.Fprintln(response, "Something went wrong!")
			fmt.Fprintln(os.Stderr, err)
			return
		}

		defer f.Close()

		output, err := os.Create(filepath.Join(Config.Path, username, fh.Filename))
		if err != nil {
			fmt.Fprintln(response, "Something went wrong!")
			fmt.Fprintln(os.Stderr, err)
			return
		}
		defer output.Close()

		_, err = io.Copy(output, f)
		if err != nil {
			fmt.Fprintln(response, err)
			fmt.Fprintln(os.Stderr, err)
		}

		fmt.Println("Uploaded file:", fh.Filename)
	}
	http.Redirect(response, request, "/userspace", 302)
}

// staticHandler takes care of images and other static files
func staticHandler(response http.ResponseWriter, request *http.Request) {
	if strings.Contains(request.URL.Path, ".png") || strings.Contains(request.URL.Path, ".css") {
		fmt.Println(request.URL.Path[1:])
		http.ServeFile(response, request, request.URL.Path[1:])
	}
}

// fileHandler takes care of getting/viewing or deleting a file
// View request: /file/view/username/filename
// Delete request: /file/delete/username/filename
func fileHandler(response http.ResponseWriter, request *http.Request) {
	id, err := GetSessionId(request)
	if err != nil {
		// TODO redirect ?
		http.Redirect(response, request, "/", 302)
		fmt.Fprintln(response, err)
		return
	}
	loggedin, err := GetUsername(id)
	if err != nil {
		// TODO redirect ?
		fmt.Fprintln(response, err)
		return
	}

	r, err := regexp.Compile("^/file/(view|delete)/" + loggedin + "/")
	if err != nil {
		fmt.Fprintln(os.Stderr, "ERROR regex: ", err)
		return
	}

	if r.MatchString(request.URL.Path) == false {
		log.Printf("UNAUTHORIZED: %s requires %s\n", loggedin, request.URL.Path)
		// TODO redirect ?
		http.Redirect(response, request, "/", 302)
		return
	}

	log.Printf("%s requires %s\n", loggedin, request.URL.Path)

	if strings.Contains(request.URL.Path, "/file/view/") {
		req := request.URL.Path[11:]
		if strings.Compare(req[0:len(loggedin)+1], loggedin+"/") == 0 {
			http.ServeFile(response, request, "files/"+req)
		} else {
			http.Redirect(response, request, "/", 404)
		}
	} else if strings.Contains(request.URL.Path, "/file/delete/") {
		file, err := url.PathUnescape(request.URL.String()[13:])
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error unescaping string: ")
			fmt.Fprintln(os.Stderr, err)
		}
		fmt.Println("file", file)
		fmt.Println(file)

		os.Remove("files/" + file)

		http.Redirect(response, request, "/userspace", 302)
	}
}
