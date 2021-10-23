package user

// User variable for holding a single user details
type UserModel struct {
	UserID               int    `bson:"userID"`
	FullName             string `bson:"fullName"`
	Email                string `bson:"email"`
	Username             string `bson:"username"`
	Password             string `bson:"password"`
	CreatedAt            int64  `bson:"createdAt"`
	IsVerified           bool   `bson:"isVerified"`
	AccVerifyToken       string `bson:"accVerifyToken"`
	AccVerifyTokenSentAt int64  `bson:"accVerifyTokenSentAt"`
	PassResetToken       string `bson:"passResetToken"`
	PassResetTokenSentAt int64  `bson:"passResetTokenSentAt"`
	TotalSolved          int    `bson:"totalSolved"`
}

// LastUsedID variable for holding the last ID used for user registration, problem submission & contest creation
type lastUsedIDModel struct {
	LastUserID       int `bson:"lastUserID"`
	LastSubmissionID int `bson:"lastSubmissionID"`
	LastContestID    int `bson:"lastContestID"`
}
