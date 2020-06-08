package views

import (
	"html/template"
	"log"
	"net/http"
)

func LoginPage(w http.ResponseWriter, r *http.Request) {
	signup := template.Must(template.ParseFiles("templates/login.html"))
	err := signup.Execute(w, nil)
	if err != nil { // if there is an error
		log.Print("template executing error: ", err) //log it
	}
}
