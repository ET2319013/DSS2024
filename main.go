package main

import (
	"fmt"
	"net/http"
	_ "reflect"

	_ "github.com/go-sql-driver/mysql"
	// _ "github.com/golang-migrate/migrate/v4"
	// I'd like to use this library https://github.com/golang-migrate/migrate/blob/master/database/postgres/TUTORIAL.md,
	// but failed
)

// MYSQL CREDS
var dbUser string = "root"
var dbPass string = "admin"
var dbHost string = "localhost"
var dbPort string = "3306"
var dbname string = "golangdb"

var curUser string = ""

//var curID int = -1

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

	db := InitDB(dbUser, dbPass, dbHost, dbPort, dbname) //Connect driver, Open or Create DB

	stmt := ExecuteFile(db, "0001_create_users.up.sql")
	fmt.Println(stmt)

	stmt1 := ExecuteFile(db, "0001_create_posts.up.sql")
	fmt.Println(stmt1)

	handleRequest()
}

func handleRequest() {
	http.HandleFunc("/", index)
	http.HandleFunc("/signup/", signUpPage)
	http.HandleFunc("/logout/", indexLoggedOut)
	http.HandleFunc("/create_user/", newUserPage)
	http.HandleFunc("/post/", newPostPage)
	http.HandleFunc("/create_post/", newPost)

	http.ListenAndServe(":8080", nil)
}
