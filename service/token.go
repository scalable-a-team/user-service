package service

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"time"
	"user-service/models"
)

type TokenService struct {
	SecretKey         []byte
	AccessExpireTime  time.Duration
	RefreshExpireTime time.Duration
}

func (tg *TokenService) GenerateAccessToken(user *models.User) (string, error) {
	// Create token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	// This is the information which frontend can use
	// The backend can also decode the token and get admin etc.
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = user.Username
	claims["iss"] = "login-issuer"
	claims["group"] = user.RoleGroupName
	claims["exp"] = time.Now().Add(tg.AccessExpireTime).Unix()

	// Generate encoded token and send it as response.
	// The signing string should be secret (a generated UUID works too)
	t, err := token.SignedString(tg.SecretKey)
	if err != nil {
		return "", err
	}

	return t, nil
}

func (tg *TokenService) GetUsernameFromToken(accessToken string) (string, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return tg.SecretKey, nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Get the user record from database or
		// run through your business logic to verify if the user can log in
		username := claims["username"].(string)
		return username, nil
	}
	return "", errors.New("invalid token")
}

func (tg *TokenService) GenerateRefreshToken(user *models.User) (string, error) {
	refreshToken := jwt.New(jwt.SigningMethodHS256)
	rtClaims := refreshToken.Claims.(jwt.MapClaims)
	rtClaims["username"] = user.Username
	rtClaims["iss"] = "login-issuer"
	rtClaims["group"] = user.RoleGroupName
	rtClaims["exp"] = time.Now().Add(tg.RefreshExpireTime).Unix()

	rt, err := refreshToken.SignedString(tg.SecretKey)
	if err != nil {
		return "", err
	}

	return rt, nil
}

func (tg *TokenService) RefreshAccessToken(refreshToken string) (string, error) {
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return tg.SecretKey, nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Get the user record from database or
		// run through your business logic to verify if the user can log in
		username := claims["username"].(string)
		var user models.User
		if err := user.RetrieveByUsername(username); err != nil {
			return "", err
		}
		user.RoleGroupName = claims["group"].(string)
		accessToken, err := tg.GenerateAccessToken(&user)
		if err != nil {
			return accessToken, err
		}
		return accessToken, nil
	}
	return "", errors.New("invalid token")
}
