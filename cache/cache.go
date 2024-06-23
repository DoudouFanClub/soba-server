package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"llm_server/database"

	"github.com/redis/go-redis/v9"
)

// Redis format used to generate conversation "Key"
// Username_Title
const (
	redisKeyFormat = "%s_%s"
)

// Redis cache to store conversation dat
type RedisCache struct {
	client *redis.Client
}

// Initializes a Redis cache connection
func CreateRedisInterface(addr string, password string, db int) (*RedisCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	if err := client.Ping(context.TODO()).Err(); err != nil {
		return nil, fmt.Errorf("unable to connect to Redis: %w", err)
	}

	return &RedisCache{client: client}, nil
}

// To be called when the Backend Server terminates
func (r *RedisCache) Terminate() error {
	return r.ClearRedisMongoInterface()
}

// Clears the entire Redis cache
func (r *RedisCache) ClearRedisMongoInterface() error {
	return r.client.FlushDB(context.TODO()).Err()
}

// Loads conversation messages into the Redis cache
func (r *RedisCache) LoadConversationData(key string, value *[]database.Message) error {
	// Marshal conversation messages to JSON
	messageArrJSONBytes, err := json.Marshal(value)

	if err != nil {
		fmt.Println(err)
		return err
	} else {
		fmt.Println("adding a message to a conversation")
		return r.client.Set(context.TODO(), key, messageArrJSONBytes, 0).Err()
	}
}

// Retrieve Conversation from MongoDB and load it into the redis cache
func (r *RedisCache) LoadConversation(mongoInterface *database.MongoInterface, username string, title string) error {
	key := fmt.Sprintf(redisKeyFormat, username, title)

	// Retrieve conversation from MongoDB
	convo, err := mongoInterface.RetrieveConversation(username, title)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("ERROR - Unable to load Conversation: %w", err)
	}

	return r.LoadConversationData(key, &convo.Messages)
}

// Removes a conversation from the Redis cache
func (r *RedisCache) UnloadConversation(username string, title string) error {
	key := fmt.Sprintf(redisKeyFormat, username, title)
	return r.client.Del(context.TODO(), key).Err()
}

// Appends a new message to the Redis cache
func (r *RedisCache) AddMessageToConversation(mongoInterface *database.MongoInterface, username string, title string, newMsg database.Message) error {
	key := fmt.Sprintf(redisKeyFormat, username, title)
	msg, err := r.GetDataMsgArray(mongoInterface, username, title)

	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("unable to retrieve conversation: %w", err)
	}
	*msg = append(*msg, newMsg)

	return r.LoadConversationData(key, msg)
}

// Retrieves conversation data from the Redis cache
func (r *RedisCache) GetDataMsgArray(mongoInterface *database.MongoInterface, username string, title string) (*[]database.Message, error) {
	key := fmt.Sprintf(redisKeyFormat, username, title)
	val, err := r.client.Get(context.TODO(), key).Bytes()

	if err != nil {
		fmt.Println("Key:", key)
		fmt.Println("Value:", val)

		return nil, fmt.Errorf("unable to retrieve conversation: %w", err)
	}
	var convo []database.Message
	err = json.Unmarshal(val, &convo)
	return &convo, err
}

// Retrieves the entire conversation from the Redis cache
func (r *RedisCache) GetDataConversation(mongoInterface *database.MongoInterface, username string, title string) (*database.Conversation, error) {
	convo, err := r.GetDataMsgArray(mongoInterface, username, title)

	if err != nil {
		return nil, err
	}	
	return &database.Conversation{
		Title:    title,
		Messages: *convo,
	}, nil
}
