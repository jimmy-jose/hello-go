package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// BaseResponse is a struct with a basic response format
// status is the status code of the response and message contains message
type BaseResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

// Currency contains data about each currency type
type Currency struct {
	ID                  int    `json:"id"`
	CurrencyCode        string `json:"currency_code"`
	CurrencyDescription string `json:"currency_description"`
	IsoCurrencyCode     string `json:"iso_currency_code"`
	Country             string `json:"country"`
}

// CurrencyPayload contains a list of currencies total count and a status message
type CurrencyPayload struct {
	Currencies    []Currency `json:"currencies"`
	TotalCount    []Currency `json:"total_count"`
	StatusMessage []Currency `json:"status_message"`
}

// CurrencyData is the data returned from currency enum api
type CurrencyData struct {
	Payload CurrencyPayload `json:"payload"`
	Status  int             `json:"status"`
}

func main() {
	handleRequests()
}

// handleRequests handles all the reqests coming in
func handleRequests() {
	router := mux.NewRouter()

	router.HandleFunc("/", home).Methods("GET")

	router.HandleFunc("/hello", returnGreeting).Methods("GET")

	router.HandleFunc("/getCurrencies", getCurrencies).Methods("GET")

	log.Fatal(http.ListenAndServe(":5000", router))
}

// returnGreeting writes a simple BaseResponse
func returnGreeting(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnGreeting")
	baseResponse := BaseResponse{Status: 200, Message: "Hello world"}
	w.Header().Set("Content-Type", "application/json") // this
	json.NewEncoder(w).Encode(baseResponse)
}

// home write a simple message to the response writer
func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the Home!")
	fmt.Println("Endpoint Hit: home")
}

// getCurrencies fetch the curriencies from spenmo server and writes it to the response writer
func getCurrencies(w http.ResponseWriter, r *http.Request) {

	c := make(chan []byte)
	go fetchCurriencies(c)

	responseData, ok := <-c

	if !ok {
		http.Error(w, "Something went wrong!", http.StatusInternalServerError)
		return
	}

	var responseObject CurrencyData
	json.Unmarshal(responseData, &responseObject)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responseObject)
}

// fetchCurriencies will fetch the curriencies from Spenmo server and push it to the channel c
// It closes the channel in case of any error
func fetchCurriencies(c chan []byte) {
	response, err := http.Get("https://apiv1.qa.spenmo.com/api/v1/enumeration/currencies")
	if err != nil {
		fmt.Print(err.Error())
		close(c)
	}
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Print(err.Error())
		close(c)
	}
	c <- data
	close(c)
}
