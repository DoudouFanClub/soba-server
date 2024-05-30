package main

/*
Database Structure
	UserData
		User1 Struct
		User2 Struct
	ConversationData
		User1 Username String
			0
			1
			2
		User2 Username String
			0
			1
*/

import (
	"fmt"
	"llm_server/database"
	"llm_server/server"
)

func main() {

	mongoDbUri := "mongodb://localhost:27017/"
	redisUri := "localhost:6379"

	host, err := server.InitLLMHost(mongoDbUri, redisUri, "", 0)

	if err != nil {
		fmt.Println(err)
		return
	}

	convo := database.Conversation{
		Title:    "MyConvo",
		Messages: make([]database.Message, 0),
	}

	host.RedisClient.ClearRedisMongoInterface()

	msg_1 := database.Message{Role: "jp", Content: "jp says hi"}
	msg_2 := database.Message{Role: "ryan", Content: "ryan says bye"}
	msg_3 := database.Message{Role: "jeremy", Content: "jeremy says bye"}
	convo.Messages = append(convo.Messages, msg_1, msg_2, msg_3)

	host.MongoClient.InsertUser("jp", "1password")
	host.MongoClient.InsertConversation("jp", convo)

	red_err := host.RedisClient.LoadConversation(host.MongoClient, "jp", "MyConvo")
	if red_err != nil {
		fmt.Println(err)
	}
	host.RedisClient.AddMessageToConversation(host.MongoClient, "jp", "MyConvo", database.Message{Role: "mae", Content: "i like to eat kinder bueno"})

	host.Tick()

	// http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
	// 	server.RegisterUser(w, r, host.MongoClient)
	// })

	// http.ListenAndServe(":8080", nil)

	defer func() {
		if err = host.MongoClient.Terminate(); err != nil {
			panic(err)
		}
		// if err = redis_svr.Terminate(); err != nil {
		// 	panic(err)
		// }
	}()
}
