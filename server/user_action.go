package server

import (
	"encoding/json"
	"fmt"
	"llm_server/balancer"
	"llm_server/cache"
	"llm_server/database"
	"llm_server/helper"
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
		
		var removeConvo database.RemoveConversation
		if err := json.NewDecoder(r.Body).Decode(&removeConvo); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			fmt.Println("Unable to receive a valid payload")
			return
		}
		
		if convoExist := mongoClient.DoesConvoExist(removeConvo.Username, removeConvo.Title); convoExist {
			if delErr := mongoClient.DeleteConversation(removeConvo.Username, removeConvo.Title); delErr != nil {
				fmt.Println("Unable to delete conversation from ConversationData...", delErr)
			}
		}

		if redisErr := redisClient.UnloadConversation(removeConvo.Username, removeConvo.Title); redisErr != nil {
			fmt.Println("Unable to unload conversation from redis... Conversation may not be active.", redisErr)
		}
	}
}

func handleSendMessage(mongoClient *database.MongoInterface, redisClient *cache.RedisCache, b *balancer.Balancer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// New
		w.Header().Set("Content-Type", "text/event-stream")
		//w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		AllowCors(w)
		defer r.Body.Close()

		if r.Method == "OPTIONS" {
			return
		}

		//var userPrompt database.FrontendMessagesPrompt
		var userPrompt database.MessagePrompt
		if err := json.NewDecoder(r.Body).Decode(&userPrompt); err != nil {
			http.Error(w, "Error decoding prompt", http.StatusBadRequest)
			return
		}
		fmt.Println("User Prompt: ", userPrompt)
		// Retrieve the last message (User Prompt) & Insert it into Redis Cache
		redisClient.AddMessageToConversation(mongoClient, userPrompt.Username, userPrompt.Title, database.Message{Role: "user", Content: userPrompt.Contents.Content})
		fmt.Println("added message to convo", userPrompt.Contents.Content)

		// Rearranging index of titles
		userTitles := mongoClient.RetrieveConversationTitles(userPrompt.Username)
		currentTitleIndex := mongoClient.RetrieveTitleIndex(userPrompt.Username, userPrompt.Title, userTitles)
		updatedTitles := helper.MoveCurrentIndexToFront(userTitles, currentTitleIndex)
		mongoClient.UpdateUser(userPrompt.Username, updatedTitles)

		convo, err := redisClient.GetDataConversation(mongoClient, userPrompt.Username, userPrompt.Title)
		if err != nil {
			fmt.Println("unable to retrieve conversation data from:", userPrompt.Title, " | err:", err)
		}

		// Send to LLM
		success, result := ReceiveMessage(&w, convo.Messages, b) // have to replace this later
		// Wait and append to cache once done
		if success {
			redisClient.AddMessageToConversation(mongoClient, userPrompt.Username, userPrompt.Title, database.Message{Role: "assistant", Content: result})
		} else {
			fmt.Println("Failed to retrieve complete prompt from LLM")
		}
	}
}
