package database

import (
	"context"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoInterface struct {
	MongoOptions *options.ClientOptions
	MongoClient  *mongo.Client
}

/*
Creates a MongoDB MongoInterface

Input:

	uri: MongoInterface connection string e.g. "mongodb://localhost:27017/"

Output:

	*mongo.Client
	Error Message
*/
func CreateMongoMongoInterface(uri string) (*MongoInterface, error) {
	MongoInterfaceAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(MongoInterfaceAPI)

	// Create a new client and connect to the MongoInterface
	client, err := mongo.Connect(context.TODO(), opts)

	MongoInterface := MongoInterface{
		MongoOptions: opts,
		MongoClient:  client,
	}

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &MongoInterface, err
}

func (s *MongoInterface) Terminate() error {
	return s.MongoClient.Disconnect(context.TODO())
}

func (s *MongoInterface) findUser(username string) *UserData {
	filter := bson.M{"username": username}
	coll := s.MongoClient.Database("UserData").Collection("Users")

	var user UserData
	err := coll.FindOne(context.TODO(), filter).Decode(&user)

	if err == nil {
		return &user
	}

	return nil
}

func (s *MongoInterface) updateUser(username string, messages []int) error {
	filter := bson.M{"username": username}
	update := bson.M{"$set": bson.M{"conversations": messages}}
	coll := s.MongoClient.Database("UserData").Collection("Users")

	result, err := coll.UpdateOne(context.TODO(), filter, update)

	if result.ModifiedCount > 1 || result.ModifiedCount == 0 {
		return fmt.Errorf(username, "may have generated duplicates of ConversationIDs: %w", os.ErrInvalid)
	}

	return err
}

func (s *MongoInterface) findConversation(username string, conversationid int) *Conversation {
	filter := bson.M{"id": conversationid}
	coll := s.MongoClient.Database("ConversationData").Collection(username)

	var convo Conversation
	err := coll.FindOne(context.TODO(), filter).Decode(&convo)
	if err == nil {
		return &convo
	}

	return nil
}

func (s *MongoInterface) updateConversation(username string, convo Conversation) error {
	filter := bson.M{"id": convo.ConversationId}
	update := bson.M{"$set": bson.M{"messages": convo.Messages}}
	coll := s.MongoClient.Database("ConversationData").Collection(username)

	result, err := coll.UpdateOne(context.TODO(), filter, update)

	if result.ModifiedCount > 1 || result.ModifiedCount == 0 {
		return fmt.Errorf(username, "may have inserted duplicates of ConversationData: %w", os.ErrInvalid)
	}

	return err
}

func (s *MongoInterface) DoesUserExist(username string) bool {
	user := s.findUser(username)
	return user != nil
}

func (s *MongoInterface) InsertUser(username string, password string) error {
	coll := s.MongoClient.Database("UserData").Collection("Users")

	data := UserData{
		Username:        username,
		Password:        password,
		ConversationIDs: make([]int, 0),
	}

	bsonUserData, err := bson.Marshal(data)
	if err != nil {
		fmt.Println(err)
		return err
	}

	user_exist := s.DoesUserExist(username)

	if user_exist == false {
		_, err = coll.InsertOne(context.TODO(), bsonUserData)
		return err
	}
	return nil
}

/*
Inserts or Updates a Pre-Existing Conversation Thread from Redis Cache Data

Also Appends ConversationId to UserData tied to Username Conversation is newly created
*/
func (s *MongoInterface) InsertConversation(username string, conversation Conversation) error {
	convo := s.findConversation(username, conversation.ConversationId)

	if convo == nil {
		bsonConvoData, err := bson.Marshal(conversation)
		if err != nil {
			return err
		}

		user := s.findUser(username)
		if user != nil {
			fmt.Println("insert convo before")
			coll := s.MongoClient.Database("ConversationData").Collection(username)
			_, err = coll.InsertOne(context.TODO(), bsonConvoData)
			s.InsertConversationId(username, conversation)
			fmt.Println("insert convoid after")
			if err != nil {
				return err
			}
		}
	} else {
		s.updateConversation(username, conversation)
	}

	return nil
}

func (s *MongoInterface) DeleteConversation(username string, conversationid int) error {
	filter := bson.M{"id": conversationid}
	convoCollection := s.MongoClient.Database("ConversationData").Collection(username)

	result, err := convoCollection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}

	s.DeleteConversationId(username, conversationid)

	if result.DeletedCount == 0 {
		fmt.Println("Warning: Unable to remove Conversation -", conversationid)
	} else if result.DeletedCount > 1 {
		fmt.Println("Error: Multiple Conversations tied to -", conversationid)
	}

	return nil
}

func (s *MongoInterface) RetrieveConversation(username string, conversationid int) (*Conversation, error) {
	filter := bson.M{"id": conversationid}
	convoCollection := s.MongoClient.Database("ConversationData").Collection(username)

	var convo Conversation
	err := convoCollection.FindOne(context.TODO(), filter).Decode(&convo)
	if err != nil {
		fmt.Println("ConversationId stored within UserData Database however Conversation was not found in ConversationData Database...")
		return nil, err
	}

	return &convo, nil
}

func (s *MongoInterface) RetrieveConversations(user UserData) ([]Conversation, error) {
	conversations := make([]Conversation, len(user.ConversationIDs))

	for i, convoid := range user.ConversationIDs {
		convo, err := s.RetrieveConversation(user.Username, convoid)
		if err != nil {
			return conversations, err
		}
		conversations[i] = *convo
	}

	return conversations, nil
}

func (s *MongoInterface) GenerateConversationId(user UserData) int {
	lastConvoId := user.ConversationIDs[len(user.ConversationIDs)-1]
	return lastConvoId + 1
}
