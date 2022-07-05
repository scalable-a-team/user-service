package models

import (
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/html"
	"strings"
	"user-service/db"
	"user-service/forms"
)

type UserProfile struct {
	FirstName string
	LastName  string
}

type User struct {
	Username string
	Password string
	Profile  UserProfile
}

func (u *User) CreateUser(registerForm forms.UserSignUp) (*User, error) {

	var profile = UserProfile{
		FirstName: registerForm.FirstName,
		LastName:  registerForm.LastName,
	}
	var user = User{
		Username: registerForm.Username,
		Password: registerForm.Password,
		Profile:  profile,
	}
	if err := db.GetDB().Create(&user).Error; err != nil {
		return &User{}, err
	}
	return &user, nil
}

func (u *User) Login(form forms.UserSignIn) (bool, error) {
	if err := db.GetDB().Where("username = ?", form.Username).First(u).Error; err != nil {
		return false, err
	}

	//Compare the password form and database if match
	bytePassword := []byte(form.Password)
	byteHashedPassword := []byte(u.Password)

	if err := bcrypt.CompareHashAndPassword(byteHashedPassword, bytePassword); err != nil {
		return false, err
	}

	return true, nil
}

func (u *User) BeforeCreate() error {

	//turn password into hash
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)

	//remove spaces in username
	u.Username = html.EscapeString(strings.TrimSpace(u.Username))

	return nil

}
