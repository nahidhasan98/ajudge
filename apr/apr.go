package apr

import (
	"github.com/gorilla/mux"
)

func Init(router *mux.Router) aprInterfacer {
	aprService := newAprService()
	makeHTTPHandlers(router, aprService)

	return aprService
}
