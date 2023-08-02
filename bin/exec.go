package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/quevivasbien/bird-game/db"
)

const AWS_REGION = "us-east-1"

func Help() {
	space := strings.Repeat(" ", 3)
	fmt.Println("Commands")
	fmt.Println("- resetdb" + space + "(re-initialize database)")
	fmt.Println("- listusers" + space + "(list all users currently in database)")
	fmt.Println("- deluser [name]" + space + "(delete user)")
	fmt.Println("- makeadmin [name] [password]" + space + "(create admin account)")
}

func ResetDB() {
	tables, err := db.GetTables(AWS_REGION)
	if err != nil {
		panic(fmt.Sprint("Problem getting existing tables:", err))
	}
	err = tables.Reset()
	if err != nil {
		panic(fmt.Sprint("Problem while resetting tables:", err))
	}
	fmt.Println("Successfully reset database")
}

func ListUsers() {
	tables, err := db.GetTables(AWS_REGION)
	if err != nil {
		panic(fmt.Sprint("Problem getting existing tables:", err))
	}
	users, err := tables.AllUsers()
	if err != nil {
		panic(fmt.Sprint("Problem getting existing tables:", err))
	}
	fmt.Printf("Users (%d):\n", len(users))
	for i, user := range users {
		fmt.Printf("%d: %v\n", i+1, user)
	}
}

func MakeAdmin(name string, password string) {
	tables, err := db.GetTables(AWS_REGION)
	if err != nil {
		panic(fmt.Sprint("Problem getting existing tables:", err))
	}
	err = tables.PutUser(db.User{Name: name, Password: password, Admin: true})
	if err != nil {
		panic(fmt.Sprint("Problem creating admin user on database:", err))
	}
	fmt.Println("Successfully created admin user")
}

func DelUser(name string) {
	tables, err := db.GetTables(AWS_REGION)
	if err != nil {
		panic(fmt.Sprint("Problem getting existing tables:", err))
	}
	err = tables.DeleteUser(name)
	if err != nil {
		panic(fmt.Sprint("Problem deleting user:", err))
	}
	fmt.Println("Successfully deleted user")
}

func main() {
	if len(os.Args) < 2 {
		Help()
		return
	}

	command := os.Args[1]
	if command == "resetdb" {
		ResetDB()
	} else if command == "listusers" {
		ListUsers()
	} else if command == "deluser" {
		if len(os.Args) < 3 {
			fmt.Println("Missing username to delete")
		} else {
			DelUser(os.Args[2])
		}
	} else if command == "makeadmin" {
		if len(os.Args) < 4 {
			fmt.Println("Missing name and/or password for admin user")
		} else {
			MakeAdmin(os.Args[2], os.Args[3])
		}
	} else {
		Help()
	}
}
