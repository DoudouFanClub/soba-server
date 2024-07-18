package database

type ConversationRequest struct {
	Username  string `json:"username" bson:"username"`
	Title     string `json:"title" bson:"title"`
	PrevTitle string `json:"prevtitle" bson:"prevtitle"`
	// need to add a new field to specify the model
}

type ConversationTitlesRequest struct {
	Username string   `json:"username" bson:"username"`
	Titles   []string `json:"titles" bson:"titles"`
}

type Conversation struct {
	Title    string    `json:"title" bson:"title"`
	Messages []Message `json:"messages" bson:"messages"`
}

type RemoveConversation struct {
	Username string `json:"username" bson:"username"`
	Title    string `json:"title" bson:"title"`
}

func SetConversationId(conversation *Conversation, title string) {
	conversation.Title = title
}

func InsertMessage(conversation *Conversation, msg Message) {
	conversation.Messages = append(conversation.Messages, msg)
}
