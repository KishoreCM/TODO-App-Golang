package controllers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/KishoreCM/TodoApp/models"
	"github.com/KishoreCM/TodoApp/views"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

type Response struct {
	Message string
}

func IsEmpty(data string) bool {
	if len(data) == 0 {
		return true
	} else {
		return false
	}
}

func displaySignUpPage(w http.ResponseWriter, message Response) {
	signup := template.Must(template.ParseFiles("templates/signup.html"))
	if templateErr := signup.ExecuteTemplate(w, "signup.html", message); templateErr != nil { // if there is an error
		log.Print("template executing error: ", templateErr) //log it
	}
}

func displayLoginPage(w http.ResponseWriter, message Response) {
	signup := template.Must(template.ParseFiles("templates/login.html"))
	if templateErr := signup.ExecuteTemplate(w, "login.html", message); templateErr != nil { // if there is an error
		log.Print("template executing error: ", templateErr) //log it
	}
}

func RedirectToHome(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/signup", http.StatusSeeOther)
}

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		views.LoginPage(w, r)
	} else {
		r.ParseForm()
		email, password := r.FormValue("email"), r.FormValue("password")
		//err := json.NewDecoder(r.Body).Decode(user)
		if IsEmpty(email) || IsEmpty(password) {
			//var res = map[string]interface{}{"message": "Invalid request"}
			//json.NewEncoder(w).Encode(res)
			err := Response{
				Message: "Please Enter the Missing Credentials",
			}
			displayLoginPage(w, err)
			return
		}
		FindOne(w, r, email, password)
	}
}

func FindOne(w http.ResponseWriter, r *http.Request, email, password string) {
	user := &models.User{}

	if err := db.Where("Email = ?", email).First(user).Error; err != nil {
		var err = Response{
			Message: "Email address not found",
		}
		displayLoginPage(w, err)
		return
	}

	errf := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if errf != nil && errf == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
		var err = Response{
			Message: "Invalid login credentials. Please try again",
		}
		displayLoginPage(w, err)
		return
	}

	expirationTime := time.Now().Add(30 * time.Second)
	cookieExpirationTime := time.Now().Add(40 * time.Second)

	claims := &models.Claims{
		UserID: user.ID,
		Name:   user.Name,
		Email:  user.Email,
		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), claims)

	tokenString, error := token.SignedString([]byte("secret"))
	if error != nil {
		fmt.Println(error)
		err := Response{
			Message: "JWT Signing Failed",
		}
		fmt.Fprintln(w, err)
	}

	type ResponseWithToken struct {
		Message string
		Token   string
	}
	/*
		res := ResponseWithToken{
			Message: "logged in!",
			Token:   tokenString,
		}
	*/
	http.SetCookie(w, &http.Cookie{
		Name:    "todoToken",
		Value:   tokenString,
		Expires: cookieExpirationTime,
	})

	//fmt.Fprintln(w, res)
	http.Redirect(w, r, "/todo/todos", http.StatusSeeOther)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	clearCookie := &http.Cookie{
		Name:   "todoToken",
		Value:  "",
		MaxAge: -1,
	}
	http.SetCookie(w, clearCookie)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

//CreateUser function -- create a new user
func CreateUser(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	name, email, password, confirmPassword := r.FormValue("name"), r.FormValue("email"), r.FormValue("password"), r.FormValue("confirmpassword")

	if IsEmpty(name) || IsEmpty(email) || IsEmpty(password) || IsEmpty(confirmPassword) {
		//fmt.Fprintf(w, "Please Enter the Missing Credentials")
		err := Response{
			Message: "Please Enter the Missing Credentials",
		}
		displaySignUpPage(w, err)
		return
	}

	if password != confirmPassword {
		//fmt.Fprintf(w, "Password Mismatch")
		err := Response{
			Message: "Password Mismatch",
		}
		displaySignUpPage(w, err)
		return
	}

	user := &models.User{
		Name:  name,
		Email: email,
	}
	//json.NewDecoder(r.Body).Decode(user)

	pass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
		err := Response{
			Message: "Password Encryption failed",
		}
		json.NewEncoder(w).Encode(err)
	}

	user.Password = string(pass)

	createdUser := db.Create(user)
	var errMessage = createdUser.Error

	if createdUser.Error != nil {
		fmt.Println(errMessage)
	}
	res := Response{
		Message: "User Registered Successfully!",
	}
	json.NewEncoder(w).Encode(res)

}

//FetchUser function
func FetchUsers(w http.ResponseWriter, r *http.Request) {

	/* For testing purpose, get the token from cookie and check it is getting refreshed for every 30 seconds

	cookie, err := r.Cookie("todoToken")
	if err != nil {
		if err == http.ErrNoCookie {
			// If the cookie is not set, return an unauthorized status
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
	}

	tokenString := cookie.Value

	if refreshToken := r.Context().Value("refreshToken"); refreshToken != nil {
		var resWithRefreshToken = map[string]interface{}{
			"refreshToken": refreshToken.(string),
			"users":        users,
		}

		json.NewEncoder(w).Encode(resWithRefreshToken)
		return
	}

	var res = map[string]interface{}{
		"token": tokenString,
		"users": users,
	}
	json.NewEncoder(w).Encode(res)
	*/
	var users []models.User
	db.Find(&users)

}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}
	params := mux.Vars(r)
	var id = params["id"]
	db.First(&user, id)
	json.NewDecoder(r.Body).Decode(user)
	db.Save(&user)
	json.NewEncoder(w).Encode(&user)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var id = params["id"]
	var user models.User
	db.First(&user, id)
	db.Delete(&user)
	json.NewEncoder(w).Encode("User deleted")
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var id = params["id"]
	var user models.User
	db.First(&user, id)
	json.NewEncoder(w).Encode(&user)
}
