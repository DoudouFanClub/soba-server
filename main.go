package main

import (
	"fmt"
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

	host.RedisClient.ClearRedisMongoInterface()

	host.Tick()

	defer func() {
		if err = host.MongoClient.Terminate(); err != nil {
			panic(err)
		}
	}()
}
