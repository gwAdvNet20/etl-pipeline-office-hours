package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

//BodyStruct represents the body of a request. Add json fields below as needed
type BodyStruct struct {
	FName string `json:"fname"`
}

//YOU WILL NEED TO CHANGE THIS TO THE REAL FUNCTIONALITY
//handleRoute handles the root route
func handleRoute(w http.ResponseWriter, r *http.Request) {

	//LOGIC for this endpoint goes here
	//object for body
	var BS BodyStruct

	//decode data in body
	err := json.NewDecoder(r.Body).Decode(&BS)
	if err != nil {
		//return error if failure
		log.Println("Error parsing json:", err)
		e := NewError(http.StatusBadRequest, err.Error())
		http.Error(w, e.json, http.StatusBadRequest)
		return
	}

	log.Println("BODY IS ", BS)

	//store 1 for value of filename
	LogStore.StoreValue(BS.FName, 1)

	//create a response using shared response lib and get json
	res := Response{200, "Success"}
	jOut, _ := res.JSON()

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, jOut)
}

//YOU WILL NEED TO CHANGE THIS TO THE REAL FUNCTIONALITY
//handleRouteParameter
func handleRouteParameter(w http.ResponseWriter, r *http.Request) {

	//fetch parameters from url
	params := mux.Vars(r)

	//fetch parameter matchin value in {}
	fname := params["PARAM_NAME"]
	log.Println("Fname is ", fname)

	//retrieve value for filename
	results, err := LogStore.fetchValues(fname)
	if err != nil {
		//return error if failure
		log.Println("Error parsing json:", err)
		e := NewError(http.StatusBadRequest, err.Error())
		http.Error(w, e.json, http.StatusBadRequest)
		return
	}
	resultMap := map[string]int{"result": len(results)}
	//create response and get json
	res := ResponseInt{200, "Success", resultMap}
	jOut, _ := res.JSON()

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(jOut))
}
