package main

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"fmt"
	"log"
	"math/rand"
	"os"
	"regexp"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Row struct {
	Email     string
	Password  string
	Created   string
	LastLogin string
}

var (
	user           string = ""
	password       string = ""
	connectionType string = ""
	hostname       string = ""
	port           string = ""
	dbName         string = ""
	tableName      string = ""
)

func main() {
	fmt.Println(os.Args)
	if len(os.Args) != 3 {
		log.Fatal("Usage: ./mod-db <add | delete> <email>")
	} else if strings.ToLower(os.Args[1]) != "add" && strings.ToLower(os.Args[1]) != "delete" {
		log.Fatal("only add or delete command allowed")
	} else if emailRegex := regexp.MustCompile("^[a-z0-9._%+]+@[A-Za-z]+.[a-z]{2,4}$"); !emailRegex.MatchString(os.Args[2]) {
		log.Fatal("invalid email")
	}
	connString := fmt.Sprintf("%s:%s@%s(%s:%s)/%s", user, password, connectionType, hostname, port, dbName)
	conn, err := sql.Open("mysql", connString)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	if strings.ToLower(os.Args[1]) == "add" {
		addUser(os.Args[2], conn)
	}
}

func generatePassword() string {
	rand.Seed(time.Now().UnixNano())
	var passwordLength = 16
	var charset = []byte(`abcdefghijklmnopqrstuvwxyz!@#$%^&*()-_<>,./?"'{}~\|+=ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789`)
	pwd := make([]byte, passwordLength)
	for i := range pwd {
		// randomly select 1 character from given charset
		pwd[i] = charset[rand.Intn(len(charset))]
	}
	return string(pwd)
}

func addUser(email string, conn *sql.DB) {
	var version string
	err := conn.QueryRow("SELECT VERSION()").Scan(&version)
	if err != nil {
		log.Fatal(err)
	} else {
		timestamp := time.Now()
		newUser := Row{
			Email:     email,
			Password:  generatePassword(),
			Created:   timestamp.Format("2006-01-02 15:04:05"),
			LastLogin: timestamp.Format("2006-01-02 15:04:05"),
		}

		query := fmt.Sprintf("INSERT INTO %s (email, password, created, last_login) VALUES(?, ?, ?, ?)", tableName)
		insert, err := conn.Prepare(query)
		if err != nil {
			log.Fatal(err)
		}

		passwordHash := md5.Sum([]byte(newUser.Password))
		resp, err := insert.Exec(newUser.Email, hex.EncodeToString(passwordHash[:]), newUser.Created, newUser.LastLogin)
		if err != nil {
			log.Fatal(err)
		}
		insert.Close()

		affectedRows, _ := resp.RowsAffected()
		fmt.Printf("Rows affected: %d\nUsername: %s\nPassword: %s\n(Save this password as it will not be saved anywhere)", affectedRows, newUser.Email, newUser.Password)
	}
}
