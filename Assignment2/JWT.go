package main

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type User struct {
	Name string
	Id   int
	jwt.RegisteredClaims
}

func main() {

	u := User{
		Name: "Ishaan",
		Id:   1,
	}

	tok, _ := CreateJwtToken(u)

	fmt.Println("Token is: ", tok)

	ParsedUser, _ := ParseToken(tok)

	fmt.Println("Parsed User: \n", ParsedUser)

}

func CreateJwtToken(user User) (string, error) {

	claim := &User{
		Name: user.Name,
		Id:   user.Id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	return token.SignedString([]byte("JGNU3T48FNRVNRG5R8GERV"))

}

func ParseToken(token string) (*User, error) {
	//claim := &User{}
	var claim User

	tok, err := jwt.ParseWithClaims(token, &claim, func(t *jwt.Token) (interface{}, error) {
		return []byte("JGNU3T48FNRVNRG5R8GERV"), nil
	})

	if err != nil {
		return nil, err
	}

	if !tok.Valid {
		return nil, err
	}

	return &claim, nil

}
