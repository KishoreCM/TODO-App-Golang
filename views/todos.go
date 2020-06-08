package views

import (
	"html/template"
	"log"
	"net/http"
)

func TodosPage(w http.ResponseWriter, r *http.Request) {
	todos := template.Must(template.ParseFiles("templates/todos.html"))
	err := todos.Execute(w, nil)
	if err != nil { // if there is an error
		log.Print("template executing error: ", err) //log it
	}
}
