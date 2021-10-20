package auth

import "github.com/gorilla/mux"

func Init(router *mux.Router) authInterfacer {
	repo := newRepository()
	authService := newAuthService(repo)
	makeHTTPHandlers(router, authService)

	return authService
}
