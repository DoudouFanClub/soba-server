package database

type Conversation struct {
	ConversationId int       `bson:"id"`
	Messages       []Message `bson:"messages"`
}

func SetConversationId(conversation *Conversation, conversationid int) {
	conversation.ConversationId = conversationid
}

func InsertMessage(conversation *Conversation, msg Message) {
	conversation.Messages = append(conversation.Messages, msg)
}
