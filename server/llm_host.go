package server

import (
	"encoding/json"
	"fmt"
	"llm_server/balancer"
	"llm_server/cache"
	"llm_server/database"
	"llm_server/socket"
	"net/http"
	"os"
)

type LLM_Host struct {
	MongoClient *database.MongoInterface
	RedisClient *cache.RedisCache
}

type EndpointConfig struct {
	EndpointsData []socket.Endpoint `json:"endpoints_data"`
}

func InitLLMHost(databaseUri string, redisAddr string, redisPassword string, db int) (*LLM_Host, error) {

	// Init MongoDb Connection
	mongoClient, err_1 := database.CreateMongoInterface(databaseUri)
	if err_1 != nil {
		return nil, err_1
	}

	// Init Redis Cache Connection
	redisClient, err_2 := cache.CreateRedisInterface(redisAddr, redisPassword, db)
	if err_2 != nil {
		return nil, err_2
	}

	b := balancer.CreateBalancer()
	

	file, _ := os.Open("config/endpoints.cfg")
	defer file.Close()
	decoder := json.NewDecoder(file)
	configuration := EndpointConfig{}
	err := decoder.Decode(&configuration)
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Println(configuration)
	
	for _, val := range configuration.EndpointsData {
		b.Add(val)
	}
	
	// Server Access Callbacks
	http.HandleFunc("/login", handleLogin(mongoClient))
	http.HandleFunc("/logout", handleLogout(mongoClient, redisClient))
	http.HandleFunc("/register", handleRegister(mongoClient))

	// User Action Callbacks
	http.HandleFunc("/new_chat", handleNewChat(mongoClient))
	//http.HandleFunc("/rename_chat", handleRenameChat(mongoClient, redisClient))
	//http.HandleFunc("/delete_chat", handleRenameChat(mongoClient, redisClient))
	http.HandleFunc("/send_message", handleSendMessage(mongoClient, redisClient, b))
	
	// Server Action Callbacks
	http.HandleFunc("/load_chat", handleLoadChat(mongoClient, redisClient))
	http.HandleFunc("/retrieve_convo_titles", handleRetrieveConvoTitles(mongoClient))

	return &LLM_Host{
		MongoClient: mongoClient,
		RedisClient: redisClient,
	}, nil
}

func (h *LLM_Host) Tick() {
	http.ListenAndServe(":8080", nil)
}
