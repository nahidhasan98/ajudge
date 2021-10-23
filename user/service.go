package user

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/nahidhasan98/ajudge/errorhandling"
	"github.com/nahidhasan98/ajudge/mail"
	"golang.org/x/crypto/bcrypt"
)

type userInterfacer interface {
	isAvailableEmail(email string) bool
	isAvailableUsername(username string) bool
	AddUser(fullName, email, username, password string) (*UserModel, error)
}
type user struct {
	repoService repoInterfacer
}

func (u *user) AddUser(fullName, email, username, password string) (*UserModel, error) {
	password = makePasswordHash(password) // hashing password
	newToken := generateToken()           // generating token for account verification
	currTime := time.Now().Unix()
	lastUserID := u.repoService.getLastUserID()

	// preparing data for inserting to DB
	userData := &UserModel{
		UserID:               lastUserID + 1,
		FullName:             fullName,
		Email:                email,
		Username:             username,
		Password:             password,
		CreatedAt:            currTime,
		IsVerified:           false,
		AccVerifyToken:       newToken,
		AccVerifyTokenSentAt: currTime,
		PassResetToken:       "",
		PassResetTokenSentAt: 0,
	}

	err := u.repoService.createUser(userData)
	if err != nil {
		return nil, err
	}

	err = u.repoService.updateLastUserID()
	errorhandling.Check(err)

	// sending mail to user email with a verification link
	verificationLink := "https://ajudge.net/verify-email/token=" + newToken
	mail := mail.Init()
	mail.SendMailForRegistration(email, username, verificationLink)

	return userData, nil
}

func (u *user) isAvailableEmail(email string) bool {
	_, err := u.repoService.getUserByEmail(email)
	return err != nil
}

func (u *user) isAvailableUsername(username string) bool {
	_, err := u.repoService.getUserByUsername(username)
	return err != nil
}

// makePasswordHash function
func makePasswordHash(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	errorhandling.Check(err)
	return string(bytes)
}

// generateToken function
func generateToken() string {
	b := make([]byte, 16)
	rand.Read(b)

	hasher := md5.New()
	hasher.Write(b)
	return hex.EncodeToString(hasher.Sum(nil))
}

func newUserService(repo repoInterfacer) userInterfacer {
	return &user{
		repoService: repo,
	}
}
