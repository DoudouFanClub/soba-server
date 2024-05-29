package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"llm_server/database"

	"github.com/redis/go-redis/v9"
)

// struct definitions
// contains redis cache

// (reference ptr to redis cache) convert cache into conversation

// (reference ptr to redis cache) add message under a cached conversation

// (reference ptr to redis cache) load conversation to cache

// can consider:
// - update message
//   - edit existing user message and POP ai response and re-query
// - delete user prompt - reset history for response and user prompt

// type RedisCache struct {
// 	RedisOption *redis.Options
// 	RedisClient *redis.Client
// }

// func CreateRedisMongoInterface(addr string, networkType string, password string, db int) *RedisCache {
// 	opts := redis.Options{
// 		Addr:     addr,
// 		Network:  networkType,
// 		Password: password,
// 		DB:       db,
// 	}

// 	rbd := redis.NewClient(&opts)

// 	return &RedisCache{
// 		RedisOption: &opts,
// 		RedisClient: rbd,
// 	}
// }

// func (r *RedisCache) Terminate() error {
// 	return r.ClearRedisMongoInterface()
// }

// func (r *RedisCache) ClearRedisMongoInterface() error {
// 	err := r.RedisClient.FlushDB(context.TODO()).Err()

// 	if err != nil {
// 		return fmt.Errorf("ERROR - Unable to Flush Database for Redis Cache: %w", err)
// 	}

// 	return nil
// }

// func (r *RedisCache) LoadConversationBytes(mongoMongoInterface *database.MongoInterface, key string, value []database.Message) error {
// 	// Marshal conversation messages to JSON
// 	messageArrJSONBytes, err := json.Marshal(value)
// 	if err != nil {
// 		return fmt.Errorf("ERROR - Unable to Marshal Conversation from MongoDB - ID %s: %w", key, err)
// 	}

// 	// Store the JSON in Redis
// 	err = r.RedisClient.Set(context.TODO(), key, messageArrJSONBytes, 0).Err() // eventually replace id with a string equivalent
// 	if err != nil {
// 		return fmt.Errorf("ERROR - Unable to load Message Array for [%s's] Conversation onto the Redis Cache: %w", key, err)
// 	}

// 	return nil
// }

const (
	redisKeyFormat = "%s_%d"
)

type RedisCache struct {
	client *redis.Client
}

func CreateRedisMongoInterface(addr string, password string, db int) (*RedisCache, error) {
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
func (r *RedisCache) LoadConversation(mongoMongoInterface *database.MongoInterface, username string, conversationId int) error {
	key := fmt.Sprintf(redisKeyFormat, username, conversationId)

	// Retrieve conversation from MongoDB
	convo, err := mongoMongoInterface.RetrieveConversation(username, conversationId)
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

func (r *RedisCache) AddMessageToConversation(mongoMongoInterface *database.MongoInterface, username string, conversationId int, newMsg database.Message) error {
	key := fmt.Sprintf(redisKeyFormat, username, conversationId)
	msg, err := r.GetDataMsgArray(mongoMongoInterface, username, conversationId)

	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("unable to retrieve conversation: %w", err)
	}
	*msg = append(*msg, newMsg)

	return r.LoadConversationBytes(key, msg)
}

func (r *RedisCache) RemoveMessageFromConversation(mongoMongoInterface *database.MongoInterface, username string, conversationId int, newMsg database.Message) error {
	return nil
}

// Note that return type is a Pointer type
// If memory usage ends up being an issue, we may want to return by copy rather than reference
func (r *RedisCache) GetDataMsgArray(mongoMongoInterface *database.MongoInterface, username string, conversationId int) (*[]database.Message, error) {
	key := fmt.Sprintf(redisKeyFormat, username, conversationId)
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
func (r *RedisCache) GetDataConversation(mongoMongoInterface *database.MongoInterface, username string, conversationId int) (*database.Conversation, error) {
	convo, err := r.GetDataMsgArray(mongoMongoInterface, username, conversationId)

	if err != nil {
		return nil, err
	}
	return &database.Conversation{
		ConversationId: conversationId,
		Messages:       *convo,
	}, nil
}
