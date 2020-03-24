package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

//handleLinesCount handles fetching line counts
func handleLinesCount(w http.ResponseWriter, r *http.Request) {

	//fetch parameters from url
	params := mux.Vars(r)

	fname := params["fname"]
	url := "http://localhost:" + viper.GetString("services.ms-line-count") + "/lines/count/" + fname

	log.Println("Fetching URL: ", url)
	//make request to ms
	resp, err := http.Get(url)
	if err != nil {
		e := NewError(http.StatusBadRequest, err.Error())
		http.Error(w, e.json, http.StatusBadRequest)
		return
	}
	defer resp.Body.Close()

	//decode reesponse body
	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	//if statusCode is above 300 then its an error, parse and return
	if result["statusCode"].(float64) > 300 {
		log.Println("Error", result["error"].(string))
		e := NewError(http.StatusInternalServerError, result["error"].(string))
		http.Error(w, e.json, http.StatusBadRequest)
		return
	}

	//convert response to proper response
	resOut := ResponseInt{int(result["statusCode"].(float64)), result["message"].(string), ConvertMapInterfaceToMapInt(result["data"])}
	jOut, _ := resOut.JSON()

	//return response to client
	w.WriteHeader(int(result["statusCode"].(float64)))
	fmt.Fprintf(w, jOut)
}

//handleServeUploadPage serves the static html file
func handleServeUploadPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./html/upload.html")
}

//handleUploadLog handles uploading of log file and triggers etl pipeline
func handleUploadLog(w http.ResponseWriter, r *http.Request) {
	// Parse our multipart form, 10 << 20 specifies a maximum
	// upload of 10 MB files.
	r.ParseMultipartForm(10 << 20)

	// FormFile returns the first file for the given key `myFile`
	// it also returns the FileHeader so we can get the Filename,
	// the Header and the size of the file
	file, handler, err := r.FormFile("myFile")
	if err != nil {
		log.Println("Error Retrieving the File")
		log.Println(err)
		return
	}

	defer file.Close()
	fmt.Printf("Received Uploaded File: %+v\n", handler.Filename)

	// read all of the contents of our uploaded file into a
	// byte array
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Println(err)
	}

	//clean log file and store in db
	processLogFile(fileBytes, handler.Filename)
	result := runPipeline(handler.Filename)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "<h1>Pipeline Status</h1>")
	fmt.Fprintf(w, "Log File Uploaded Successfully<br>")
	//iterate over results
	for k, v := range result {
		var status string
		//check status
		if v {
			status = `<font size="3" color="green">Completed</font>`
		} else {
			status = `<font size="3" color="red">Failed</font>`
		}
		//print status
		fmt.Fprintf(w, "<strong>"+k+"</strong>: "+status+"<br>")
	}

}
