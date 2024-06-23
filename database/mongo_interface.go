package database

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoInterface struct {
	MongoClient  *mongo.Client
}

// Creates a new MongoDB client.
func CreateMongoInterface(uri string) (*MongoInterface, error) {
	clientOptions := options.Client().ApplyURI(uri)

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	return &MongoInterface{
		MongoClient: client,
	}, nil
}

// Disconnects the MongoDB client
func (s *MongoInterface) Terminate() error {
	return s.MongoClient.Disconnect(context.TODO())
}

// Retrieves a user from MongoDB by username
func (s *MongoInterface) findUser(username string) (*UserData, error) {
	filter := bson.M{"username": username}
	coll := s.MongoClient.Database("UserData").Collection("Users")

	var user UserData
	err := coll.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	return &user, nil
}

// Updates user messages in MongoDB
func (s *MongoInterface) updateUser(username string, messages []string) error {
	filter := bson.M{"username": username}
	update := bson.M{"$set": bson.M{"conversations": messages}}
	coll := s.MongoClient.Database("UserData").Collection("Users")

	result, err := coll.UpdateOne(context.TODO(), filter, update)

	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	if result.ModifiedCount != 1 {
		return fmt.Errorf("expected 1 document to be modified, got %d", result.ModifiedCount)
	}

	return nil
}

// findConversation retrieves a conversation from MongoDB by title under a user's collection.
func (s *MongoInterface) findConversation(username string, title string) (*Conversation, error) {
	filter := bson.M{"title": title}
	coll := s.MongoClient.Database("ConversationData").Collection(username)

	var convo Conversation
	err := coll.FindOne(context.TODO(), filter).Decode(&convo)
	if err != nil {
		return nil, fmt.Errorf("failed to find conversation: %w", err)
	}

	return &convo, nil
}

// Updates a conversation in MongoDB
func (s *MongoInterface) updateConversation(username string, convo Conversation) error {
	filter := bson.M{"title": convo.Title}
	update := bson.M{"$set": bson.M{"messages": convo.Messages}}
	coll := s.MongoClient.Database("ConversationData").Collection(username)

	result, err := coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return fmt.Errorf("failed to update conversation: %w", err)
	}

	if result.ModifiedCount != 1 {
		return fmt.Errorf("expected 1 document to be modified, got %d", result.ModifiedCount)
	}

	return nil
}

// Checks if a user exists in MongoDB
func (s *MongoInterface) DoesUserExist(username string) bool {
	user, err := s.findUser(username)
	if err != nil {
		return false
	}
	return user != nil
}

// Checks if a conversation topic exists
func (s *MongoInterface) DoesConvoExist(username string, title string) bool {
	convo, err := s.findConversation(username, title)
	if err != nil {
		return false
	}
	return convo != nil
}

// Retrieves a conversation from Username
func (s *MongoInterface) GetConvo(username string, title string) Conversation {
	convo, err := s.findConversation(username, title)
	if err != nil {
		return Conversation{}
	}
	return *convo
}

// Verifies whether the User's login is valid
func (s *MongoInterface) IsUserLoginValid(userInput UserData) bool {
	user, err := s.findUser(userInput.Username)
	if err != nil || user == nil {
		return false
	}
	return user.Password == userInput.Password
}

// Inserts a new User into MongoDB upon successful registration
func (s *MongoInterface) InsertUser(username string, password string) error {
	coll := s.MongoClient.Database("UserData").Collection("Users")

	data := UserData{
		Username:        username,
		Password:        password,
		ConversationIDs: make([]string, 0),
	}

	bsonUserData, err := bson.Marshal(data)
	if err != nil {
		fmt.Println(err)
		return err
	}

	user_exist := s.DoesUserExist(username)

	if !user_exist {
		_, err = coll.InsertOne(context.TODO(), bsonUserData)
		return err
	}
	return nil
}

// Inserts or updates an existing conversation in MongoDB
func (s *MongoInterface) InsertConversation(username string, conversation Conversation) error {
	convo, _ := s.findConversation(username, conversation.Title)

	// Just insert if it doesn't yet exist
	// Else, update it
	if convo == nil {
		bsonConvoData, err := bson.Marshal(conversation)
		if err != nil {
			return err
		}

		user, _ := s.findUser(username)
		if user != nil {
			coll := s.MongoClient.Database("ConversationData").Collection(username)
			_, err = coll.InsertOne(context.TODO(), bsonConvoData)
			if err != nil {
				return fmt.Errorf("unable to insert new conversation to user in ConversationData: %w", err)
			}
			err = s.InsertConversationId(username, conversation)
			if err != nil {
				return fmt.Errorf("unable to insert new conversation to user in UserData: %w", err)
			}
		}
	} else {
		s.updateConversation(username, conversation)
	}

	return nil
}

// Deletes a conversation from MongoDB
func (s *MongoInterface) DeleteConversation(username string, title string) error {
	filter := bson.M{"title": title}
	convoCollection := s.MongoClient.Database("ConversationData").Collection(username)

	// Removes the Conversation from Username
	result, err := convoCollection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return fmt.Errorf("failed to delete conversation: %w", err)
	}
	
	// Removes the TitleId attached to the UserData
	err = s.DeleteConversationId(username, title)
	if err != nil {
		return fmt.Errorf("unable to remove conversation id from user %w", err)
	}

	if result.DeletedCount != 1 {
		return fmt.Errorf("expected 1 document to be deleted, got %d", result.DeletedCount)
	}

	return nil
}

// Retrieves a Conversation from a User based on the Title
func (s *MongoInterface) RetrieveConversation(username string, title string) (*Conversation, error) {
	filter := bson.M{"title": title}
	convoCollection := s.MongoClient.Database("ConversationData").Collection(username)

	var convo Conversation
	err := convoCollection.FindOne(context.TODO(), filter).Decode(&convo)
	if err != nil {
		fmt.Println("ConversationId stored within UserData Database however Conversation was not found in ConversationData Database...")
		return nil, err
	}

	return &convo, nil
}

// Retrieves all the Conversations topics attached to a User
func (s *MongoInterface) RetrieveConversations(user UserData) ([]Conversation, error) {
	conversations := make([]Conversation, len(user.ConversationIDs))

	for i, title := range user.ConversationIDs {
		convo, err := s.RetrieveConversation(user.Username, title)
		if err != nil {
			return conversations, fmt.Errorf("unable to find User when retrieving conversation titles: %w", err)
		}
		conversations[i] = *convo
	}

	return conversations, nil
}

// Retrieves all the Conversations titles attached to a User
func (s *MongoInterface) RetrieveConversationTitles(userName string) []string {
	user, err := s.findUser(userName)
	if err != nil {
		fmt.Println("unable to find User when retrieving conversation titles:", err)
		return make([]string, 0)
	}
	return user.ConversationIDs
}