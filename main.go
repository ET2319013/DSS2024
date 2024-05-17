package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	_ "reflect"

	_ "github.com/go-sql-driver/mysql"
)

// MYSQL CREDS
var dbUser string = "root"
var dbPass string = "admin"
var dbHost string = "localhost"
var dbPort string = "3306"
var dbname string = "golangdb"

var curUser string = ""

//var db *sql.DB

type TUser struct {
	Password string
	Email    string
}

type Joke struct {
	Who       string
	Punchline string
}

func main() {

	db1, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, ""))
	if err != nil {
		panic(err.Error())
	}

	_, err = db1.Exec("CREATE DATABASE  IF NOT EXISTS " + dbname)
	if err != nil {
		panic(err)
	}
	db1.Close()

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbname))
	if err != nil {
		panic(err.Error())
	}

	stmt, err := db.Exec("CREATE TABLE IF NOT EXISTS tbl_user (id INT AUTO_INCREMENT PRIMARY KEY, email VARCHAR(50), password VARCHAR(50))")
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(stmt)

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
		"templates/blog.html",
		"templates/signup.html",
		"templates/signup_success.html")
	if err != nil {
		panic(err)
	}
	tmpl.ExecuteTemplate(w_page, "index", curUser)

	var blogText = "post1"

	tmpl.ExecuteTemplate(w_page, "blog", blogText)

}

func (this TUser) newUser(db *sql.DB) {

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

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbname))
	if err != nil {
		panic(err.Error())
	}
	var query = fmt.Sprintf("SELECT * FROM tbl_user WHERE email = '" + _user.Email + "'")
	check_user_in_tbl, err := db.Query(query)
	if err != nil {
		panic(err)
	}
	var check_email string
	for check_user_in_tbl.Next() {
		check_user_in_tbl.Scan(&check_email)
	}
	defer check_user_in_tbl.Close()

	curUser = _user.Email
	if check_email == "" {
		_user.newUser(db)
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
