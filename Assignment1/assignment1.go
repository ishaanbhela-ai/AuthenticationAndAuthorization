package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type user struct {
	Name     string
	Password []byte
	Phno     string
}

var users []user

var OTP int

func main() {

	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/login", handleLogin)
	mux.HandleFunc("POST /api/signup", handleSignUp)
	mux.HandleFunc("POST /api/twoFactorAuth", handleTwoFactorAuth)
	mux.HandleFunc("POST /api/twoFactorAuthVerify", handleTwoFactorAuthVerification)

	fmt.Println("Listening on port :8080")
	http.ListenAndServe(":8080", enableCORS(mux))

}

func handleSignUp(w http.ResponseWriter, r *http.Request) {

	type registerUser struct {
		Name        string `json:"Name"`
		Password    string `json:"Password"`
		ConfirmPass string `json:"ConfirmPass"`
		Phno        string `json:"Phno"`
	}

	var newUser registerUser

	json.NewDecoder(r.Body).Decode(&newUser)

	fmt.Println(newUser.Password, " - ", newUser.ConfirmPass)

	if newUser.Password != newUser.ConfirmPass {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"Message": "Password and confirm Password must be same",
		})
		return
	}

	var u user

	HashedPass, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"Message": err.Error(),
		})
		return
	}

	u.Name = newUser.Name
	u.Password = HashedPass
	u.Phno = newUser.Phno

	users = append(users, u)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"Message": "New User added",
	})

}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	type Login struct {
		Name     string `json:"Name"`
		Password string `json:"Password"`
	}

	var login Login
	var u user

	json.NewDecoder(r.Body).Decode(&login)

	for _, value := range users {
		if value.Name == login.Name {
			u = value
			break
		}
	}

	err := bcrypt.CompareHashAndPassword(u.Password, []byte(login.Password))

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"Message": err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(u)

	w.WriteHeader(http.StatusOK)
}

func handleTwoFactorAuth(w http.ResponseWriter, r *http.Request) {

	Otp, _ := rand.Int(rand.Reader, big.NewInt(10000))

	message := fmt.Sprintf("OTP Sent to mobile number %d", Otp.Int64())

	fmt.Println(message)

	OTP = int(Otp.Int64())

	json.NewEncoder(w).Encode(map[string]string{
		"message": message,
	})
}

func handleTwoFactorAuthVerification(w http.ResponseWriter, r *http.Request) {
	type verify struct {
		Otp int `json:"Otp"`
	}

	var req verify

	json.NewDecoder(r.Body).Decode(&req)

	if req.Otp != OTP {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"Message": "OTP Does not match",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"Message": "2 Factor Authentication Successful",
	})
}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "*")

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight request
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
