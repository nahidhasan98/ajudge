package apr

import (
	"embed"

	"github.com/gorilla/mux"
)

func Init(router *mux.Router, frontendFS embed.FS) aprInterfacer {
	aprService := newAprService()
	makeHTTPHandlers(router, aprService, frontendFS)

	return aprService
}
