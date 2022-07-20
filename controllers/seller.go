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
// @Summary SellerLogin user
// @Schemes
// @Description Return JWT access and refresh pair, alongside user profile
// @Tags example
// @Accept json
// @Produce json
// @Param data body forms.UserSignIn true "SellerLogin input"
// @Success 200 {object} forms.LoginResponse
// @Router /seller/login [post]
func SellerLogin(c *gin.Context) {
	var loginData forms.UserSignIn
	if err := c.ShouldBind(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var userModel = models.Seller{}
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
		RoleGroupName: enums.Seller,
	}
	tokenString, err := middlewares.GetSellerJwtMiddleware().GenerateAccessToken(&tokenUserInput)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	refreshTokenString, err := middlewares.GetSellerJwtMiddleware().GenerateRefreshToken(&tokenUserInput)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var loginResponse = forms.LoginResponse{
		Token:   tokenString,
		Refresh: refreshTokenString,
		User:    generateSellerData(userModel),
	}
	c.JSON(http.StatusOK, loginResponse)
}

// PingExample godoc
// @Summary Register customer
// @Schemes
// @Description Register seller account
// @Tags example
// @Accept json
// @Produce json
// @Param data body forms.UserSignUp true "Signup input"
// @Success 200 {object} forms.LoginResponse
// @Router /seller/register [post]
func SellerRegister(c *gin.Context) {
	var input forms.UserSignUp

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var userModel = new(models.Seller)
	newUser, err := userModel.CreateAccount(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tokenUserInput := service.TokenUserInput{
		Username:      newUser.Username,
		UserID:        newUser.ID,
		RoleGroupName: enums.Seller,
	}

	tokenString, err := middlewares.GetSellerJwtMiddleware().GenerateAccessToken(&tokenUserInput)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	refreshTokenString, err := middlewares.GetSellerJwtMiddleware().GenerateRefreshToken(&tokenUserInput)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var loginResponse = forms.LoginResponse{
		Token:   tokenString,
		Refresh: refreshTokenString,
		User:    generateSellerData(*newUser),
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
// @Router /seller/refresh_token [post]
func SellerRefreshToken(c *gin.Context) {
	var input forms.RefreshTokenRequest
	claims, err := middlewares.GetSellerJwtMiddleware().ValidateRefreshAccessToken(input.RefreshToken)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	username := claims["username"].(string)
	userId := claims["userid"].(uuid.UUID)
	var user models.Seller
	if err := user.RetrieveByUserID(c.Request.Context(), userId); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	accessToken, err := middlewares.GetSellerJwtMiddleware().GenerateAccessToken(
		&service.TokenUserInput{
			Username:      username,
			UserID:        userId,
			RoleGroupName: enums.Seller,
		},
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"access_token": accessToken})
}

// PingExample godoc
// @Summary Get Seller SellerProfile
// @Schemes
// @Description Get seller profile from Authorization JWT header
// @Tags example
// @Accept json
// @Produce json
// @Security JWT Key
// @param Authorization header string true "Bearer YourJWTToken"
// @Success 200 {object} forms.UserResponse
// @Router /seller/profile [get]
func GetSellerProfile(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	tokenString := authHeader[len("Bearer "):]
	userID, err := middlewares.GetSellerJwtMiddleware().GetUserIDFromToken(tokenString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user := models.Seller{}

	err = user.RetrieveByUserIDWithProfile(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	loginResponse := generateSellerData(user)
	c.JSON(http.StatusOK, loginResponse)
}

func generateSellerData(userModel models.Seller) forms.UserResponse {
	return forms.UserResponse{
		ID:       userModel.ID,
		Username: userModel.Username,
		Profile: forms.UserProfileResponse{
			FirstName: userModel.SellerProfile.FirstName,
			LastName:  userModel.SellerProfile.LastName,
		},
		Group: forms.UserGroupResponse{
			Name: enums.Seller,
		},
		WalletBalance: userModel.SellerWallet.Balance,
	}
}
