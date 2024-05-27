package database

type UserData struct {
	Username        string `bson:"username"`
	Password        string `bson:"password"`
	ConversationIDs []int  `bson:"conversations"`
}
