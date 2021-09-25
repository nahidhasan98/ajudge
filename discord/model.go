package discord

type discordModel struct {
	SubID            int    `bson:"subID"`
	MessageID        string `bson:"messageID"`
	Message          string `bson:"message"`
	NotificationType string `bson:"notificationType"`
}
