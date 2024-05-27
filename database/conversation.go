package database

type Conversation struct {
	ConversationId int       `bson:"id"`
	Messages       []Message `bson:"messages"`
}

func InsertMessage(conversation *Conversation, msg Message) {
	conversation.Messages = append(conversation.Messages, msg)
}
