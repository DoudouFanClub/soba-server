package database

type Message struct {
	Role    string `json:"role" bson:"role"`
	Content string `json:"content" bson:"content"`
}

type MessagePrompt struct {
	Username string  `bson:"username"`
	Title    string  `bson:"title"`
	Contents Message `bson:"content"`
}

type FrontendMessagesPrompt struct {
	Username string    `json:"username" bson:"username"`
	Title    string    `json:"title" bson:"title"`
	Contents []Message `json:"contents" bson:"contents"`
}