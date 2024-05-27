package database

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Server struct {
	ClientOptions *options.ClientOptions
	Client        *mongo.Client
}

/*
Creates a MongoDB Server

Input:

	uri: Server connection string e.g. "mongodb://localhost:27017/"

Output:

	*mongo.Client
	Error Message
*/
func CreateMongoServer(uri string) (*mongo.Client, error) {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), opts)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return client, err
}

func DoesUserExist(client *mongo.Client, username string) bool {
	filter := bson.M{"username": username}
	coll := client.Database("UserData").Collection("Users")

	var user UserData
	err := coll.FindOne(context.TODO(), filter).Decode(&user)

	if err != nil {
		fmt.Println(err)
		return false
	}

	return true
}

func InsertUser(client *mongo.Client, username string, password string) error {
	coll := client.Database("UserData").Collection("Users")

	data := UserData{
		Username:        username,
		Password:        password,
		ConversationIDs: make([]int, 0),
	}

	bson_user_data, err := bson.Marshal(data)

	if err != nil {
		fmt.Println(err)
		return err
	}

	coll.InsertOne(context.TODO(), bson_user_data)

	return nil
}

/*
	add update user details
	add does conversation exist
	add update conversation
*/

func InsertConversation(client *mongo.Client, username string, conversation Conversation) error {
	filter := bson.M{"id": conversation.ConversationId}
	convo_collection := client.Database("ConversationData").Collection("Conversations")

	var convo Conversation
	err := convo_collection.FindOne(context.TODO(), filter).Decode(&convo)

	// If err == nil - Does not exist yet, insert
	// Else - Overwrite
	if err != nil {
		bson_convo_data, b_err := bson.Marshal(conversation)

		if b_err != nil {
			return b_err
		}

		convo_collection.InsertOne(context.TODO(), bson_convo_data)
	} else {
		convo = conversation
	}

	return nil
}

func DeleteConversation(client *mongo.Client, conversationid int) error {
	filter := bson.M{"id": conversationid}
	convo_collection := client.Database("ConversationData").Collection("Conversations")

	result, err := convo_collection.DeleteOne(context.TODO(), filter)

	if err != nil {
		return err
	} else {
		if result.DeletedCount == 0 {
			fmt.Println("Warning: Unable to remove Conversation -", conversationid)
		} else if result.DeletedCount > 1 {
			fmt.Println("Error: Multiple Conversations tied to -", conversationid)
		}
	}

	return nil
}

func retrieveConversation(client *mongo.Client, conversationid int) (*Conversation, error) { // double check this return when null
	filter := bson.M{"id": conversationid}
	convo_collection := client.Database("ConversationData").Collection("Conversations")

	var convo Conversation
	err := convo_collection.FindOne(context.TODO(), filter).Decode(&convo)

	if err != nil {
		fmt.Println("ConversationId stored within UserData Database however Conversation was not found in ConversationData Database...")
		return nil, err
	}

	return &convo, err
}

func RetrieveConversations(client *mongo.Client, userdata UserData) ([]Conversation, error) {
	conversations := make([]Conversation, len(userdata.ConversationIDs))

	// Iterate all User Conversation Ids and retrieve their Conversations
	for i, convoid := range userdata.ConversationIDs {
		convo, err := retrieveConversation(client, convoid)

		// Valid Conversation Found
		if err != nil {
			return conversations, err
		} else {
			conversations[i] = *convo
		}

	}

	return conversations, nil
}

/*
	- query db for all conversations for a username
	- generate conversation id from "topic"
	  - username_1, username_2, ...
	-
*/
