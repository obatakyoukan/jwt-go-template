package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/subosito/gotenv"

	"./controllers"
	"./driver"

	"github.com/gorilla/mux"
)

var db *sql.DB

func init() {
	gotenv.Load()
}

func main() {

	db = driver.ConnectDB()
	controller := controllers.Controller{}

	router := mux.NewRouter()

	router.HandleFunc("/signup", controller.Signup(db)).Methods("POST")
	router.HandleFunc("/signup", controller.Login(db)).Methods("POST")
	router.HandleFunc("/protected", controller.TokenVerifyMiddleWare(controller.ProtectedEndpoint())).Methods("GET")

	log.Println("Listen on port 8000...")
	log.Fatal(http.ListenAndServe(":8000", router))
}
