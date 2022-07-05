package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"user-service/forms"
	"user-service/models"
)

var userModel = new(models.User)

func Register(c *gin.Context) {

	var input forms.UserSignUp

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if _, err := userModel.CreateUser(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}
