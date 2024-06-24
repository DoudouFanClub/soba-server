package server

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func AllowCors(w http.ResponseWriter) {
	// Specify Content Type to receive as Json Format
	w.Header().Set("Content-Type", "application/json")
	// Set CORS headers to allow requests from all origins - Different Ports
	w.Header().Set("Access-Control-Allow-Origin", "*")
	// Allow Content-Type header
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
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

	// Write the JSON response
	_, err = w.Write(jsonResponse)
	if err != nil {
		fmt.Println("Error writing response:", err)
		return false
	}

	return true
}