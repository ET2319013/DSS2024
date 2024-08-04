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
var curID int = -1

//TODO: post comments
//var db *sql.DB

type TUser struct {
	Password string
	Email    string
}

type TPost struct {
	Email   string
	Title   string
	Content string
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

	stmt1, err := db.Exec("CREATE TABLE IF NOT EXISTS tbl_posts (id INT AUTO_INCREMENT PRIMARY KEY, email VARCHAR(50), title VARCHAR(50), post VARCHAR(1024))")
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(stmt1)

	handleRequest()
}

func handleRequest() {
	http.HandleFunc("/", index)
	http.HandleFunc("/signup/", signUp_page)
	http.HandleFunc("/logout/", index_logged_out)
	http.HandleFunc("/create_user/", newUser_page)
	http.HandleFunc("/post/", newPost_page)
	http.HandleFunc("/create_post/", newPost)

	http.ListenAndServe(":8080", nil)
}

func index(w_page http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html",
		"templates/header.html",
		"templates/footer.html",
		"templates/blog.html",
		"templates/signup.html",
		"templates/post.html",
		"templates/signup_success.html")
	if err != nil {
		panic(err)
	}
	tmpl.ExecuteTemplate(w_page, "index", curUser)

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbname))
	if err != nil {
		panic(err.Error())
	}
	var query = fmt.Sprintf("SELECT post FROM tbl_posts WHERE email = '" + curUser + "'")
	check_post_in_tbl, err := db.Query(query)
	if err != nil {
		panic(err)
	}
	var blogText string
	for check_post_in_tbl.Next() {
		check_post_in_tbl.Scan(&blogText)
		tmpl.ExecuteTemplate(w_page, "blog", blogText)
	}
	defer check_post_in_tbl.Close()
}

func index_logged_out(w_page http.ResponseWriter, r *http.Request) {

	curUser = ""
	index(w_page, r)
}

func (usr TUser) newUser(db *sql.DB) {

	insert, err := db.Query(fmt.Sprintf("INSERT INTO tbl_user (email, password) VALUES ('%s', '%s')", usr.Email, usr.Password))
	if err != nil {
		panic(err)
	}
	defer insert.Close()
}

func newPost_page(w_page http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html",
		"templates/header.html",
		"templates/footer.html",
		"templates/post.html")
	if err != nil {
		panic(err)
	}
	tmpl.ExecuteTemplate(w_page, "post", nil)
}

func newPost(w_page http.ResponseWriter, r *http.Request) {
	var _post TPost
	_post.Email = curUser
	_post.Title = r.FormValue("inputTitle")
	_post.Content = r.FormValue("inputContent")
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbname))
	if err != nil {
		panic(err.Error())
	}
	var query = fmt.Sprintf("INSERT INTO tbl_posts (email, title, post) VALUES ('%s', \"%s\", \"%s\")", _post.Email, _post.Title, _post.Content)
	insert, err := db.Query(query)
	if err != nil {
		panic(err)
	}
	defer insert.Close()
	tmpl, err := template.ParseFiles("templates/index.html",
		"templates/header.html",
		"templates/footer.html",
		"templates/signup.html",
		"templates/post.html",
		"templates/blog.html",
		"templates/signup_success.html")
	if err != nil {
		panic(err)
	}
	tmpl.ExecuteTemplate(w_page, "index", curUser)

	query = fmt.Sprintf("SELECT post FROM tbl_posts WHERE email = '" + curUser + "'")
	check_post_in_tbl, err := db.Query(query)
	if err != nil {
		panic(err)
	}
	var blogText string
	for check_post_in_tbl.Next() {
		check_post_in_tbl.Scan(&blogText)
		tmpl.ExecuteTemplate(w_page, "blog", blogText)
	}
	defer check_post_in_tbl.Close()
}

func newUser_page(w_page http.ResponseWriter, r *http.Request) {
	var _user TUser
	_user.Email = r.FormValue("inputEmail")
	_user.Password = r.FormValue("inputPassword")

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbname))
	if err != nil {
		panic(err.Error())
	}
	var query = fmt.Sprintf("SELECT email, password FROM tbl_user WHERE email = '" + _user.Email + "'")
	check_user_in_tbl, err := db.Query(query)
	if err != nil {
		panic(err)
	}
	var check_email string
	var check_password string
	for check_user_in_tbl.Next() {
		check_user_in_tbl.Scan(&check_email, &check_password)
	}
	defer check_user_in_tbl.Close()

	if check_email == "" {
		curUser = _user.Email
		_user.newUser(db)
	} else {

		if check_password != _user.Password {
			_user.Email = "ERROR_WRONG_PASSWORD"
		} else {
			curUser = _user.Email
		}
	}

	tmpl, err := template.ParseFiles("templates/index.html",
		"templates/header.html",
		"templates/footer.html",
		"templates/signup.html",
		"templates/post.html",
		"templates/blog.html",
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
