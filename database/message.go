package database

type Message struct {
	Role    string `bson:"role"`
	Content string `bson:"content"`
}

type MessagePrompt struct {
	Username string `bson:"username"`
	Title    string `bson:"title"`
	Content  string `bson:"content"`
}

type FrontendTest struct {
	Text string `bson:"text"`
}