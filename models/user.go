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

type RoleGroup struct {
	gorm.Model
	Name string
	User []User
}

func (ug *RoleGroup) GetGroup(groupName string) error {
	if err := db.GetDB().Where(RoleGroup{Name: groupName}).FirstOrCreate(ug).Error; err != nil {
		return err
	}
	return nil
}

type User struct {
	gorm.Model
	Username      string
	Password      string
	Profile       Profile `gorm:"OnDelete:CASCADE"`
	RoleGroupID   uint
	RoleGroupName string `gorm:"-"`
}

type Profile struct {
	gorm.Model
	FirstName string
	LastName  string
	UserID    uint
}

func (u *User) RetrieveByUsername(username string) error {
	if err := db.GetDB().Where("username = ?", username).First(u).Error; err != nil {
		return err
	}
	return nil
}

func (u *User) IsUsernameExist(username string) (bool, error) {
	var userExists bool
	userResultError := db.GetDB().
		Model(&User{}).
		Select("count(*) > 0").
		Where("username = ?", username).
		Find(&userExists).Error
	if userResultError != nil {
		return false, userResultError
	}
	return userExists, nil
}

func (u *User) CreateCustomer(registerForm forms.UserSignUp) (*User, error) {
	userExists, userResultError := u.IsUsernameExist(registerForm.Username)
	if userResultError != nil {
		return &User{}, userResultError
	}
	if userExists {
		return &User{}, errors.New("user already exists")
	}

	var userGroup RoleGroup
	if err := userGroup.GetGroup(enums.Customer); err != nil {
		return &User{}, err
	}
	var profile = Profile{
		FirstName: registerForm.FirstName,
		LastName:  registerForm.LastName,
	}
	var user = User{
		Username:      registerForm.Username,
		Password:      registerForm.Password,
		Profile:       profile,
		RoleGroupID:   userGroup.ID,
		RoleGroupName: enums.Customer,
	}

	if err := db.GetDB().Create(&user).Error; err != nil {
		return &User{}, err
	}
	return &user, nil
}

func (u *User) Login(form forms.UserSignIn) (bool, error) {
	if err := db.GetDB().
		Where("username = ?", form.Username).
		Preload("Profile").
		First(u).Error; err != nil {
		return false, err
	}

	//Compare the password form and database if match
	bytePassword := []byte(form.Password)
	byteHashedPassword := []byte(u.Password)

	if err := bcrypt.CompareHashAndPassword(byteHashedPassword, bytePassword); err != nil {
		return false, err
	}
	var roleGroup RoleGroup
	if err := db.GetDB().Where("id = ?", u.RoleGroupID).FirstOrCreate(&roleGroup).Error; err != nil {
		return false, err
	}
	u.RoleGroupName = roleGroup.Name
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
