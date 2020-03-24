package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/viper"
)

//LogStore is a database connection
var LogStore Database

//ServiceName is the name of this service. It will match the service in the config file.
var ServiceName string

//LogLine represents fields in a given log line
type LogLine struct {
	Name          string
	RawLog        string
	RemoteAddr    string
	TimeLocal     string
	RequestType   string
	RequestPath   string
	Status        int
	BodyBytesSent int
	HTTPReferer   string
	HTTPUserAgent string
	Created       time.Time
}

//LogFile represents a logfile with multiple lines
type LogFile struct {
	Logs []LogLine
}

//ReadConfig reads the config from a file
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

	//Start DB connection
	LogStore = Database{}
	ServiceName = "monolith"

	//read config from file using viper.
	ReadConfig()

	var err error

	//Open/Create the DB file for data storage
	LogStore.db, err = sql.Open("sqlite3", "ETL.db")
	if err != nil {
		log.Fatal(err)
	}
	defer LogStore.db.Close()

	//Create table if not found
	LogStore.dbInit()

	//Define routes.
	r := mux.NewRouter()
	r.HandleFunc("/browser/count", BasicAuth(handleBrowserCount, "read", "Please enter your username and password for this site")).Methods("GET")
	r.HandleFunc("/visitor/count", BasicAuth(handleVisitorCount, "read", "Please enter your username and password for this site")).Methods("GET")
	r.HandleFunc("/", BasicAuth(handleServeUploadPage, "write", "Please enter your username and password for this site"))
	r.HandleFunc("/upload/log", BasicAuth(handleUploadLog, "write", "Please enter your username and password for this site"))
	log.Println("Listening on: ", viper.GetString("services."+ServiceName))
	log.Fatal(http.ListenAndServe(":"+viper.GetString("services."+ServiceName), r))
}

//check if element found in golang string
func checkExists(list []string, v string) bool {
	for _, a := range list {
		if a == v {
			return true
		}
	}
	return false
}
