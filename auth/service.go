package auth

import (
	"errors"

	"github.com/nahidhasan98/ajudge/model"
	"golang.org/x/crypto/bcrypt"
)

type authInterfacer interface {
	Authenticate(username, password string) (*model.UserData, error)
}
type auth struct {
	repoService repoInterfacer
}

func (a *auth) Authenticate(username, password string) (*model.UserData, error) {
	userData, err := a.repoService.getUserByUsername(username)
	if err != nil {
		return nil, errors.New("username not found")
	}

	err = checkPasswordHash(password, userData.Password)
	if err != nil {
		return nil, errors.New("invalid password")
	}

	return userData, nil
}

// function used above by this particular file
func checkPasswordHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func newAuthService(repo repoInterfacer) authInterfacer {
	return &auth{
		repoService: repo,
	}
}
