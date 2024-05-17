package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	_ "reflect"

	_ "github.com/go-sql-driver/mysql"
)

const (
	user   = "mysql"
	pass   = "mysql"
	host   = "localhost"
	port   = "3306"
	dbname = "golangdb"
)

var db *sql.DB

type TUser struct {
	Password string
	Email    string
}

func main() {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, pass, host, port, dbname))
	if err != nil {
		db.Close()
		panic(err)
	}

	fmt.Println("Connected to MySql")
	handleRequest()
}

func handleRequest() {
	http.HandleFunc("/", index)
	http.HandleFunc("/signup/", signUp_page)
	http.HandleFunc("/create_user/", newUser_page)

	http.ListenAndServe(":8080", nil)
}

func index(w_page http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html",
		"templates/header.html",
		"templates/footer.html",
		"templates/signup.html",
		"templates/signup_success.html")
	if err != nil {
		panic(err)
	}
	tmpl.ExecuteTemplate(w_page, "index", nil)
}

func (this TUser) newUser() {

	insert, err := db.Query(fmt.Sprintf("INSERT INTO tbl_user (email, password) VALUES ('%s', '%s')", this.Email, this.Password))
	if err != nil {
		panic(err)
	}
	defer insert.Close()
}

func newUser_page(w_page http.ResponseWriter, r *http.Request) {
	var _user TUser
	_user.Email = r.FormValue("inputEmail")
	_user.Password = r.FormValue("inputPassword")

	// _, err := db.Exec("CREATE TABLE IF NOT EXISTS tbl_user (id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY, email TEXT NOT NULL, password TEXT NOT NULL)")
	// if err != nil {
	// 	panic(err)
	// }

	//	var query = fmt.Sprintf(, _user.Email)
	check_user_in_tbl, err := db.Query("SELECT * FROM tbl_user WHERE email = $1", _user.Email)
	if err != nil {
		panic(err)
	}
	var check_email string
	for check_user_in_tbl.Next() {
		check_user_in_tbl.Scan(&check_email)
	}
	defer check_user_in_tbl.Close()

	if check_email == "" {
		_user.newUser()
	} else {
		_user.Email = "ERROR_USER_ALREADY_EXISTS"
	}

	tmpl, err := template.ParseFiles("templates/index.html",
		"templates/header.html",
		"templates/footer.html",
		"templates/signup.html",
		"templates/signup_success.html")
	if err != nil {
		panic(err)
	}
	tmpl.ExecuteTemplate(w_page, "signup_success", _user)
}

func signUp_page(w_page http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html",
		"templates/header.html",
		"templates/footer.html",
		"templates/signup.html")
	if err != nil {
		panic(err)
	}
	tmpl.ExecuteTemplate(w_page, "signup", nil)
}
