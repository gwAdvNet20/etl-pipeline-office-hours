package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/viper"
)

//runPipeline calls functions of other microservices to calculate data metrics
func runPipeline(fname string) map[string]bool {
	//store results of the pipeline
	results := make(map[string]bool)

	//parseFile
	//make call to datacleaner microservice here
	results["Data Cleaner"] = dataCleaner(fname)

	//count lines
	results["Count Lines"] = countLines(fname)

	//count browsers
	//make call to browserCounts microservice here

	//count visitors
	//make call to visitorCounts microservice here

	//count websites
	//make call to websiteCounter microservice here

	return results
}

//ReqStruct represents the body of a request. Add json fields below as needed
type ReqStruct struct {
	FName   string `json:"fname"`
	ANumber int    `json:"anumber"`
}

func dataCleaner(fname string) bool {

	requestBody, err := json.Marshal(ReqStruct{
		FName:   fname,
		ANumber: 123,
	})
	if err != nil {
		log.Println("Error parsing countLines request body:", requestBody)
		return false
	}

	url := "http://localhost:" + viper.GetString("services.ms-data-cleaning") + "/timsfunc"

	log.Println("Posting URL: ", url, " with ", string(requestBody))
	//make request to ms
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		log.Println("Error making post request to ", url, ": ", err)
		return false
	}
	defer resp.Body.Close()

	//decode reesponse body
	var result map[string]interface{}
	// var result Response
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		log.Println("Error decoding json:", err)
		return false
	}
	//if statusCode is above 300 then its an error, parse and return
	if result["statusCode"].(float64) > 300 {
		log.Println("Error", result["error"].(string))
		return false
	}
	//otherwise succesful return true
	return true
}

func countLines(fname string) bool {

	requestBody, err := json.Marshal(map[string]string{
		"fname": fname,
	})
	if err != nil {
		log.Println("Error parsing countLines request body:", requestBody)
		return false
	}

	url := "http://localhost:" + viper.GetString("services.ms-line-count") + "/lines/count"

	log.Println("Posting URL: ", url, " with ", string(requestBody))
	//make request to ms
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		log.Println("Error making post request to ", url, ": ", err)
		return false
	}
	defer resp.Body.Close()

	//decode reesponse body
	var result map[string]interface{}
	// var result Response
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		log.Println("Error decoding json:", err)
		return false
	}
	//if statusCode is above 300 then its an error, parse and return
	if result["statusCode"].(float64) > 300 {
		log.Println("Error", result["error"].(string))
		return false
	}
	//otherwise succesful return true
	return true
}

//THIS FUNCTIONALITY NEEDS TO BE MOVED TO ms-data-cleaning
//TO TEST YOUR IMPLEMENTATION, COMMENT THIS OUT AND SEE IF LINE COUNT STILL WORKS.
//processLogFile takes in an uploaded logfile, stores the data, processes stats.
func processLogFile(rawLogFile []byte, fname string) bool {
	var lines []string
	scanner := bufio.NewScanner(strings.NewReader(string(rawLogFile)))

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	logFile := parseFile(lines, fname)

	//Store parsed logs
	LogStore.StoreLogLine(logFile)

	return true
}

//THIS FUNCTIONALITY NEEDS TO BE MOVED TO ms-data-cleaning
//TO TEST YOUR IMPLEMENTATION, COMMENT THIS OUT AND SEE IF LINE COUNT STILL WORKS.
//parseFile will take a slice of strings and parse the fields.
func parseFile(lines []string, fname string) LogFile {

	//list to store log lines
	lf := LogFile{}

	for _, line := range lines {

		lineSplit := strings.Split(line, " ")
		userAgent := strings.Join(lineSplit[11:], " ")
		status, _ := strconv.Atoi(lineSplit[8])
		totalBytes, _ := strconv.Atoi(lineSplit[9])
		tempLine := LogLine{
			fname,
			line,
			lineSplit[0],
			lineSplit[3] + " " + lineSplit[4],
			lineSplit[5],
			lineSplit[6],
			status,
			totalBytes,
			lineSplit[10],
			userAgent,
			time.Now(),
		}

		lf.Logs = append(lf.Logs, tempLine)
	}
	return lf
}
