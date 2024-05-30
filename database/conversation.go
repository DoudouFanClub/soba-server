package database

type ConversationRequest struct {
	Username string `bson:"username"`
	Title    string `bson:"title"`
}

type Conversation struct {
	Title    string    `bson:"title"`
	Messages []Message `bson:"messages"`
}

func SetConversationId(conversation *Conversation, title string) {
	conversation.Title = title
}

func InsertMessage(conversation *Conversation, msg Message) {
	conversation.Messages = append(conversation.Messages, msg)
}
