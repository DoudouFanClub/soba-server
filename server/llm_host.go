package server

import (
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
	mongoClient, err_1 := database.CreateMongoInterface(databaseUri)
	if err_1 != nil {
		return nil, err_1
	}

	// Init Redis Cache Connection
	redisClient, err_2 := cache.CreateRedisInterface(redisAddr, redisPassword, db)
	if err_2 != nil {
		return nil, err_2
	}

	// Server Access Callbacks
	http.HandleFunc("/login", handleLogin(mongoClient))
	http.HandleFunc("/logout", handleLogout(mongoClient, redisClient))
	http.HandleFunc("/register", handleRegister(mongoClient))

	// User Action Callbacks
	http.HandleFunc("/new_chat", handleNewChat(mongoClient))
	//http.HandleFunc("/rename_chat", handleRenameChat(mongoClient, redisClient))
	//http.HandleFunc("/delete_chat", handleRenameChat(mongoClient, redisClient))
	http.HandleFunc("/send_message", handleSendMessage(mongoClient, redisClient))
	
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
