package middlewares

import (
	"os"
	"time"
	"user-service/service"
)

var customerTokenService *service.TokenService
var sellerTokenService *service.TokenService

func InitCustomerJWTMiddleware() *service.TokenService {
	secretKey := os.Getenv("JWT_CUSTOMER_SECRET_KEY")

	customerTokenService = &service.TokenService{
		SecretKey:         []byte(secretKey),
		AccessExpireTime:  time.Minute * 15,
		RefreshExpireTime: time.Hour * 5,
	}
	return customerTokenService
}

func GetCustomerJwtMiddleware() *service.TokenService {
	return customerTokenService
}

func InitSellerJWTMiddleware() *service.TokenService {
	secretKey := os.Getenv("JWT_SELLER_SECRET_KEY")

	sellerTokenService = &service.TokenService{
		SecretKey:         []byte(secretKey),
		AccessExpireTime:  time.Minute * 15,
		RefreshExpireTime: time.Hour * 5,
	}
	return sellerTokenService
}

func GetSellerJwtMiddleware() *service.TokenService {
	return sellerTokenService
}
