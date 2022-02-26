package rank

import "github.com/gorilla/mux"

func Init(router *mux.Router) rankInterfacer {
	repo := newRepository()
	rankService := newRankService(repo)
	makeHTTPHandlers(router, rankService)

	return rankService
}
