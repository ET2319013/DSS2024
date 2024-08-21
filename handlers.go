package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	_ "reflect"

	_ "github.com/go-sql-driver/mysql"
	// _ "github.com/golang-migrate/migrate/v4"
	// I'd like to use this library https://github.com/golang-migrate/migrate/blob/master/database/postgres/TUTORIAL.md,
	// but failed
)

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

	db := OpenSQL(dbUser, dbPass, dbHost, dbPort, dbname)

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
	db := OpenSQL(dbUser, dbPass, dbHost, dbPort, dbname)
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

	db := OpenSQL(dbUser, dbPass, dbHost, dbPort, dbname)
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