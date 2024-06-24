package server

import (
	"encoding/json"
	"fmt"
	"llm_server/cache"
	"llm_server/database"
	"net/http"
)

func handleLoadChat(mongoClient *database.MongoInterface, redisClient *cache.RedisCache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		AllowCors(w)
		defer r.Body.Close()

		if r.Method == "OPTIONS" {
			return
		}

		var loadConvo database.ConversationRequest
		if err := json.NewDecoder(r.Body).Decode(&loadConvo); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		var convo database.Conversation
		convoExist := mongoClient.DoesConvoExist(loadConvo.Username, loadConvo.Title)
		if convoExist {
			redisClient.LoadConversation(mongoClient, loadConvo.Username, loadConvo.Title)
			convo = mongoClient.GetConvo(loadConvo.Username, loadConvo.Title)
		}

		successfulWrite := ParseJson(w, map[string]interface{}{"response": convo})
		if successfulWrite {
			fmt.Println("Chat load attempt:", loadConvo.Username, " | Title: ", loadConvo.Title, " | Status:", convoExist)
		} else {
			fmt.Println("Failed to respond to chat load attempt for:", loadConvo.Username, " | Title: ", loadConvo.Title)
		}
	}
}

func handleRetrieveConvoTitles(mongoClient *database.MongoInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		AllowCors(w)
		defer r.Body.Close()

		if r.Method == "OPTIONS" {
			return
		}

		var userRequest database.ConversationTitlesRequest
		if err := json.NewDecoder(r.Body).Decode(&userRequest); err != nil {
			http.Error(w, "Error retrieving conversation titles", http.StatusBadRequest)
			return
		}

		titles := mongoClient.RetrieveConversationTitles(userRequest.Username)

		successfulWrite := ParseJson(w, map[string]interface{}{"response": titles})
		if successfulWrite {
			fmt.Println("Successfully retrieved conversation titles for:", userRequest.Username)
		} else {
			fmt.Println("Failed to retrieve conversation titles for:", userRequest.Username)
		}
	}
}