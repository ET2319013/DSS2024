go get -u golang.org/x/crypto/bcrypt
go build -o main.exe sql_functions.go main.go handlers.go
start main.exe
