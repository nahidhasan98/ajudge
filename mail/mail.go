package mail

func Init() mailInterfacer {
	return newMailService()
}
