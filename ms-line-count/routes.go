package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

//CountLinesReq is a body for requesting counting of lines.
type CountLinesReq struct {
	FName string `json:"fname"`
}

//fetchBrowserCounts
func handleCountLines(w http.ResponseWriter, r *http.Request) {
	log.Println("Handling Count lines")
	//object for body
	var CLR CountLinesReq

	//decode data in body
	err := json.NewDecoder(r.Body).Decode(&CLR)
	if err != nil {
		e := NewError(http.StatusBadRequest, err.Error())
		http.Error(w, e.json, http.StatusBadRequest)
		return
	}

	//fetch data for file name
	lf, err := LogStore.fetchData(CLR.FName)
	if err != nil {
		e := NewError(http.StatusBadRequest, err.Error())
		http.Error(w, e.json, http.StatusBadRequest)
		return
	}
	//if no lines then file not found, throw error
	if len(lf.Logs) < 1 {
		e := NewError(http.StatusNotFound, CLR.FName+" not found!")
		http.Error(w, e.json, http.StatusNotFound)
		return
	}
	log.Println("Number of lines:", len(lf.Logs))

	//store num count
	err = LogStore.StoreCountLines(CLR.FName, len(lf.Logs))
	if err != nil {
		log.Println(err)
		e := NewError(http.StatusBadRequest, err.Error())
		http.Error(w, e.json, http.StatusBadRequest)
		return
	}

	//create response and get json
	res := Response{201, "Success"}
	jOut, _ := res.JSON()

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, jOut)
}

//handleLineCount
func handleLineCount(w http.ResponseWriter, r *http.Request) {

	//fetch parameters from url
	params := mux.Vars(r)

	fname := params["fname"]

	lc, err := LogStore.fetchLineCount(fname)
	if err != nil {
		e := NewError(http.StatusBadRequest, err.Error())
		http.Error(w, e.json, e.StatusCode)
		return
	}

	//No results returned. Return not found
	if len(lc) < 1 {
		e := NewError(http.StatusNotFound, "File "+fname+" not found!")
		http.Error(w, e.json, http.StatusBadRequest)
		return
	}

	//Create map for response of results
	lineCountMap := make(map[string]int)
	lineCountMap[fname] = lc[0].Count

	//create response and get json
	res := ResponseInt{201, "Success", lineCountMap}
	jOut, _ := res.JSON()

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(jOut))
}
