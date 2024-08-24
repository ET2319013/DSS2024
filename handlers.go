package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	_ "reflect"

	"golang.org/x/crypto/bcrypt"

	_ "github.com/go-sql-driver/mysql"
	// _ "github.com/golang-migrate/migrate/v4"
	// I'd like to use this library https://github.com/golang-migrate/migrate/blob/master/database/postgres/TUTORIAL.md,
	// but failed
)

func index(wPage http.ResponseWriter, r *http.Request) {
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
	tmpl.ExecuteTemplate(wPage, "index", curUser)

	db := OpenSQL(dbUser, dbPass, dbHost, dbPort, dbname)

	var query = fmt.Sprintf("SELECT post FROM tbl_posts WHERE email = '" + curUser + "'")
	checkPostInTable, err := db.Query(query)
	if err != nil {
		panic(err)
	}
	var blogText string
	for checkPostInTable.Next() {
		checkPostInTable.Scan(&blogText)
		tmpl.ExecuteTemplate(wPage, "blog", blogText)
	}
	defer checkPostInTable.Close()
}

func indexLoggedOut(wPage http.ResponseWriter, r *http.Request) {

	curUser = ""
	index(wPage, r)
}

func (usr TUser) newUser(db *sql.DB) {

	insert, err := db.Query(fmt.Sprintf("INSERT INTO tbl_user (email, password) VALUES ('%s', '%s')", usr.Email, usr.Password))
	if err != nil {
		panic(err)
	}
	defer insert.Close()
}

func newPostPage(wPage http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html",
		"templates/header.html",
		"templates/footer.html",
		"templates/post.html")
	if err != nil {
		panic(err)
	}
	tmpl.ExecuteTemplate(wPage, "post", nil)
}

func newPost(wPage http.ResponseWriter, r *http.Request) {
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
	tmpl.ExecuteTemplate(wPage, "index", curUser)

	query = fmt.Sprintf("SELECT post FROM tbl_posts WHERE email = '" + curUser + "'")
	checkPostInTable, err := db.Query(query)
	if err != nil {
		panic(err)
	}
	var blogText string
	for checkPostInTable.Next() {
		checkPostInTable.Scan(&blogText)
		tmpl.ExecuteTemplate(wPage, "blog", blogText)
	}
	defer checkPostInTable.Close()
}

func newUserPage(wPage http.ResponseWriter, r *http.Request) {
	var _user TUser
	_user.Email = r.FormValue("inputEmail")
	password := []byte(r.FormValue("inputPassword"))
	// Hashing the password with the default cost of 10
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	_user.Password = string(hashedPassword)

	db := OpenSQL(dbUser, dbPass, dbHost, dbPort, dbname)
	var query = fmt.Sprintf("SELECT email, password FROM tbl_user WHERE email = '" + _user.Email + "'")
	checkUserInTable, err := db.Query(query)
	if err != nil {
		panic(err)
	}
	var checkEmail string
	var checkPassword string
	for checkUserInTable.Next() {
		checkUserInTable.Scan(&checkEmail, &checkPassword)
	}
	defer checkUserInTable.Close()

	if checkEmail == "" {
		curUser = _user.Email
		_user.newUser(db)
	} else {
		err = bcrypt.CompareHashAndPassword(hashedPassword, password)

		if err != nil {
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
	tmpl.ExecuteTemplate(wPage, "signup_success", _user)
}

func signUpPage(wPage http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html",
		"templates/header.html",
		"templates/footer.html",
		"templates/signup.html")
	if err != nil {
		panic(err)
	}
	tmpl.ExecuteTemplate(wPage, "signup", nil)
}
