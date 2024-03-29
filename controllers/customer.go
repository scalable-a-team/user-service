package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"user-service/enums"
	"user-service/forms"
	"user-service/middlewares"
	"user-service/models"
	"user-service/service"
)

// PingExample godoc
// @Summary BuyerLogin user
// @Schemes
// @Description Return JWT access and refresh pair, alongside user profile
// @Tags example
// @Accept json
// @Produce json
// @Param data body forms.UserSignIn true "BuyerLogin input"
// @Success 200 {object} forms.LoginResponse
// @Router /customer/login [post]
func BuyerLogin(c *gin.Context) {
	var loginData forms.UserSignIn
	if err := c.ShouldBind(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var userModel = models.Buyer{}
	isSuccess, err := userModel.Login(c.Request.Context(), loginData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if !isSuccess {
		c.JSON(http.StatusBadRequest, gin.H{"error": "authentication failed"})
		return
	}

	tokenUserInput := service.TokenUserInput{
		Username:      userModel.Username,
		UserID:        userModel.ID,
		RoleGroupName: enums.Buyer,
		Firstname:     userModel.BuyerProfile.FirstName,
		Lastname:      userModel.BuyerProfile.LastName,
	}
	tokenString, err := middlewares.GetCustomerJwtMiddleware().GenerateAccessToken(&tokenUserInput)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	refreshTokenString, err := middlewares.GetCustomerJwtMiddleware().GenerateRefreshToken(&tokenUserInput)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var loginResponse = forms.LoginResponse{
		Token:   tokenString,
		Refresh: refreshTokenString,
		User:    generateBuyerData(userModel),
	}
	c.JSON(http.StatusOK, loginResponse)
}

// PingExample godoc
// @Summary Register customer
// @Schemes
// @Description Register buyer account
// @Tags example
// @Accept json
// @Produce json
// @Param data body forms.UserSignUp true "Signup input"
// @Success 200 {object} forms.LoginResponse
// @Router /customer/register [post]
func RegisterCustomer(c *gin.Context) {
	var input forms.UserSignUp

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var userModel = new(models.Buyer)
	newUser, err := userModel.CreateAccount(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tokenUserInput := service.TokenUserInput{
		Username:      newUser.Username,
		UserID:        newUser.ID,
		RoleGroupName: enums.Buyer,
		Firstname:     newUser.BuyerProfile.FirstName,
		Lastname:      newUser.BuyerProfile.LastName,
	}
	tokenString, err := middlewares.GetCustomerJwtMiddleware().GenerateAccessToken(&tokenUserInput)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	refreshTokenString, err := middlewares.GetCustomerJwtMiddleware().GenerateRefreshToken(&tokenUserInput)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var loginResponse = forms.LoginResponse{
		Token:   tokenString,
		Refresh: refreshTokenString,
		User:    generateBuyerData(*newUser),
	}
	c.JSON(http.StatusOK, loginResponse)
}

// PingExample godoc
// @Summary Refresh token handler
// @Schemes
// @Description Return JWT access token given refresh token
// @Tags example
// @Accept json
// @Produce json
// @Param data body forms.RefreshTokenRequest true "Receive refresh token"
// @Success 200 {string} refresh_token
// @Router /customer/refresh_token [post]
func BuyerRefreshTokenHandler(c *gin.Context) {
	var input forms.RefreshTokenRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	claims, err := middlewares.GetCustomerJwtMiddleware().ValidateRefreshAccessToken(input.RefreshToken)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID, err := uuid.Parse(claims["userid"].(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var user models.Buyer
	if err := user.RetrieveByUserIDWithProfile(c.Request.Context(), userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tokenUserInput := service.TokenUserInput{
		Username:      user.Username,
		UserID:        user.ID,
		RoleGroupName: enums.Buyer,
		Firstname:     user.BuyerProfile.FirstName,
		Lastname:      user.BuyerProfile.LastName,
	}
	accessToken, err := middlewares.GetCustomerJwtMiddleware().GenerateAccessToken(&tokenUserInput)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"access_token": accessToken})
}

// PingExample godoc
// @Summary Get Buyer BuyerProfile
// @Schemes
// @Description Get customer profile from Authorization JWT header
// @Tags example
// @Accept json
// @Produce json
// @Security JWT Key
// @param Authorization header string true "Bearer YourJWTToken"
// @Success 200 {object} forms.UserResponse
// @Router /customer/profile [get]
func GetBuyerProfileHandler(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	tokenString := authHeader[len("Bearer "):]
	userID, err := middlewares.GetCustomerJwtMiddleware().GetUserIDFromToken(tokenString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user := models.Buyer{}

	err = user.RetrieveByUserIDWithProfile(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	loginResponse := generateBuyerData(user)
	c.JSON(http.StatusOK, loginResponse)
}

func generateBuyerData(userModel models.Buyer) forms.UserResponse {
	return forms.UserResponse{
		ID:       userModel.ID,
		Username: userModel.Username,
		Profile: forms.UserProfileResponse{
			FirstName: userModel.BuyerProfile.FirstName,
			LastName:  userModel.BuyerProfile.LastName,
		},
		Group: forms.UserGroupResponse{
			Name: enums.Buyer,
		},
		WalletBalance: userModel.BuyerWallet.Balance,
	}
}

// PingExample godoc
// @Summary Topup customer wallet balance to purchase stuff
// @Schemes
// @Description Increase custom wallet balance
// @Tags example
// @Accept json
// @Produce json
// @Security JWT Key
// @param Authorization header string true "Bearer YourJWTToken"
// @Param data body forms.AddWalletBalanceInput true "Increment balance by certain amount"
// @Success 200 {object} forms.AddWalletBalanceResponse
// @Router /customer/increase_balance [post]
func AddBuyerWalletBalance(c *gin.Context) {
	var input forms.AddWalletBalanceInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	authHeader := c.GetHeader("Authorization")
	tokenString := authHeader[len("Bearer "):]

	userID, err := middlewares.GetCustomerJwtMiddleware().GetUserIDFromToken(tokenString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user := models.Buyer{}

	err = user.RetrieveByUserIDWithProfile(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedBalance, err := user.AddBalance(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, forms.AddWalletBalanceResponse{NewBalance: updatedBalance})
}
