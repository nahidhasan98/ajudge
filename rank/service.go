package rank

type rankInterfacer interface {
	getRankList(rankType string) ([]rankModel, error)
}

type rank struct {
	repoService repoInterfacer
}

func (r *rank) getRankList(rankType string) ([]rankModel, error) {
	var rankList []rankModel
	var err error

	if rankType == "oj" {
		rankList, err = r.repoService.getOJRank()
	} else if rankType == "user" {
		rankList, err = r.repoService.getUserRank()
	}
	if err != nil {
		return nil, err
	}

	return rankList, nil
}

func newRankService(repo repoInterfacer) rankInterfacer {
	return &rank{
		repoService: repo,
	}
}
