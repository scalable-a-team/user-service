package models

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"html"
	"strings"
	"time"
	"user-service/db"
	"user-service/forms"
)

type Buyer struct {
	ID           uuid.UUID `gorm:"primarykey;type:uuid;uniqueIndex:username_unique"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Username     string `gorm:"uniqueIndex:username_unique"`
	Password     string
	BuyerProfile BuyerProfile `gorm:"OnDelete:CASCADE"`
	BuyerWallet  BuyerWallet
}

type BuyerProfile struct {
	gorm.Model
	FirstName string
	LastName  string
	BuyerID   uuid.UUID `gorm:"type:uuid;primaryKey"`
}

type BuyerWallet struct {
	gorm.Model
	Balance decimal.Decimal `gorm:"type:decimal(12,2);"`
	BuyerID uuid.UUID       `gorm:"type:uuid;primaryKey"`
}

func (u *Buyer) RetrieveByUsername(c context.Context, username string) error {
	if err := db.GetDB(c).Where("username = ?", username).First(u).Error; err != nil {
		return err
	}
	return nil
}

func (u *Buyer) RetrieveByUsernameWithProfile(c context.Context, username string) error {
	if err := db.GetDB(c).
		Where("username = ?", username).
		Preload("BuyerProfile").
		Preload("BuyerWallet").
		First(u).
		Error; err != nil {
		return err
	}
	return nil
}

func (u *Buyer) IsUsernameExist(c context.Context, username string) (bool, error) {
	var userExists bool
	userResultError := db.GetDB(c).
		Model(&Buyer{}).
		Select("count(*) > 0").
		Where("username = ?", username).
		Find(&userExists).Error
	if userResultError != nil {
		return false, userResultError
	}
	return userExists, nil
}

func (u *Buyer) CreateAccount(c context.Context, registerForm forms.UserSignUp) (*Buyer, error) {
	userExists, userResultError := u.IsUsernameExist(c, registerForm.Username)
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
			Balance: decimal.NewFromInt(0),
		},
	}

	if err := db.GetDB(c).Create(&user).Error; err != nil {
		return &Buyer{}, err
	}
	return &user, nil
}

func (u *Buyer) Login(c context.Context, form forms.UserSignIn) (bool, error) {
	if err := db.GetDB(c).
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
	u.ID = uuid.New()
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

func (u *Buyer) AddBalance(c context.Context, input forms.AddWalletBalanceInput) (decimal.Decimal, error) {
	updatedBalance := u.BuyerWallet.Balance
	err := db.GetDB(c).Transaction(func(tx *gorm.DB) error {
		tmpWallet := BuyerWallet{}
		// do some database operations in the transaction (use 'tx' from this point, not 'db')
		if err := tx.
			Model(&tmpWallet).
			Where("ID = ?", u.BuyerWallet.ID).
			Clauses(clause.Locking{Strength: "UPDATE"}).Find(&tmpWallet).Error; err != nil {
			return err
		}
		updatedBalance = tmpWallet.Balance.Add(input.AddBalance)
		if err := tx.
			Model(&tmpWallet).
			Where("ID = ?", u.BuyerWallet.ID).
			Update("Balance", updatedBalance).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return decimal.NewFromInt(0), err
	}
	return updatedBalance, nil
}

type Seller struct {
	ID            uuid.UUID `gorm:"primarykey;type:uuid;uniqueIndex:username_unique"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Username      string `gorm:"uniqueIndex:username_unique"`
	Password      string
	SellerProfile SellerProfile `gorm:"OnDelete:CASCADE"`
	SellerWallet  SellerWallet
}

type SellerWallet struct {
	gorm.Model
	Balance  decimal.Decimal `gorm:"type:decimal(12,2);"`
	SellerID uuid.UUID       `gorm:"type:uuid;primaryKey"`
}

type SellerProfile struct {
	gorm.Model
	FirstName string
	LastName  string
	SellerID  uuid.UUID `gorm:"type:uuid;primaryKey"`
}

func (u *Seller) RetrieveByUsername(c context.Context, username string) error {
	if err := db.GetDB(c).Where("username = ?", username).First(u).Error; err != nil {
		return err
	}
	return nil
}

func (u *Seller) RetrieveByUsernameWithProfile(c context.Context, username string) error {
	if err := db.GetDB(c).
		Where("username = ?", username).
		Preload("SellerProfile").
		Preload("SellerWallet").
		First(u).
		Error; err != nil {
		return err
	}
	return nil
}

func (u *Seller) IsUsernameExist(c context.Context, username string) (bool, error) {
	var userExists bool
	userResultError := db.GetDB(c).
		Model(&Seller{}).
		Select("count(*) > 0").
		Where("username = ?", username).
		Find(&userExists).Error
	if userResultError != nil {
		return false, userResultError
	}
	return userExists, nil
}

func (u *Seller) CreateAccount(c context.Context, registerForm forms.UserSignUp) (*Seller, error) {
	userExists, userResultError := u.IsUsernameExist(c, registerForm.Username)
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
			Balance: decimal.NewFromInt(0),
		},
	}

	if err := db.GetDB(c).Create(&user).Error; err != nil {
		return &Seller{}, err
	}
	return &user, nil
}

func (u *Seller) Login(c context.Context, form forms.UserSignIn) (bool, error) {
	if err := db.GetDB(c).
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
	u.ID = uuid.New()
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
