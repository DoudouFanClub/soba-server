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
		var user database.UserData
		err := json.NewDecoder(r.Body).Decode(&user)

		if err != nil {
			fmt.Println(err)
			return
		}

		loginValid := mongo_svr.IsUserLoginValid(user)

		// Login Invalid, check if wrong password
		if loginValid {
			// Load Conversations
			fmt.Println("Successful login")
		} else {
			userExist := mongo_svr.DoesUserExist(user.Username)
			if userExist {
				// Prompt user wrong password
				userExist = true
				fmt.Println("User exists, BUT incorrect password")
			} else {
				fmt.Println("Would you like to create a new account")
			}
		}
	})

	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		var user database.UserData
		err := json.NewDecoder(r.Body).Decode(&user)

		if err != nil {
			fmt.Println(err)
			return
		}

		userExist := mongo_svr.DoesUserExist(user.Username)

		if !userExist {
			fmt.Println("user does not yet exist")
			mongo_svr.InsertUser(user.Username, user.Password)
		} else {
			fmt.Println("user already exists:", user.Username)
		}
	})

	http.HandleFunc("/new_chat", func(w http.ResponseWriter, r *http.Request) {
		// force users to insert their chat name to facilitate loading of conversation titles on the sidebar
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

	http.HandleFunc("/load_chat", func(w http.ResponseWriter, r *http.Request) {
		var load_convo database.ConversationRequest
		err := json.NewDecoder(r.Body).Decode(&load_convo)

		if err != nil {
			fmt.Println("Error loading chat: %w", err.Error())
			return
		}

		// load the convo into cache
		convo_exist := mongo_svr.DoesConvoExist(load_convo.Username, load_convo.Title)
		if convo_exist {
			fmt.Println("Loaded chat")
			redis_svr.LoadConversation(mongo_svr, load_convo.Username, load_convo.Title)
		} else {
			fmt.Println("Chat does not exist")
		}
	})

	http.HandleFunc("/delete_chat", func(w http.ResponseWriter, r *http.Request) {

	})

	http.HandleFunc("/user_query", func(w http.ResponseWriter, r *http.Request) {
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
