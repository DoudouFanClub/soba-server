package main

import (
	"fmt"
	"llm_server/database"

	"context"

	"go.mongodb.org/mongo-driver/bson"
)

const uri = "mongodb://localhost:27017/"

func main() {

	client, err := database.CreateMongoServer(uri)

	status := database.DoesUserExist(client, "jeremy")
	status2 := database.DoesUserExist(client, "ryan")

	if status {
		fmt.Println("Jeremy was here")
	} else {
		fmt.Println("Uh oh he ded")
	}

	if status2 {
		fmt.Println("Ryan was here")
		database.InsertUser(client, "ryan", "byoray wew")
	} else {
		fmt.Println("Uh oh he ded too")
	}

	if err != nil {
		panic(err)
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
	coll := client.Database("ace").Collection("ttace")

	coll.InsertOne(context.Background(), bson.D{{"name", "ACE"}})

}

/*
package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

type Person struct {
	Name  string
	Email string
}

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	person := Person{
		Name:  "Ace",
		Email: "aceacett@simp.com",
	}

	json_bytes, err := json.Marshal(person)
	if err != nil {
		fmt.Println("noob")
		return
	}

	err2 := rdb.Set(ctx, "key", json_bytes, 0).Err()
	if err2 != nil {
		panic(err)
	}

	val, err := rdb.Get(ctx, "key").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("key", val)

	val2, err := rdb.Get(ctx, "key2").Result()
	if err == redis.Nil {
		fmt.Println("key2 does not exist")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println("key2", val2)
	}
	// Output: key value
	// key2 does not exist
}
*/
