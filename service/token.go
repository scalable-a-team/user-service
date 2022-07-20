package service

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"time"
)

type TokenService struct {
	SecretKey         []byte
	ISS               string
	AccessExpireTime  time.Duration
	RefreshExpireTime time.Duration
}

type TokenUserInput struct {
	Username      string
	UserID        uuid.UUID
	RoleGroupName string
	Firstname     string
	Lastname      string
}

func (tg *TokenService) GenerateAccessToken(user *TokenUserInput) (string, error) {
	// Create token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	// This is the information which frontend can use
	// The backend can also decode the token and get admin etc.
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = user.Username
	claims["userid"] = user.UserID
	claims["iss"] = tg.ISS
	claims["group"] = user.RoleGroupName
	claims["firstname"] = user.Firstname
	claims["lastname"] = user.Lastname
	claims["exp"] = time.Now().Add(tg.AccessExpireTime).Unix()

	// Generate encoded token and send it as response.
	// The signing string should be secret (a generated UUID works too)
	fmt.Println("generating token")
	fmt.Println(string(tg.SecretKey))
	t, err := token.SignedString(tg.SecretKey)
	if err != nil {
		return "", err
	}

	return t, nil
}

func (tg *TokenService) GetUserIDFromToken(accessToken string) (uuid.UUID, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return tg.SecretKey, nil
	})
	if err != nil {
		return uuid.Nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Get the user record from database or
		// run through your business logic to verify if the user can log in
		userID, err := uuid.Parse(claims["userid"].(string))
		if err != nil {
			return uuid.Nil, err
		}
		return userID, nil
	}
	return uuid.Nil, errors.New("invalid token")
}

func (tg *TokenService) GenerateRefreshToken(user *TokenUserInput) (string, error) {
	refreshToken := jwt.New(jwt.SigningMethodHS256)
	rtClaims := refreshToken.Claims.(jwt.MapClaims)
	rtClaims["username"] = user.Username
	rtClaims["userid"] = user.UserID
	rtClaims["iss"] = tg.ISS
	rtClaims["group"] = user.RoleGroupName
	rtClaims["exp"] = time.Now().Add(tg.RefreshExpireTime).Unix()

	rt, err := refreshToken.SignedString(append(tg.SecretKey, []byte("refresh")...))
	if err != nil {
		return "", err
	}

	return rt, nil
}

func (tg *TokenService) ValidateRefreshAccessToken(refreshToken string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return append(tg.SecretKey, []byte("refresh")...), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Get the user record from database or
		// run through your business logic to verify if the user can log in
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
