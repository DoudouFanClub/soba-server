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

// Adds a conversation ID to a user's list of conversation IDs
func (s *MongoInterface) InsertConversationId(username string, conversation Conversation) error {
	filter := bson.M{"username": username}
	userCollection := s.MongoClient.Database("UserData").Collection("Users")

	var user UserData
	err := userCollection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		return fmt.Errorf("unable to insert conversation ID (user not registered): %w", err)
	}

	user.ConversationIDs = append(user.ConversationIDs, conversation.Title)
	return s.updateUser(username, user.ConversationIDs)
}

// Removes a conversation ID from a user's list of conversation IDs
func removeConversationId(user *UserData, titleToRemove string) error {
	index := -1
	for i, title := range user.ConversationIDs {
		if title == titleToRemove {
			index = i
			break
		}
	}

	if index == -1 {
		return fmt.Errorf("attempting to remove a conversation that was not found in user data")
	}

	// Remove the conversation ID by shifting elements to the left and reducing the slice size by 1.
	copy(user.ConversationIDs[index:], user.ConversationIDs[index+1:])
	user.ConversationIDs = user.ConversationIDs[:len(user.ConversationIDs)-1]
	return nil
}

// Removes a conversation ID from a user's list of conversation IDs in MongoDB
func (s *MongoInterface) DeleteConversationId(username string, title string) error {
	filter := bson.M{"username": username}
	userCollection := s.MongoClient.Database("UserData").Collection("Users")

	var user UserData
	err := userCollection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		return fmt.Errorf("failed to find user: %w", err)
	}

	err = removeConversationId(&user, title)
	if err != nil {
		return err
	}

	return nil
}
