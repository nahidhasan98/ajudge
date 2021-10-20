package auth

//User variable for holding a single user details
type User struct {
	UserID   int    `bson:"userID"`
	Username string `bson:"username"`
	Password string `bson:"password"`
}
