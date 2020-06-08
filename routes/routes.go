package routes

import (
	"net/http"

	"github.com/KishoreCM/TodoApp/controllers"
	"github.com/KishoreCM/TodoApp/utils/auth"
	"github.com/KishoreCM/TodoApp/views"
	"github.com/gorilla/mux"
)

func Handlers() *mux.Router {

	r := mux.NewRouter().StrictSlash(true)

	r.Use(CommonMiddleware)

	r.HandleFunc("/", controllers.RedirectToHome).Methods("GET")
	r.HandleFunc("/signup", views.SignUpPage)

	r.HandleFunc("/register", controllers.CreateUser).Methods("POST")
	r.HandleFunc("/login", controllers.Login).Methods("GET", "POST")

	// Auth route
	s := r.PathPrefix("/todo").Subrouter()
	s.Use(auth.JwtVerify)
	s.HandleFunc("/todos", controllers.GetTodos).Methods("GET", "POST")
	s.HandleFunc("/createTodo", controllers.CreateTodo).Methods("POST")
	s.HandleFunc("/updateTodo/{id}", controllers.UpdateTodo).Methods("POST")
	s.HandleFunc("/deleteTodo/{id}", controllers.DeleteTodo).Methods("DELETE")
	s.HandleFunc("/logout", controllers.Logout).Methods("GET")
	s.HandleFunc("/user", controllers.FetchUsers).Methods("GET")
	s.HandleFunc("/user/{id}", controllers.GetUser).Methods("GET")
	s.HandleFunc("/user/{id}", controllers.UpdateUser).Methods("PUT")
	s.HandleFunc("/user/{id}", controllers.DeleteUser).Methods("DELETE")
	return r
}

// CommonMiddleware --Set content-type
func CommonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//w.Header().Add("Content-Type", "application/json")
		w.Header().Add("Content-Type", "text/html")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Access-Control-Request-Headers, Access-Control-Request-Method, Connection, Host, Origin, User-Agent, Referer, Cache-Control, X-header")
		next.ServeHTTP(w, r)
	})
}
