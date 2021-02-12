package vjudge

import (
	"net/http"

	"github.com/nahidhasan98/ajudge/errorhandling"
	"github.com/nahidhasan98/ajudge/model"
)

//GETRequest function for http GET request
func GETRequest(apiURL string) *http.Response {
	defer errorhandling.Recovery() //for panic() error Recovery

	req, err := http.NewRequest("GET", apiURL, nil)
	errorhandling.Check(err)
	req.Header.Add("Content-Type", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")

	response, err := model.Client.Do(req)
	errorhandling.Check(err)

	return response
}
