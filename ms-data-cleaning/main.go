package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/viper"
)

//LogStore is a database connection
var LogStore Database

//ServiceName is the name of this service. It will match the service in the config file.
var ServiceName string

//ReadConfig reads the config from a file DO NOT MODIFY
func ReadConfig() {
	// Set the file name of the configurations file
	viper.SetConfigName("config")
	// Set the path to look for the configurations file
	viper.AddConfigPath("../")
	//Set the config type
	viper.SetConfigType("yml")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s", err)
	}
}

func main() {

	//DEFINE THE SERVICE NAME. CHANGE THIS TO MATCH A LINE IN CONFIG FILE
	ServiceName = "ms-data-cleaning"

	//Start DB connection
	LogStore = Database{}

	//read config from file using viper.
	ReadConfig()

	var err error

	//Open/Create the DB file for data storage. This is shared across all microservices. DO NOT CHANGE
	LogStore.db, err = sql.Open("sqlite3", "../ETL.db")
	if err != nil {
		log.Fatal(err)
	}
	defer LogStore.db.Close()

	//init database
	LogStore.dbInit()

	//Define routes.
	r := mux.NewRouter()

	//Define your routes here. You may need to add more routes here.
	r.HandleFunc("/timsfunc", handleRoute).Methods("POST")
	r.HandleFunc("/route/with/{PARAM_NAME}", handleRouteParameter).Methods("GET")

	//Serve the webserver. You should not change this
	log.Println("Listening on: ", viper.GetString("services."+ServiceName))
	log.Fatal(http.ListenAndServe(":"+viper.GetString("services."+ServiceName), r))
}
