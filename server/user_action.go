package server

import (
	"encoding/json"
	"fmt"
	"llm_server/cache"
	"llm_server/database"
	"net/http"
)

func handleNewChat(mongoClient *database.MongoInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		AllowCors(w)
		defer r.Body.Close()

		if r.Method == "OPTIONS" {
			return
		}

		var createConvo database.ConversationRequest
		if err := json.NewDecoder(r.Body).Decode(&createConvo); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		var serverResponse string
		convoExist := mongoClient.DoesConvoExist(createConvo.Username, createConvo.Title)
		if !convoExist {
			convo := database.Conversation{
				Title:    createConvo.Title,
				Messages: make([]database.Message, 0),
			}
			mongoClient.InsertConversation(createConvo.Username, convo)
			serverResponse = "success"
		} else {
			serverResponse = "failure"
		}

		successfulWrite := ParseJson(w, map[string]interface{}{"response": serverResponse})
		if successfulWrite {
			fmt.Println("New chat creation attempt:", createConvo.Username, serverResponse)
		} else {
			fmt.Println("Failed to respond to new chat creation attempt for:", createConvo.Username)
		}
	}
}

func handleRenameChat(mongoClient *database.MongoInterface, redisClient *cache.RedisCache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		AllowCors(w)
		defer r.Body.Close()

		if r.Method == "OPTIONS" {
			return
		}
	}
}

func handleDeleteChat(mongoClient *database.MongoInterface, redisClient *cache.RedisCache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		AllowCors(w)
		defer r.Body.Close()
		
		if r.Method == "OPTIONS" {
			return
		}

		// check if loaded in cache
		//     if yes, clear username_title from cache
		// clear database
	}
}

func handleSendMessage(mongoClient *database.MongoInterface, redisClient *cache.RedisCache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		AllowCors(w)
		defer r.Body.Close()

		if r.Method == "OPTIONS" {
			return
		}

		var userPrompt database.FrontendMessagesPrompt
		if err := json.NewDecoder(r.Body).Decode(&userPrompt); err != nil {
			http.Error(w, "Error decoding prompt", http.StatusBadRequest)
			return
		}

		// Retrieve the last message (User Prompt) & Insert it into Redis Cache
		lastUserPrompt := userPrompt.Contents[len(userPrompt.Contents)-1].Content
		redisClient.AddMessageToConversation(mongoClient, userPrompt.Username, userPrompt.Title, database.Message{Role: "User", Content: lastUserPrompt})
		fmt.Println("added message to convo")

		// Send to LLM
		success, result := ReceiveMessage(&w, r, &b)
		// Wait and append to cache once done
		if success {
			redisClient.AddMessageToConversation(mongoClient, userPrompt.Username, userPrompt.Title, database.Message{Role: "assistant", Content: result})
		} else {
			fmt.Println("Failed to retrieve complete prompt from LLM")
		}
	}
}