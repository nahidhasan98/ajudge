package discord

type discordModel struct {
	SubID            int    `bson:"subID"`
	MessageID        string `bson:"messageID"`
	Message          string `bson:"message"`
	NotificationType string `bson:"notificationType"`
}

// User variable for holding a single user details
type UserModel struct {
	UserID               int    ``
	FullName             string ``
	Email                string ``
	Username             string ``
	Password             string ``
	CreatedAt            int64  ``
	IsVerified           bool   ``
	AccVerifyToken       string ``
	AccVerifyTokenSentAt int64  ``
	PassResetToken       string ``
	PassResetTokenSentAt int64  ``
	TotalSolved          int    ``
}
