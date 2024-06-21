package server

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func AllowCors(w http.ResponseWriter) {

	w.Header().Set("Content-Type", "application/json")

	// Set CORS headers to allow requests from all origins
	w.Header().Set("Access-Control-Allow-Origin", "*")
	//w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS") // Allow POST requests
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type") // Allow Content-Type header
}


func ParseJson(w http.ResponseWriter, jsonField ...map[string]interface{}) bool {

	var responseData map[string]interface{} = make(map[string]interface{})

	for _, jsonParam := range jsonField {
		for key, value := range jsonParam {
			responseData[key] = value
		}
	}

	jsonResponse, err := json.Marshal(responseData)
	if err != nil {
		// Return a 500 Internal Server Error response if there's an error encoding the response data
		http.Error(w, fmt.Sprintf("Error encoding response: %s", err.Error()), http.StatusInternalServerError)
		return false
	}

	// Write the JSON response to the http.ResponseWriter
	_, err = w.Write(jsonResponse)
	if err != nil {
		// Handle the error if unable to write the response
		fmt.Println("Error writing response:", err)
		return false
	}

	return true
}