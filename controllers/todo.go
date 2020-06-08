package controllers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/KishoreCM/TodoApp/models"
	"github.com/KishoreCM/TodoApp/templates"
	"github.com/KishoreCM/TodoApp/utils"
	"github.com/gorilla/mux"
)

var db = utils.ConnectDB()

func GetTodos(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		userId := r.Context().Value("userId").(uint)

		var todos []models.Todo
		db.Where("user_id = ?", userId).Find(&todos)
		checkboxes := templates.CheckBox{
			Todos: todos,
		}

		todoPage := template.Must(template.ParseFiles("templates/todos.html"))
		if templateErr := todoPage.ExecuteTemplate(w, "todos.html", checkboxes); templateErr != nil { // if there is an error
			log.Print("template executing error: ", templateErr) //log it
		}
		return
	}
}

func CreateTodo(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("userId").(uint)
	r.ParseForm()
	newTodo := r.FormValue("newTodo")
	createTodo := &models.Todo{
		UserID:      int(userId),
		Description: newTodo,
		Done:        false,
	}
	db.Create(createTodo)
	fmt.Println("Todo Created!")
	http.Redirect(w, r, "/todo/todos", http.StatusSeeOther)
}

func UpdateTodo(w http.ResponseWriter, r *http.Request) {
	//db.LogMode(true)
	r.ParseForm()
	todoCheckbox := r.FormValue("todoCheckbox")
	params := mux.Vars(r)
	formId := params["id"]
	var todo models.Todo
	if todoCheckbox == "true" {
		db.Model(&todo).Where("id = ?", formId).Update("done", true)
	} else {
		db.Model(&todo).Where("id = ?", formId).Update("done", false)
	}
	http.Redirect(w, r, "/todo/todos", http.StatusSeeOther)
}

func DeleteTodo(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	formId := params["id"]
	var todo models.Todo
	db.Find(&todo, formId)
	db.Delete(&todo)
	http.Redirect(w, r, "/todo/todos", http.StatusSeeOther)
}
