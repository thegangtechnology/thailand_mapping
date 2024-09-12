package vaults

import (
	"go-template/config"
	"go-template/models"

	log "github.com/sirupsen/logrus"

	"github.com/lestrrat-go/jwx/jwt"
	"gorm.io/gorm"
)

//go:generate mockgen -source=$GOFILE -destination=../../mocks/vaults_test/$GOFILE -package=vaults_test

type UserImpl struct {
	db *gorm.DB
}

type User interface {
	Stringify(token jwt.Token) (name, email *string, err error)
	FetchByEmail(user *models.User) error
	Create(user *models.User) error
	Save(user *models.User) error

	WithTrx(trxHandle *gorm.DB) User
}

func NewUser(db *gorm.DB) User {
	return UserImpl{
		db: db,
	}
}

func (u UserImpl) WithTrx(trxHandle *gorm.DB) User {
	if trxHandle == nil {
		log.Error("Transaction Database not found")

		return u
	}

	u.db = trxHandle

	return u
}

func (u UserImpl) Stringify(token jwt.Token) (name, email *string, err error) {
	nameToken, ok := token.Get("name")
	if !ok {
		return nil, nil, config.ErrNameNotFoundInClaims
	}

	emailToken, ok := token.Get("email")
	if !ok {
		return nil, nil, config.ErrEmailNotFoundInClaims
	}

	nameStr, ok := nameToken.(string)
	if !ok {
		return nil, nil, config.ErrNameCast
	}

	emailStr, ok := emailToken.(string)
	if !ok {
		return nil, nil, config.ErrEmailCast
	}

	return &nameStr, &emailStr, nil
}

func (u UserImpl) FetchByEmail(user *models.User) error {
	err := u.db.Where("email = ?", user.Email).First(&user).Error

	return err
}

func (u UserImpl) Create(user *models.User) error {
	err := u.db.Create(&user).Error

	return err
}

func (u UserImpl) Save(user *models.User) error {
	err := u.db.Save(&user).Error

	return err
}
