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
	ServiceName = "ms-line-count"

	//Start DB connection
	LogStore = Database{}

	//read config from file using viper.
	ReadConfig()

	var err error

	//Open/Create the DB file for data storage
	LogStore.db, err = sql.Open("sqlite3", "../ETL.db")
	if err != nil {
		log.Fatal(err)
	}
	defer LogStore.db.Close()

	//init database
	LogStore.dbInit()

	//Define routes.
	r := mux.NewRouter()
	r.HandleFunc("/lines/count", handleCountLines).Methods("POST")
	r.HandleFunc("/lines/count/{fname}", handleLineCount).Methods("GET")
	log.Println("Listening on: ", viper.GetString("services."+ServiceName))
	log.Fatal(http.ListenAndServe(":"+viper.GetString("services."+ServiceName), r))
}
