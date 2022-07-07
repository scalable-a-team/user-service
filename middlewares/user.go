package middlewares

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"os"
	"time"
	"user-service/enums"
	"user-service/forms"
	"user-service/models"
)

var identityKey = os.Getenv("JWT_IDENTITY_KEY")

func PayloadFunc(data interface{}) jwt.MapClaims {
	if v, ok := data.(*models.User); ok {
		return jwt.MapClaims{
			identityKey: v.Username,
			"group":     v.Group.Name,
			"iss":       "login-issuer",
		}
	}
	return jwt.MapClaims{}
}

func CustomerAuthMiddleware() (*jwt.GinJWTMiddleware, error) {
	secretKey := os.Getenv("JWT_CUSTOMER_SECRET_KEY")

	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "jwt",
		Key:         []byte(secretKey),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		IdentityKey: identityKey,
		PayloadFunc: PayloadFunc,
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginData forms.UserSignIn
			if err := c.ShouldBind(&loginData); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			var userModel = models.User{}

			isSuccess, err := userModel.Login(loginData)
			if err != nil || !isSuccess || userModel.Group.Name != enums.Customer {
				return nil, jwt.ErrFailedAuthentication
			}
			return &userModel, nil
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			return true
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		TokenLookup:   "header: Authorization, query: token, cookie: jwt",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	})
	return authMiddleware, err
}

func SellerAuthMiddleware() (*jwt.GinJWTMiddleware, error) {
	secretKey := os.Getenv("JWT_SELLER_SECRET_KEY")

	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "jwt",
		Key:         []byte(secretKey),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		IdentityKey: identityKey,
		PayloadFunc: PayloadFunc,
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginData forms.UserSignIn
			if err := c.ShouldBind(&loginData); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			var userModel = models.User{}

			isSuccess, err := userModel.Login(loginData)
			if err != nil || !isSuccess || userModel.Group.Name != enums.Seller {
				return nil, jwt.ErrFailedAuthentication
			}
			return &userModel, nil
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			return true
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		TokenLookup:   "header: Authorization, query: token, cookie: jwt",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	})
	return authMiddleware, err
}
