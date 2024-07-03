package server

import (
	"encoding/json"
	"fmt"
	"llm_server/balancer"
	"llm_server/database"
	"net/http"
	"time"
)

/*
	rest api
	- post request is encrypted
*/

// Register user also verifies user
func RegisterUser(w http.ResponseWriter, r *http.Request, m *database.MongoInterface) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	var user database.UserData

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	insert_err := m.InsertUser(user.Username, user.Password)
	if insert_err != nil {
		http.Error(w, insert_err.Error(), http.StatusBadRequest)
		return
	}
}

// Load chat upon user click

// Receive chat from user (this endpoint also returns a stream of html objects based on what llm replies)
// This also needs to manage the availability of the llm inferer
// pseudo load balancer and message queue needs to be used here
// caching the users prompt and the llm's reply
func ReceiveMessage(w *http.ResponseWriter, userChatHistory []database.Message, b *balancer.Balancer) (bool, string) {

	fmt.Println("Receive Message History:", userChatHistory)
	chatHistoryBytes, _ := json.Marshal(userChatHistory)
	chatHistoryStr := string(chatHistoryBytes) + "<doudousmacksoba>"
	// this shouldn't be a blocking call as the net/http documentation
	// says that each request is done on a different goroutine
	tries := 20
	for !b.Available() {
		// wait 3 seconds before polling again
		time.Sleep(3 * time.Second)
		tries -= 1
		// keep trying to send for 1 minute, else drop out temporarily
		if tries < 0 {
			fmt.Println("Connection timed out")
			return false, ""
		}
	}

	success, msg := b.Send([]byte(chatHistoryStr), w)
	if !success {
		return false, msg
	}

	return success, msg
}
