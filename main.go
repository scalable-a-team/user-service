package main

import (
	"context"
	"fmt"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"log"
	"net/http"
	"os"
	"time"
	"user-service/controllers"
	"user-service/db"
	_ "user-service/docs"
	"user-service/middlewares"
	"user-service/models"
	"user-service/otl"
)

// @title           Buyer Service API
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

// @securityDefinitions.apikey ApiKeyAuth  Authorization
//@in header
//@name Authorization
func main() {
	cleanUp := otl.InitTracer()
	defer cleanUp(context.Background())
	port := os.Getenv("PORT")

	if _, isCitusEnabled := os.LookupEnv("CITUS_ENABLED"); isCitusEnabled {
		time.Sleep(30 * time.Second)
	}
	dbInstance := db.Init()
	err := dbInstance.AutoMigrate(
		&models.Seller{},
		&models.SellerProfile{},
		&models.SellerWallet{},
		&models.Buyer{},
		&models.BuyerProfile{},
		&models.BuyerWallet{},
	)
	if err != nil {
		fmt.Println(err)
	}
	if _, isCitusEnabled := os.LookupEnv("CITUS_ENABLED"); isCitusEnabled {
		if err := dbInstance.Exec("SELECT create_distributed_table('sellers', 'id')").Error; err != nil {
			fmt.Println("some issue creating distributed table")
			fmt.Println(err)
			panic("create distributed failed")
		}
	}

	r := gin.Default()
	r.Use(otelgin.Middleware("UserService"))
	// the jwt middleware
	middlewares.InitCustomerJWTMiddleware()
	middlewares.InitSellerJWTMiddleware()
	r.GET("/api/user/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET("api/user/debug", getClaims)

	customerRouter := r.Group("/api/user/customer")
	customerRouter.POST("/login", controllers.BuyerLogin)
	customerRouter.POST("/register", controllers.RegisterCustomer)
	customerRouter.POST("/refresh_token", controllers.BuyerRefreshTokenHandler)
	customerRouter.GET("/profile", controllers.GetBuyerProfileHandler)
	customerRouter.POST("/increase_balance", controllers.AddBuyerWalletBalance)

	sellerRouter := r.Group("/api/user/seller")
	sellerRouter.POST("/login", controllers.SellerLogin)
	sellerRouter.POST("/register", controllers.SellerRegister)
	sellerRouter.GET("/refresh_token", controllers.SellerRefreshToken)
	sellerRouter.GET("/profile", controllers.GetSellerProfile)

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
