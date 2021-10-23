package user

import "github.com/gorilla/mux"

func Init(router *mux.Router) userInterfacer {
	repo := newRepository()
	userService := newUserService(repo)
	makeHTTPHandlers(router, userService)

	return userService
}
