package main

import (
	"fmt"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"user-service/controllers"
	"user-service/db"
	"user-service/middlewares"
	"user-service/models"
)

func main() {
	port := os.Getenv("PORT")

	dbInstance := db.Init()
	dbInstance.AutoMigrate(&models.User{}, &models.UserProfile{})
	r := gin.Default()

	// the jwt middleware
	customerAuthMiddleware, err := middlewares.CustomerAuthMiddleware()
	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}
	// When you use jwt.New(), the function is already automatically called for checking,
	// which means you don't need to call it again.
	errInit := customerAuthMiddleware.MiddlewareInit()
	if errInit != nil {
		log.Fatal("customerAuthMiddleware.MiddlewareInit() Error:" + errInit.Error())
	}

	sellerAuthMiddleware, err := middlewares.SellerAuthMiddleware()
	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}
	errInit = sellerAuthMiddleware.MiddlewareInit()
	if errInit != nil {
		log.Fatal("sellerAuthMiddleware.MiddlewareInit() Error:" + errInit.Error())
	}

	r.GET("api/auth/debug", getClaims)

	customerRouter := r.Group("/api/auth/customer")
	customerRouter.POST("/login", customerAuthMiddleware.LoginHandler)
	customerRouter.POST("/register", controllers.RegisterCustomer)
	customerRouter.GET("/refresh_token", customerAuthMiddleware.RefreshHandler)

	sellerRouter := r.Group("/api/auth/seller")
	sellerRouter.POST("/login", sellerAuthMiddleware.LoginHandler)
	//sellerRouter.POST("/register", controllers.RegisterCustomer)
	sellerRouter.GET("/refresh_token", sellerAuthMiddleware.RefreshHandler)

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
