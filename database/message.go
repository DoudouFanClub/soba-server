package database

type Message struct {
	Role    string `bson:"role"`
	Content string `bson:"content"`
}
