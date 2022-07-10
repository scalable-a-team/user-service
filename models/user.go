package models

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"html"
	"strings"
	"user-service/db"
	"user-service/forms"
)

type Buyer struct {
	gorm.Model
	Username     string
	Password     string
	BuyerProfile BuyerProfile `gorm:"OnDelete:CASCADE"`
	BuyerWallet  BuyerWallet
}

type BuyerProfile struct {
	gorm.Model
	FirstName string
	LastName  string
	BuyerID   uint
}

type BuyerWallet struct {
	gorm.Model
	Balance uint
	BuyerID uint
}

func (u *Buyer) RetrieveByUsername(username string) error {
	if err := db.GetDB().Where("username = ?", username).First(u).Error; err != nil {
		return err
	}
	return nil
}

func (u *Buyer) RetrieveByUsernameWithProfile(username string) error {
	if err := db.GetDB().
		Where("username = ?", username).
		Preload("BuyerProfile").
		Preload("BuyerWallet").
		First(u).
		Error; err != nil {
		return err
	}
	return nil
}

func (u *Buyer) IsUsernameExist(username string) (bool, error) {
	var userExists bool
	userResultError := db.GetDB().
		Model(&Buyer{}).
		Select("count(*) > 0").
		Where("username = ?", username).
		Find(&userExists).Error
	if userResultError != nil {
		return false, userResultError
	}
	return userExists, nil
}

func (u *Buyer) CreateAccount(registerForm forms.UserSignUp) (*Buyer, error) {
	userExists, userResultError := u.IsUsernameExist(registerForm.Username)
	if userResultError != nil {
		return &Buyer{}, userResultError
	}
	if userExists {
		return &Buyer{}, errors.New("user already exists")
	}

	var profile = BuyerProfile{
		FirstName: registerForm.FirstName,
		LastName:  registerForm.LastName,
	}
	var user = Buyer{
		Username:     registerForm.Username,
		Password:     registerForm.Password,
		BuyerProfile: profile,
		BuyerWallet: BuyerWallet{
			Balance: 0,
		},
	}

	if err := db.GetDB().Create(&user).Error; err != nil {
		return &Buyer{}, err
	}
	return &user, nil
}

func (u *Buyer) Login(form forms.UserSignIn) (bool, error) {
	if err := db.GetDB().
		Where("username = ?", form.Username).
		Preload("BuyerProfile").
		Preload("BuyerWallet").
		First(u).Error; err != nil {
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

func (u *Buyer) BeforeCreate(tx *gorm.DB) error {

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

type Seller struct {
	gorm.Model
	Username      string
	Password      string
	SellerProfile SellerProfile `gorm:"OnDelete:CASCADE"`
	SellerWallet  SellerWallet
}

type SellerWallet struct {
	gorm.Model
	Balance  uint
	SellerID uint
}

type SellerProfile struct {
	gorm.Model
	FirstName string
	LastName  string
	SellerID  uint
}

func (u *Seller) RetrieveByUsername(username string) error {
	if err := db.GetDB().Where("username = ?", username).First(u).Error; err != nil {
		return err
	}
	return nil
}

func (u *Seller) RetrieveByUsernameWithProfile(username string) error {
	if err := db.GetDB().
		Where("username = ?", username).
		Preload("SellerProfile").
		Preload("SellerWallet").
		First(u).
		Error; err != nil {
		return err
	}
	return nil
}

func (u *Seller) IsUsernameExist(username string) (bool, error) {
	var userExists bool
	userResultError := db.GetDB().
		Model(&Seller{}).
		Select("count(*) > 0").
		Where("username = ?", username).
		Find(&userExists).Error
	if userResultError != nil {
		return false, userResultError
	}
	return userExists, nil
}

func (u *Seller) CreateAccount(registerForm forms.UserSignUp) (*Seller, error) {
	userExists, userResultError := u.IsUsernameExist(registerForm.Username)
	if userResultError != nil {
		return &Seller{}, userResultError
	}
	if userExists {
		return &Seller{}, errors.New("user already exists")
	}

	var profile = SellerProfile{
		FirstName: registerForm.FirstName,
		LastName:  registerForm.LastName,
	}
	var user = Seller{
		Username:      registerForm.Username,
		Password:      registerForm.Password,
		SellerProfile: profile,
		SellerWallet: SellerWallet{
			Balance: 0,
		},
	}

	if err := db.GetDB().Create(&user).Error; err != nil {
		return &Seller{}, err
	}
	return &user, nil
}

func (u *Seller) Login(form forms.UserSignIn) (bool, error) {
	if err := db.GetDB().
		Where("username = ?", form.Username).
		Preload("SellerProfile").
		Preload("SellerWallet").
		First(u).Error; err != nil {
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

func (u *Seller) BeforeCreate(tx *gorm.DB) error {

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
