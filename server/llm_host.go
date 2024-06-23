package server

import (
	"encoding/json"
	"fmt"
	"llm_server/balancer"
	"llm_server/cache"
	"llm_server/database"
	"net/http"
)

type LLM_Host struct {
	MongoClient *database.MongoInterface
	RedisClient *cache.RedisCache
}

var b = balancer.CreateBalancer()

func InitLLMHost(databaseUri string, redisAddr string, redisPassword string, db int) (*LLM_Host, error) {

	// Init MongoDb Connection
	mongo_svr, err_1 := database.CreateMongoInterface(databaseUri)
	if err_1 != nil {
		return nil, err_1
	}

	// Init Redis Cache Connection
	redis_svr, err_2 := cache.CreateRedisInterface(redisAddr, redisPassword, db)
	if err_2 != nil {
		return nil, err_2
	}

	// Init net/http Callbacks
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		AllowCors(w)

		// Close the request body after reading from it
		defer r.Body.Close()

		// Handle preflight OPTIONS requests
		if r.Method == "OPTIONS" {
			return
		}

		var user database.UserData
		err := json.NewDecoder(r.Body).Decode(&user)

		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(user.Username, user.Password, user.ConversationIDs)

		loginValid := mongo_svr.IsUserLoginValid(user)

		var loginStatus string

		// Login Invalid, check if wrong password
		if loginValid {
			// Load Conversations
			fmt.Println("Successful login")
			loginStatus = "success"
		} else {
			userExist := mongo_svr.DoesUserExist(user.Username)
			if userExist {
				// Prompt user wrong password
				userExist = true
				loginStatus = "failure"
			} else {
				loginStatus = "invalid"
			}
		}

		// Respond to the Frontend that Login was "loginStatus"
		successfulWrite := ParseJson(w, map[string]interface{}{"response": loginStatus})
		if successfulWrite {
			fmt.Println("Successfully logged in:", user.Username)
		} else {
			fmt.Println("Failed to log in:", user.Username)
		}
	})

	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		AllowCors(w)

		// Close the request body after reading from it
		defer r.Body.Close()

		// Handle preflight OPTIONS requests
		if r.Method == "OPTIONS" {
			return
		}

		var convo database.ConversationRequest
		err := json.NewDecoder(r.Body).Decode(&convo)

		if err != nil {
			fmt.Println(err)
			return
		}

		// Unload the Redis Cache for "Useranme"
		// And updates the MongoDB with the latest conversation
		var updatedDatabaseOnLogout bool
		redisConversation, err := redis_svr.GetDataConversation(mongo_svr, convo.Username, convo.Title)
		if err != nil {
			fmt.Println("Warning: Unable to retrieve conversation from Redis Cache to update Database: ", convo.Title)
			fmt.Println(err)
			updatedDatabaseOnLogout = false
		} else {
			mongo_svr.InsertConversation(convo.Username, *redisConversation)
			redis_svr.UnloadConversation(convo.Username, convo.Title)
			updatedDatabaseOnLogout = true
		}
		
		// Respond to the Frontend that Login was "loginStatus"
		successfulWrite := ParseJson(w, map[string]interface{}{"response": updatedDatabaseOnLogout})
		if successfulWrite {
			fmt.Println("Successfully logged out:", convo.Username)
		} else {
			fmt.Println("Failed to log in:", convo.Username)
		}
	})

	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		AllowCors(w)

		// Close the request body after reading from it
		defer r.Body.Close()

		// Handle preflight OPTIONS requests
		if r.Method == "OPTIONS" {
			return
		}

		var user database.UserData
		err := json.NewDecoder(r.Body).Decode(&user)

		if err != nil {
			fmt.Println(err)
			return
		}

		userExist := mongo_svr.DoesUserExist(user.Username)

		// Create a Response Json Format
		var serverResponse string

		// Handle Registration of User
		if !userExist {
			fmt.Println("user does not yet exist") // Remove afterwards

			mongo_svr.InsertUser(user.Username, user.Password)
			serverResponse = "success"
		} else {
			fmt.Println("user already exists:", user.Username) // Remove afterwards
			serverResponse = "failure"
		}

		// Send a response to the frontend
		successfulWrite := ParseJson(w, map[string]interface{}{"response": serverResponse})
		if successfulWrite {
			fmt.Println("Successfully logged in:", user.Username)
		} else {
			fmt.Println("Failed to log in:", user.Username)
		}
	})

	http.HandleFunc("/new_chat", func(w http.ResponseWriter, r *http.Request) {
		AllowCors(w)

		// Close the request body after reading from it
		defer r.Body.Close()

		// Handle preflight OPTIONS requests
		if r.Method == "OPTIONS" {
			return
		}

		// Force users to insert their chat name to facilitate loading of conversation titles on the sidebar
		var create_convo database.ConversationRequest
		err := json.NewDecoder(r.Body).Decode(&create_convo)

		if err != nil {
			fmt.Println(err)
			return
		}

		convo_exist := mongo_svr.DoesConvoExist(create_convo.Username, create_convo.Title)

		// Create a Response Json Format
		var serverResponse string

		if !convo_exist {
			convo := database.Conversation{
				Title:    create_convo.Title,
				Messages: make([]database.Message, 0),
			}
			mongo_svr.InsertConversation(create_convo.Username, convo)
			serverResponse = "success"
			fmt.Println("Created a new conversation")
		} else {
			fmt.Println("Conversation already exists")
			serverResponse = "failure"
		}

		// Send a response to the frontend
		successfulWrite := ParseJson(w, map[string]interface{}{"response": serverResponse})
		if successfulWrite {
			fmt.Println("Successfully logged in:", create_convo.Username)
		} else {
			fmt.Println("Failed to log in:", create_convo.Username)
		}
	})

	http.HandleFunc("/rename_chat", func(w http.ResponseWriter, r *http.Request) {
		AllowCors(w)

		// Close the request body after reading from it
		defer r.Body.Close()

		// Handle preflight OPTIONS requests
		if r.Method == "OPTIONS" {
			return
		}

	})

	http.HandleFunc("/load_chat", func(w http.ResponseWriter, r *http.Request) {
		AllowCors(w)

		// Close the request body after reading from it
		defer r.Body.Close()

		// Handle preflight OPTIONS requests
		if r.Method == "OPTIONS" {
			return
		}

		var load_convo database.ConversationRequest
		err := json.NewDecoder(r.Body).Decode(&load_convo)

		if err != nil {
			fmt.Println("Error loading chat: %w", err.Error())
			return
		}

		// Load the convo into cache and retreve user's convo
		var convo database.Conversation
		convo_exist := mongo_svr.DoesConvoExist(load_convo.Username, load_convo.Title)
		if convo_exist {
			fmt.Println("Loaded chat")
			redis_svr.LoadConversation(mongo_svr, load_convo.Username, load_convo.Title)
			convo = mongo_svr.GetConvo(load_convo.Username, load_convo.Title)
		} else {
			fmt.Println("Chat does not exist")
		}

		// Send a response to the frontend
		successfulWrite := ParseJson(w, map[string]interface{}{"response": convo})
		if successfulWrite {
			fmt.Println("Successfully retrieved conversation for:", load_convo.Username, "  |  Title:", load_convo.Title)
		} else {
			fmt.Println("Failed to retrieve conversation for:", load_convo.Username, "  |  Title:", load_convo.Title)
		}
	})

	http.HandleFunc("/delete_chat", func(w http.ResponseWriter, r *http.Request) {
		AllowCors(w)

		// Close the request body after reading from it
		defer r.Body.Close()

		// Handle preflight OPTIONS requests
		if r.Method == "OPTIONS" {
			return
		}

		// check if loaded in cache
		//     if yes, clear username_title from cache
		// clear database
	})

	http.HandleFunc("/send_message", func(w http.ResponseWriter, r *http.Request) {
		AllowCors(w)

		// Close the request body after reading from it
		defer r.Body.Close()

		// Handle preflight OPTIONS requests
		if r.Method == "OPTIONS" {
			return
		}

		
		var user_prompt database.FrontendMessagesPrompt
		err := json.NewDecoder(r.Body).Decode(&user_prompt)
		
		if err != nil {
			fmt.Println("Error decoding prompt: %w", err.Error())
			return
		}
		
		// Retrieve the last message (User Prompt) & Insert it into Redis Cache
		LastUserPrompt := user_prompt.Contents[len(user_prompt.Contents) - 1].Content
		redis_svr.AddMessageToConversation(mongo_svr, user_prompt.Username, user_prompt.Title, database.Message{Role: "User", Content: LastUserPrompt})

		// Send to LLM
		success, result := ReceiveMessage(&w, r, &b)
		// Wait and append to cache once done
		if success {
			redis_svr.AddMessageToConversation(mongo_svr, user_prompt.Username, user_prompt.Title, database.Message{Role: "assistant", Content: result})
		} else {
			fmt.Println("noob")
		}

		// Still need to use http.writer to respond
	})

	http.HandleFunc("/retrieve_convo_titles", func(w http.ResponseWriter, r *http.Request) {
		AllowCors(w)

		// Close the request body after reading from it
		defer r.Body.Close()

		// Handle preflight OPTIONS requests
		if r.Method == "OPTIONS" {
			return
		}

		var userRequest database.ConversationTitlesRequest
		err := json.NewDecoder(r.Body).Decode(&userRequest)

		if err != nil {
			fmt.Println("Error Retrieving Conversation Titles: %w", err.Error())
			return
		}

		titles := mongo_svr.RetrieveConversationTitles(userRequest.Username)

		// Send a response to the frontend
		successfulWrite := ParseJson(w, map[string]interface{}{"response": titles})
		if successfulWrite {
			fmt.Println("Successfully retrieved conversation titles for:", userRequest.Username)
		} else {
			fmt.Println("Failed to retrieve conversation for:", userRequest.Username)
		}
	})

	return &LLM_Host{
		MongoClient: mongo_svr,
		RedisClient: redis_svr,
	}, nil
}

func (h *LLM_Host) Tick() {
	http.ListenAndServe(":8080", nil)
}
