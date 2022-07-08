package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"user-service/enums"
	"user-service/forms"
	"user-service/middlewares"
	"user-service/models"
)

// PingExample godoc
// @Summary Login user
// @Schemes
// @Description Return JWT access and refresh pair, alongside user profile
// @Tags example
// @Accept json
// @Produce json
// @Param data body forms.UserSignIn true "Login input"
// @Success 200 {object} forms.LoginResponse
// @Router /user/login [post]
func Login(c *gin.Context) {
	var loginData forms.UserSignIn
	if err := c.ShouldBind(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var userModel = models.User{}
	isSuccess, err := userModel.Login(loginData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if !isSuccess {
		c.JSON(http.StatusBadRequest, gin.H{"error": "authentication failed"})
		return
	}

	if userModel.RoleGroupName != enums.Customer {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user group for this endpoint"})
		return
	}

	tokenString, err := middlewares.GetCustomerJwtMiddleware().GenerateAccessToken(&userModel)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	refreshTokenString, err := middlewares.GetCustomerJwtMiddleware().GenerateRefreshToken(&userModel)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var loginResponse = forms.LoginResponse{
		Token:   tokenString,
		Refresh: refreshTokenString,
		User:    generateUserData(userModel),
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
// @Router /user/register [post]
func RegisterCustomer(c *gin.Context) {
	var input forms.UserSignUp

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var userModel = new(models.User)
	newUser, err := userModel.CreateCustomer(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tokenString, err := middlewares.GetCustomerJwtMiddleware().GenerateAccessToken(newUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	refreshTokenString, err := middlewares.GetCustomerJwtMiddleware().GenerateRefreshToken(newUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var loginResponse = forms.LoginResponse{
		Token:   tokenString,
		Refresh: refreshTokenString,
		User:    generateUserData(*newUser),
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
// @Router /user/refresh_token [post]
func RefreshTokenHandler(c *gin.Context) {
	var input forms.RefreshTokenRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := middlewares.GetCustomerJwtMiddleware().RefreshAccessToken(input.RefreshToken)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"access_token": token})
}

func generateUserData(userModel models.User) forms.UserResponse {
	return forms.UserResponse{
		ID:       userModel.ID,
		Username: userModel.Username,
		Profile: forms.UserProfileResponse{
			FirstName: userModel.Profile.FirstName,
			LastName:  userModel.Profile.LastName,
		},
		Group: forms.UserGroupResponse{
			Name: userModel.RoleGroupName,
		},
	}
}
