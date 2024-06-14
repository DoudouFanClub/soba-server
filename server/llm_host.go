package server

import (
	"encoding/json"
	"fmt"
	"llm_server/cache"
	"llm_server/database"
	"net/http"
)

type LLM_Host struct {
	MongoClient *database.MongoInterface
	RedisClient *cache.RedisCache
}

func InitLLMHost(databaseUri string, redisAddr string, redisPassword string, db int) (*LLM_Host, error) {

	// Init MongoDb Connection
	mongo_svr, err_1 := database.CreateMongoMongoInterface(databaseUri)
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
		w.Header().Set("Content-Type", "application/json")

        // Set CORS headers to allow requests from all origins
        w.Header().Set("Access-Control-Allow-Origin", "*")
        //w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS") // Allow POST requests
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type") // Allow Content-Type header

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
				fmt.Println("User exists, BUT incorrect password")
				loginStatus = "failure"
			} else {
				fmt.Println("Would you like to create a new account")
				loginStatus = "invalid"
			}
		}

		responseData := map[string]interface{}{"response": loginStatus,}
		jsonResponse, err := json.Marshal(responseData)
		if err != nil {
			// Return a 500 Internal Server Error response if there's an error encoding the response data
			http.Error(w, fmt.Sprintf("Error encoding response: %s", err.Error()), http.StatusInternalServerError)
			return
		}

		// Write the JSON response to the http.ResponseWriter
		_, err = w.Write(jsonResponse)
		if err != nil {
			// Handle the error if unable to write the response
			fmt.Println("Error writing response:", err)
			return
		}
	})

	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

        // Set CORS headers to allow requests from all origins
        w.Header().Set("Access-Control-Allow-Origin", "*")
        //w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS") // Allow POST requests
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type") // Allow Content-Type header

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

		

		// call redis fn clear the cache & store the updated convo into DB
	})

	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

        // Set CORS headers to allow requests from all origins
        w.Header().Set("Access-Control-Allow-Origin", "*")
        //w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS") // Allow POST requests
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type") // Allow Content-Type header

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
		responseData := map[string]string{};

		// Handle Registration of User
		if !userExist {
			fmt.Println("user does not yet exist") // Remove afterwards
			
			mongo_svr.InsertUser(user.Username, user.Password)
			responseData["response"] = "success"
			} else {
				fmt.Println("user already exists:", user.Username) // Remove afterwards
				responseData["response"] = "failure"
		}

		// Send a response to the frontend
		jsonResponse, err := json.Marshal(responseData)
		if err != nil {
			// Return a 500 Internal Server Error response if there's an error encoding the response data
			http.Error(w, fmt.Sprintf("Error encoding response: %s", err.Error()), http.StatusInternalServerError)
			return
		}

		// Write the JSON response to the http.ResponseWriter
		_, err = w.Write(jsonResponse)
		if err != nil {
			// Handle the error if unable to write the response
			fmt.Println("Error writing response:", err)
			return
		}
	})

	http.HandleFunc("/new_chat", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

        // Set CORS headers to allow requests from all origins
        w.Header().Set("Access-Control-Allow-Origin", "*")
        //w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS") // Allow POST requests
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type") // Allow Content-Type header

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

		if !convo_exist {
			convo := database.Conversation{
				Title:    create_convo.Title,
				Messages: make([]database.Message, 0),
			}
			mongo_svr.InsertConversation(create_convo.Username, convo)
			fmt.Println("Created a new conversation")
		} else {
			fmt.Println("Conversation already exists")
		}
	})

	http.HandleFunc("/rename_chat", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

        // Set CORS headers to allow requests from all origins
        w.Header().Set("Access-Control-Allow-Origin", "*")
        //w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS") // Allow POST requests
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type") // Allow Content-Type header

		// Close the request body after reading from it
		defer r.Body.Close()

		// Handle preflight OPTIONS requests
		if r.Method == "OPTIONS" {
			return
		}


	})

	http.HandleFunc("/load_chat", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

        // Set CORS headers to allow requests from all origins
        w.Header().Set("Access-Control-Allow-Origin", "*")
        //w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS") // Allow POST requests
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type") // Allow Content-Type header

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
		//responseData := map[string]interface{}{"response": convo.Messages,}
		jsonResponse, err := json.Marshal(convo)
		if err != nil {
			// Return a 500 Internal Server Error response if there's an error encoding the response data
			http.Error(w, fmt.Sprintf("Error encoding response: %s", err.Error()), http.StatusInternalServerError)
			return
		}

		// Write the JSON response to the http.ResponseWriter
		_, err = w.Write(jsonResponse)
		if err != nil {
			// Handle the error if unable to write the response
			fmt.Println("Error writing response:", err)
			return
		}
	})

	http.HandleFunc("/delete_chat", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

        // Set CORS headers to allow requests from all origins
        w.Header().Set("Access-Control-Allow-Origin", "*")
        //w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS") // Allow POST requests
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type") // Allow Content-Type header

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

	http.HandleFunc("/user_query", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

        // Set CORS headers to allow requests from all origins
        w.Header().Set("Access-Control-Allow-Origin", "*")
        //w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS") // Allow POST requests
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type") // Allow Content-Type header

		// Close the request body after reading from it
		defer r.Body.Close()

		// Handle preflight OPTIONS requests
		if r.Method == "OPTIONS" {
			return
		}

		var user_prompt database.MessagePrompt
		err := json.NewDecoder(r.Body).Decode(&user_prompt)

		if err != nil {
			fmt.Println("Error decoding prompt: %w", err.Error())
			return
		}

		redis_svr.AddMessageToConversation(mongo_svr, user_prompt.Username, user_prompt.Title, database.Message{Role: "User", Content: user_prompt.Content})

		// send to LLM

		// Wait and append to cache once done
		// Also forward the data to frontend
	})

	http.HandleFunc("/retrieve_convo_titles", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

        // Set CORS headers to allow requests from all origins
        w.Header().Set("Access-Control-Allow-Origin", "*")
        //w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS") // Allow POST requests
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type") // Allow Content-Type header

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
		responseData := map[string]interface{}{"response": titles,}
		jsonResponse, err := json.Marshal(responseData)
		if err != nil {
			// Return a 500 Internal Server Error response if there's an error encoding the response data
			http.Error(w, fmt.Sprintf("Error encoding response: %s", err.Error()), http.StatusInternalServerError)
			return
		}

		// Write the JSON response to the http.ResponseWriter
		_, err = w.Write(jsonResponse)
		if err != nil {
			// Handle the error if unable to write the response
			fmt.Println("Error writing response:", err)
			return
		}
	})

	
	http.HandleFunc("/testpost", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

        // Set CORS headers to allow requests from all origins
        w.Header().Set("Access-Control-Allow-Origin", "*")
        //w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS") // Allow POST requests
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type") // Allow Content-Type header

		// Close the request body after reading from it
		defer r.Body.Close()

		// Handle preflight OPTIONS requests
		if r.Method == "OPTIONS" {
			return
		}

		var user_prompt database.FrontendTest
		err := json.NewDecoder(r.Body).Decode(&user_prompt)

		if err != nil {
			fmt.Println("Error decoding prompt: %w", err.Error())
			fmt.Println(user_prompt)
			fmt.Println("returning cause of an error gg", r.Body)
			return
		}
		
		fmt.Println("Received an input:", user_prompt.Text)





		// Send a response to the frontend
		responseData := map[string]string{"message": "Data received successfully"}
		jsonResponse, err := json.Marshal(responseData)
		if err != nil {
			// Return a 500 Internal Server Error response if there's an error encoding the response data
			http.Error(w, fmt.Sprintf("Error encoding response: %s", err.Error()), http.StatusInternalServerError)
			return
		}

		// Write the JSON response to the http.ResponseWriter
		_, err = w.Write(jsonResponse)
		if err != nil {
			// Handle the error if unable to write the response
			fmt.Println("Error writing response:", err)
			return
		}


		// send to LLM

		// Wait and append to cache once done
		// Also forward the data to frontend
	})

	return &LLM_Host{
		MongoClient: mongo_svr,
		RedisClient: redis_svr,
	}, nil
}

func (h *LLM_Host) Tick() {
	http.ListenAndServe(":8080", nil)
}
