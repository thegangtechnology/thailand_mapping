package services

import (
	"go-template/config"
	"go-template/internal/core/vaults"
	"go-template/models"
	"go-template/utils/slices"

	"gorm.io/gorm"

	"github.com/labstack/echo/v4"
	"github.com/lestrrat-go/jwx/jwt"
	log "github.com/sirupsen/logrus"
)

//go:generate mockgen -source=$GOFILE -destination=../../mocks/services_test/$GOFILE -package=services_test

type UserImpl struct {
	user vaults.User
	api  vaults.API
}

type User interface {
	CheckPermissions(roles *models.Roles, userInSystem, user *models.User) (models.Permissions, models.Role, error)
	AllowAction(permissions *models.Permissions, update models.TransactionAction) bool

	ExtractUserInformation(ctx echo.Context) (*string, *string, error)
	GetUser(nameStr, emailStr *string) (*models.User, error)
}

func NewUser(userRepository vaults.User, apiRepository vaults.API) User {
	return UserImpl{
		user: userRepository,
		api:  apiRepository,
	}
}

func (u UserImpl) WithTrx(trxHandle *gorm.DB) User {
	u.user = u.user.WithTrx(trxHandle)

	return u
}

func (u UserImpl) AllowAction(permissions *models.Permissions, action models.TransactionAction) bool {
	if permissions == nil || permissions.List == nil {
		return false
	}

	for _, permission := range *permissions.List {
		if permission.Action == action {
			return permission.IsGranted
		}
	}

	return false
}

func (u UserImpl) CheckPermissions(roles *models.Roles, userInSystem,
	user *models.User) (models.Permissions, models.Role, error) {
	var (
		permissions models.Permissions
		isRole      bool
	)

	if roles == nil || roles.Roles == nil {
		return permissions, "", config.ErrNoRoles
	}

	if user == nil {
		return permissions, "", config.ErrNoUserGiven
	}

	if permissions, isRole = u.isAdmin(roles); isRole {
		return permissions, models.AdminRole, nil
	}

	allowDelete := userInSystem == nil

	if permissions, isRole = u.isUser(roles, allowDelete); isRole {
		return permissions, models.UserRole, nil
	}

	return permissions, "", config.ErrCannotDetectRoles
}

func (u UserImpl) GetUser(nameStr, emailStr *string) (*models.User, error) {
	var user models.User

	if nameStr == nil || emailStr == nil {
		return nil, config.ErrNoUserGiven
	}

	user.Email = *emailStr
	if err := u.user.FetchByEmail(&user); err != nil {
		log.WithField("userEmail", emailStr).Debugln("user not found, lazy creating...")

		user = models.User{
			UserDTO: models.UserDTO{
				DisplayName: *nameStr,
				Email:       *emailStr,
			},
		}

		if err := u.user.Create(&user); err != nil {
			return nil, config.ErrUserCreate
		}
	}

	if user.DisplayName != *nameStr {
		user.DisplayName = *nameStr
		if err := u.user.Save(&user); err != nil {
			return &user, err
		}
	}

	return &user, nil
}

func (u UserImpl) ExtractUserInformation(ctx echo.Context) (nameStr, emailStr *string, err error) {
	token, ok := ctx.Get(config.UserKey).(jwt.Token)
	if !ok {
		log.WithError(err).Error("token not found in context")

		return nil, nil, config.ErrTokenNotFound
	}

	nameStr, emailStr, err = u.user.Stringify(token)
	if err != nil {
		return nil, nil, err
	}

	return
}

func (u UserImpl) isAdmin(roles *models.Roles) (permissions models.Permissions, is bool) {
	if !slices.Contains(*roles.Roles, models.AdminRole) {
		return permissions, false
	}

	actions := []models.Permission{
		{
			Action:    models.CreateTx,
			IsGranted: true,
		},
		{
			Action:    models.UpdateTx,
			IsGranted: true,
		},
		{
			Action:    models.DeleteTx,
			IsGranted: true,
		},
		{
			Action:    models.ReadTx,
			IsGranted: true,
		},
	}
	permissions.List = &actions

	return permissions, true
}

func (u UserImpl) isUser(roles *models.Roles, allowDelete bool) (permissions models.Permissions, isUser bool) {
	if !slices.Contains(*roles.Roles, models.UserRole) {
		return permissions, false
	}

	actions := []models.Permission{
		{
			Action:    models.CreateTx,
			IsGranted: true,
		},
		{
			Action:    models.UpdateTx,
			IsGranted: true,
		},
		{
			Action:    models.DeleteTx,
			IsGranted: allowDelete,
		},
		{
			Action:    models.ReadTx,
			IsGranted: true,
		},
	}

	permissions.List = &actions

	return permissions, true
}
