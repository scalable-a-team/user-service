package models

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/html"
	"gorm.io/gorm"
	"strings"
	"user-service/db"
	"user-service/enums"
	"user-service/forms"
)

type UserGroup struct {
	ID   int
	Name string
}

func (ug *UserGroup) GetGroup(groupName string) error {
	if err := db.GetDB().Where("name = ?", groupName).FirstOrCreate(ug).Error; err != nil {
		return err
	}
	return nil
}

type UserProfile struct {
	ID        int
	FirstName string
	LastName  string
}

type User struct {
	ID        int
	Username  string
	Password  string
	Profile   UserProfile `gorm:"foreignKey:ProfileID"`
	ProfileID int
	Group     UserGroup `gorm:"foreignKey:GroupID"`
	GroupID   int
}

func (u *User) CreateCustomer(registerForm forms.UserSignUp) (*User, error) {
	var userExists bool
	userResultError := db.GetDB().
		Model(&User{}).
		Select("count(*) > 0").
		Where("username = ?", registerForm.Username).
		Find(&userExists).Error
	if userResultError != nil {
		return &User{}, userResultError
	}
	if userExists {
		return &User{}, errors.New("user already exists")
	}

	var userGroup UserGroup
	if err := userGroup.GetGroup(enums.Customer); err != nil {
		return &User{}, err
	}
	var profile = UserProfile{
		FirstName: registerForm.FirstName,
		LastName:  registerForm.LastName,
	}
	var user = User{
		Username: registerForm.Username,
		Password: registerForm.Password,
		Profile:  profile,
		Group:    userGroup,
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

func (u *User) BeforeCreate(tx *gorm.DB) error {

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
