package server

import (
	"encoding/json"
	"fmt"
	"llm_server/cache"
	"llm_server/database"
	"net/http"
)

func handleLogin(mongoClient *database.MongoInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		AllowCors(w)
		defer r.Body.Close()

		if r.Method == "OPTIONS" {
			return
		}

		var user database.UserData
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		var loginStatus string
		loginValid := mongoClient.IsUserLoginValid(user)
		if loginValid {
			loginStatus = "success"
		} else if mongoClient.DoesUserExist(user.Username) {
			loginStatus = "failure"
		} else {
			loginStatus = "invalid"
		}

		successfulWrite := ParseJson(w, map[string]interface{}{"response": loginStatus})
		if successfulWrite {
			fmt.Println("Login attempt:", user.Username, loginStatus)
		} else {
			fmt.Println("Failed to respond to login attempt for:", user.Username)
		}
	}
}

func handleLogout(mongoClient *database.MongoInterface, redisClient *cache.RedisCache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		AllowCors(w)
		defer r.Body.Close()

		if r.Method == "OPTIONS" {
			return
		}

		var convo database.ConversationRequest
		if err := json.NewDecoder(r.Body).Decode(&convo); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		var updatedDatabaseOnLogout bool
		redisConversation, err := redisClient.GetDataConversation(mongoClient, convo.Username, convo.Title)
		if err != nil {
			fmt.Println("Warning: Unable to retrieve conversation from Redis Cache:", convo.Title, err)
			updatedDatabaseOnLogout = false
		} else {
			mongoClient.InsertConversation(convo.Username, *redisConversation)
			redisClient.UnloadConversation(convo.Username, convo.Title)
			updatedDatabaseOnLogout = true
		}

		successfulWrite := ParseJson(w, map[string]interface{}{"response": updatedDatabaseOnLogout})
		if successfulWrite {
			fmt.Println("Logout attempt:", convo.Username, "| Database updated:", successfulWrite)
		} else {
			fmt.Println("Failed to process logout attempt for:", convo.Username)
		}
	}
}

func handleRegister(mongoClient *database.MongoInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		AllowCors(w)
		defer r.Body.Close()

		if r.Method == "OPTIONS" {
			return
		}

		var user database.UserData
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		var serverResponse string
		userExist := mongoClient.DoesUserExist(user.Username)
		if !userExist {
			mongoClient.InsertUser(user.Username, user.Password)
			serverResponse = "success"
		} else {
			serverResponse = "failure"
		}

		successfulWrite := ParseJson(w, map[string]interface{}{"response": serverResponse})
		if successfulWrite {
			fmt.Println("Registration attempt:", user.Username, serverResponse)
		} else {
			fmt.Println("Failed to respond to registration attempt for:", user.Username)
		}
	}
}