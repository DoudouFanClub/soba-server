package database

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)

type UserData struct {
	Username        string   `bson:"username"`
	Password        string   `bson:"password"`
	ConversationIDs []string `bson:"conversations"`
}

func (s *MongoInterface) InsertConversationId(username string, conversation Conversation) {
	filter := bson.M{"username": username}
	user_collection := s.MongoClient.Database("UserData").Collection("Users")

	var user UserData
	err := user_collection.FindOne(context.TODO(), filter).Decode(&user)

	if err != nil {
		fmt.Println("Error: Unable to Insert Conversation Id (User Not Registered):", err)
	} else {
		user.ConversationIDs = append(user.ConversationIDs, conversation.Title)
		s.updateUser(username, user.ConversationIDs)
	}
}

func removeConversationId(user *UserData, title_to_remove string) {
	index := -1
	for i, title := range user.ConversationIDs {
		if title == title_to_remove {
			index = i
			break
		}
	}

	if index == -1 {
		fmt.Println("Error: Attempting to remove a Conversation that was NOT within UserData")
		return
	}

	// Shift Elements to the Left
	copy(user.ConversationIDs[index:], user.ConversationIDs[index+1:])

	// Reduce size by 1
	user.ConversationIDs = user.ConversationIDs[:len(user.ConversationIDs)-1]
}

func (s *MongoInterface) DeleteConversationId(username string, title string) {
	filter := bson.M{"username": username}
	user_collection := s.MongoClient.Database("UserData").Collection("Users")

	var user UserData
	err := user_collection.FindOne(context.TODO(), filter).Decode(&user)

	if err != nil {
		fmt.Println(err)
	} else {
		removeConversationId(&user, title)
	}
}
