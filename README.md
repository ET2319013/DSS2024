# go-web-server
[local](http://localhost:8080/)

use MySQL server creds:

var dbUser string = "root"
var dbPass string = "admin"
var dbHost string = "localhost"
var dbPort string = "3306"

 go build -o main.exe sql_functions.go main.go handlers.go

// I'd like to use this library https://github.com/golang-migrate/migrate/blob/master/database/postgres/TUTORIAL.md, 
// but failed. so I just use db.Exec(string(data)) from files used to be migration files
