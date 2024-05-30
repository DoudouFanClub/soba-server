package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"llm_server/database"

	"github.com/redis/go-redis/v9"
)

const (
	redisKeyFormat = "%s_%s"
)

type RedisCache struct {
	client *redis.Client
}

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

func (r *RedisCache) Terminate() error {
	return r.ClearRedisMongoInterface()
}

func (r *RedisCache) ClearRedisMongoInterface() error {
	return r.client.FlushDB(context.TODO()).Err()
}

func (r *RedisCache) LoadConversationBytes(key string, value *[]database.Message) error {
	// Marshal conversation messages to JSON
	messageArrJSONBytes, err := json.Marshal(value)

	if err != nil {
		fmt.Println(err)
		return err
	} else {
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

	return r.LoadConversationBytes(key, &convo.Messages)
}

func (r *RedisCache) UnloadConversation(username string, conversationId int) error {
	key := fmt.Sprintf(redisKeyFormat, username, conversationId)
	return r.client.Del(context.TODO(), key).Err()
}

func (r *RedisCache) AddMessageToConversation(mongoMongoInterface *database.MongoInterface, username string, title string, newMsg database.Message) error {
	key := fmt.Sprintf(redisKeyFormat, username, title)
	msg, err := r.GetDataMsgArray(mongoMongoInterface, username, title)

	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("unable to retrieve conversation: %w", err)
	}
	*msg = append(*msg, newMsg)

	return r.LoadConversationBytes(key, msg)
}

func (r *RedisCache) RemoveMessageFromConversation(mongoMongoInterface *database.MongoInterface, username string, title string, newMsg database.Message) error {
	return nil
}

// Note that return type is a Pointer type
// If memory usage ends up being an issue, we may want to return by copy rather than reference
func (r *RedisCache) GetDataMsgArray(mongoMongoInterface *database.MongoInterface, username string, title string) (*[]database.Message, error) {
	key := fmt.Sprintf(redisKeyFormat, username, title)
	val, err := r.client.Get(context.TODO(), key).Bytes()

	if err != nil {
		return nil, fmt.Errorf("unable to retrieve conversation: %w", err)
	}
	var convo []database.Message
	err = json.Unmarshal(val, &convo)
	return &convo, err
}

// Note that return type is a Pointer type
// If memory usage ends up being an issue, we may want to return by copy rather than reference
func (r *RedisCache) GetDataConversation(mongoMongoInterface *database.MongoInterface, username string, title string) (*database.Conversation, error) {
	convo, err := r.GetDataMsgArray(mongoMongoInterface, username, title)

	if err != nil {
		return nil, err
	}
	return &database.Conversation{
		Title:    title,
		Messages: *convo,
	}, nil
}
