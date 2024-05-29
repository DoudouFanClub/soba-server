package main

//const uri = "mongodb://localhost:27017/"

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

/*

Frontend:
- Angular

Backend:
- Database
  - MongoDB
- Cache the current Chat History
  - Redis
- Message Queue (Handle multiple users sending prompts to Local LLM)
  -

*/

// func main() {

// 	MongoInterface, err := database.CreateMongoMongoInterface(uri)

// 	MongoInterface.InsertUser("ryan", "wewewewewewe wew")

// 	status := MongoInterface.DoesUserExist("jeremy")
// 	status2 := MongoInterface.DoesUserExist("ryan")

// 	if status {
// 		fmt.Println("Jeremy was here")
// 	} else {
// 		fmt.Println("Uh oh he ded")
// 	}

// 	if status2 {
// 		fmt.Println("Ryan was here")
// 	} else {
// 		fmt.Println("Uh oh he ded too")
// 		MongoInterface.InsertUser("ryan", "wewewewewewe wew")
// 	}

// 	if err != nil {
// 		panic(err)
// 	}

// 	defer func() {
// 		if err = MongoInterface.Terminate(); err != nil {
// 			panic(err)
// 		}
// 	}()
// }

import (
	"context"
	"fmt"
	"llm_server/cache"
	"llm_server/database"
	"llm_server/server"
	"net/http"
)

var ctx = context.Background()

type Person struct {
	Name  string
	Email string
}

const uri = "mongodb://localhost:27017/"

func main() {

	mongo_svr, err := database.CreateMongoMongoInterface(uri)
	redis_svr, _ := cache.CreateRedisMongoInterface("localhost:6379", "", 0)

	if err != nil {
		fmt.Println(err)
		return
	}

	convo := database.Conversation{
		ConversationId: 0,
		Messages:       make([]database.Message, 0),
	}

	redis_svr.ClearRedisMongoInterface()

	msg_1 := database.Message{Role: "jp", Content: "jp says hi"}
	msg_2 := database.Message{Role: "ryan", Content: "ryan says bye"}
	msg_3 := database.Message{Role: "jeremy", Content: "jeremy says bye"}
	convo.Messages = append(convo.Messages, msg_1, msg_2, msg_3)

	mongo_svr.InsertUser("jp", "1password")
	mongo_svr.InsertConversation("jp", convo)

	red_err := redis_svr.LoadConversation(mongo_svr, "jp", 0)
	if red_err != nil {
		fmt.Println(err)
	}
	redis_svr.AddMessageToConversation(mongo_svr, "jp", 0, database.Message{Role: "mae", Content: "i like to eat kinder bueno"})

	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		server.RegisterUser(w, r, mongo_svr)
	})

	http.ListenAndServe(":8080", nil)

	defer func() {
		if err = mongo_svr.Terminate(); err != nil {
			panic(err)
		}
		// if err = redis_svr.Terminate(); err != nil {
		// 	panic(err)
		// }
	}()

	// rdb := redis.NewClient(&redis.Options{
	// 	Addr:     "localhost:6379",
	// 	Password: "", // no password set
	// 	DB:       0,  // use default DB
	// })
	// person := Person{
	// 	Name:  "Ace",
	// 	Email: "aceacett@simp.com",
	// }

	// json_bytes, err := json.Marshal(person)
	// if err != nil {
	// 	fmt.Println("noob")
	// 	return
	// }

	// err2 := rdb.Set(ctx, "key", json_bytes, 0).Err()
	// if err2 != nil {
	// 	panic(err)
	// }

	// val, err := rdb.Get(ctx, "key").Result()
	// if err != nil {
	// 	panic(err)
	// }
	// var data interface{}
	// json.Unmarshal([]byte(val), &data)
	// fmt.Println("key", data)

	// val2, err := rdb.Get(ctx, "key2").Result()
	// if err == redis.Nil {
	// 	fmt.Println("key2 does not exist")
	// } else if err != nil {
	// 	panic(err)
	// } else {
	// 	fmt.Println("key2", val2)
	// }
}
