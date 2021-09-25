package discord

func Init() discordInterfacer {
	discordRepo := newRepository()
	discordService := newDiscordService(discordRepo)

	return discordService
}
