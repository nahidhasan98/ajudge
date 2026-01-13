package user

import (
	"embed"

	"github.com/gorilla/mux"
)

func Init(router *mux.Router, frontendFS embed.FS) userInterfacer {
	repo := newRepository()
	userService := newUserService(repo, frontendFS)
	makeHTTPHandlers(router, userService)

	return userService
}
