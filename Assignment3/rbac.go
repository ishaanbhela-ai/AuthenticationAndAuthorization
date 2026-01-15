package main

import (
	"fmt"
)

type Document struct {
	id      int
	owner   string
	content string
}

type User struct {
	name    string
	isAdmin bool
}

var users []User
var docs []Document

func main() {
	userA := User{
		name:    "A",
		isAdmin: true,
	}

	userB := User{
		name:    "B",
		isAdmin: false,
	}

	userC := User{
		name:    "C",
		isAdmin: false,
	}

	users = append(users, userA, userB, userC)

	doc1 := Document{
		id:      1,
		owner:   "A",
		content: "Admin document owned by A",
	}

	doc2 := Document{
		id:      2,
		owner:   "B",
		content: "User B's private document",
	}

	doc3 := Document{
		id:      3,
		owner:   "C",
		content: "User C's notes",
	}

	doc4 := Document{
		id:      4,
		owner:   "A",
		content: "Another admin document",
	}

	docs = append(docs, doc1, doc2, doc3, doc4)

	var u string
	var id int
	var currentUser User
	var DocToAccess Document

	fmt.Println("Enter Your Name")
	fmt.Scan(&u)

	fmt.Println("Enter Document id to access")
	fmt.Scan(&id)

	for _, value := range users {
		if value.name == u {
			currentUser = value
		}
	}

	for _, value := range docs {
		if value.id == id {
			DocToAccess = value
		}
	}

	if currentUser.isAdmin {
		fmt.Println("Document Content: \n", DocToAccess.content)
		return
	}

	if currentUser.name == DocToAccess.owner {
		fmt.Println("Document Content: \n", DocToAccess.content)
		return
	}

	fmt.Println("Document Inaccessible, You are not the owner or Admin")

}
