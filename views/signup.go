package views

import (
	"html/template"
	"log"
	"net/http"
)

func SignUpPage(w http.ResponseWriter, r *http.Request) {
	signup := template.Must(template.ParseFiles("templates/signup.html"))
	err := signup.Execute(w, nil)
	if err != nil { // if there is an error
		log.Print("template executing error: ", err) //log it
	}
}
