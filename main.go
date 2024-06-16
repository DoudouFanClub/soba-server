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
	msg_4 := database.Message{Role: "jeremy", Content: "jeremy says byasdasdasdasdasdasdasdasdasdasdasde"}
	msg_5 := database.Message{Role: "jeremy", Content: "jeremy says byedas sssssss sssssssss ss sssssssss ssssssssss sssssss ss"}
	msg_6 := database.Message{Role: "jeremy", Content: "jeremy says bye asdasdasd sadasdas  asdasd asd asd asd asd asda sdas dasd asd asd dasdas dasd asd asd"}
	msg_7 := database.Message{Role: "jeremy", Content: "jeremy says bye sadasdas  asdasd asd asd asd asd asda sdas"}
	msg_8 := database.Message{Role: "jeremy", Content: "jeremy says bye dasdas  asdasd asd adasdas  asdasd asd a"}
	msg_9 := database.Message{Role: "jeremy", Content: "jeremy adasdas  asdas asdasadasdas  asdasadasdas  asdas"}
	msg_10 := database.Message{Role: "jeremy", Content: "jeremy says bye dasdas  asdasdasdas  asdasdasdas  asdas"}
	msg_11 := database.Message{Role: "jeremy", Content: "jeremy says byedasdas  asdasd asd asddasdas  asdasd asd asddasdas  asdasd asd asd"}
	msg_12 := database.Message{Role: "jeremy", Content: "jeremy says byesays bdassays byedasdas  asdassayedasdas  asys byedasdas  asdassays byedasdas  asdas"}
	msg_13 := database.Message{Role: "jeremy", Content: "jeremy says byeasdasadasdas  aasdasadasdas  aasdasadasdas  aasdasadasdas  aasdasadasdas  aasdasadasdas  aasdasadasdas  aasdasadasdas  a"}
	msg_14 := database.Message{Role: "jeremy", Content: "jeremy says byesays byedasdas  asdassayedasdas  asys byedasdassays byedasdas  asdassayedasdas  asys byedasdassays byedasdas  asdassayedasdas  asys byedasdassays byedasdas  asdassayedasdas  asys byedasdas"}
	msg_15 := database.Message{Role: "jeremy", Content: "jeremy says byeyed asdas  asdassayedasdyedasdas  asdassayedasdyedasdas  asdassayedasdyedasdas  asdassayedasdyedasdas  asdassayedasd"}
	msg_16 := database.Message{Role: "jeremy", Content: "jeremy says bye eyed asdas  asdassayeeyed asdas  asdassayeeyed asdas  asdassayeeyed asdas  asdassayeeyed asdas  asdassaye"}
	msg_17 := database.Message{Role: "jeremy", Content: "jeremy says bye emy says bye eyed asdasemy says bye eyed asdas"}
	msg_18 := database.Message{Role: "jeremy", Content: "jeremy says bye ys bye emy says bye eyed asdasys bye emy says bye eyed asdasys bye emy says bye eyed asdasys bye emy says bye eyed asdas"}
	msg_19 := database.Message{Role: "jeremy", Content: "jeremy says byes ays bye ys bye emy says bye eyed asdasays bye ys bye emy says bye eyed asdasays bye ys bye emy says bye eyed asdasays bye ys bye emy says bye eyed asda"}
	convo.Messages = append(convo.Messages, msg_1, msg_2, msg_3, msg_4, msg_5, msg_6, msg_7, msg_8, msg_9, msg_10, msg_11, msg_12, msg_13, msg_14, msg_15, msg_16, msg_17, msg_18, msg_19)

	host.MongoClient.InsertUser("jp", "1password")
	host.MongoClient.InsertConversation("jp", convo)

	red_err := host.RedisClient.LoadConversation(host.MongoClient, "jp", "MyConvo")
	if red_err != nil {
		fmt.Println(err)
	}
	host.RedisClient.AddMessageToConversation(host.MongoClient, "jp", "MyConvo", database.Message{Role: "mae", Content: "i like to eat kinder bueno"})

	host.Tick()

	defer func() {
		if err = host.MongoClient.Terminate(); err != nil {
			panic(err)
		}
		// if err = redis_svr.Terminate(); err != nil {
		// 	panic(err)
		// }
	}()
}
