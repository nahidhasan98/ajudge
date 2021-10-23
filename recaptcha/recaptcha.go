package recaptcha

func Init() recaptchaInterfacer {
	return newRecaptchaService()
}
