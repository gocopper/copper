package cauth2

import (
	"github.com/tusharsoni/copper/cerror"
	"github.com/tusharsoni/copper/crandom"
	"golang.org/x/crypto/bcrypt"
)

func generateSessionToken() (rawToken string, encryptedToken string, err error) {
	rawToken = crandom.GenerateRandomString(128)
	encryptedTokenData, err := bcrypt.GenerateFromPassword([]byte(rawToken), bcrypt.DefaultCost)
	if err != nil {
		return "", "", cerror.New(err, "failed to generate encrypted session token", nil)
	}

	encryptedToken = string(encryptedTokenData)

	return rawToken, encryptedToken, nil
}

func verifySessionToken(encrypted, plain string) bool {
	return bcrypt.CompareHashAndPassword([]byte(encrypted), []byte(plain)) == nil
}
