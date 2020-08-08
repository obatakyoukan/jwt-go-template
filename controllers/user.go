package controllers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"../models"
	"../utils"

	"golang.org/x/crypto/bcrypt"
)

type Controller struct{}

func (c Controller) Signup(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Decoding the POST
		var user models.User
		var error models.Error
		json.NewDecoder(r.Body).Decode(&user)

		// validations
		if user.Email == "" {
			// respond with error
			error.Message = "Email is missing."
			utils.ResponseWithError(w, http.StatusBadRequest, error)
			return
		}
		if user.Password == "" {
			// respond with error
			error.Message = "Password is missing."
			utils.ResponseWithError(w, http.StatusBadRequest, error)
			return
		}

		// Hashing the Password
		hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
		if err != nil {
			log.Fatal(err)
		}
		user.Password = string(hash)

		// Insert user's information on DB
		stmt := "INSERT INTO users (email, password) VALUES($1, $2) RETURNING id;"
		err = db.QueryRow(stmt, user.Email, user.Password).Scan(&user.ID)
		if err != nil {
			error.Message = "Server error."
			utils.ResponseWithError(w, http.StatusInternalServerError, error)
			return
		}

		// Clearing the user Password
		user.Password = ""

		// Respond
		w.Header().Set("Content-Type", "application/json")
		utils.ResponseJSON(w, user)
	}
}

func (c Controller) Login(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var user models.User
		var jwt models.JWT
		var error models.Error

		json.NewDecoder(r.Body).Decode(&user)

		if user.Email == "" {
			error.Message = "Email is missing."
			utils.ResponseWithError(w, http.StatusBadRequest, error)
			return
		}
		if user.Password == "" {
			error.Message = "Password is missing."
			utils.ResponseWithError(w, http.StatusBadRequest, error)
			return
		}

		password := user.Password

		// Extract user's informations from DB
		row := db.QueryRow("SELECT * FROM users WHERE email = $1", user.Email)
		err := row.Scan(&user.ID, &user.Email, &user.Password)

		if err != nil {
			if err == sql.ErrNoRows {
				error.Message = "The user does not exist."
				utils.ResponseWithError(w, http.StatusBadRequest, error)
				return
			} else {
				log.Fatal(err)
			}
		}

		hashedPassword := user.Password
		isValidPassword := utils.ComparePasswords(hashedPassword, []byte(password))

		if isValidPassword {

			token, err := utils.GenerateToken(user)
			if err != nil {
				log.Fatal(err)
			}

			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Authorization", token)

			jwt.Token = token
			utils.ResponseJSON(w, jwt)
		} else {
			error.Message = "Invalid Password."
			utils.ResponseWithError(w, http.StatusUnauthorized, error)
		}
	}
}
