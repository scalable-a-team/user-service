package main

import (
	"fmt"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"log"
	"net/http"
	"os"
	"user-service/controllers"
	"user-service/db"
	_ "user-service/docs"
	"user-service/middlewares"
	"user-service/models"
)

// @title           User Service API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8000
// @BasePath  /api/user

// @securityDefinitions.basic  BasicAuth
func main() {
	port := os.Getenv("PORT")

	dbInstance := db.Init()
	err := dbInstance.AutoMigrate(&models.RoleGroup{}, &models.User{}, &models.Profile{})
	if err != nil {
		fmt.Println(err)
	}
	r := gin.Default()

	// the jwt middleware
	middlewares.InitCustomerJWTMiddleware()
	middlewares.InitSellerJWTMiddleware()
	r.GET("/api/user/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET("api/user/debug", getClaims)

	customerRouter := r.Group("/api/user/customer")
	customerRouter.POST("/login", controllers.Login)
	customerRouter.POST("/register", controllers.RegisterCustomer)
	customerRouter.POST("/refresh_token", controllers.RefreshTokenHandler)

	//sellerRouter := r.Group("/api/auth/seller")
	//sellerRouter.POST("/login", sellerAuthMiddleware.LoginHandler)
	//sellerRouter.POST("/register", controllers.RegisterCustomer)
	//sellerRouter.GET("/refresh_token", sellerAuthMiddleware.RefreshHandler)

	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal(err)
	}
}

func getClaims(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	for k, vals := range c.Request.Header {
		fmt.Println("%s", k)
		for _, v := range vals {
			fmt.Println("\t%s", v)
		}
	}
	fmt.Println(claims)
	c.Status(http.StatusNoContent)
	return
}
