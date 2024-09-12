package services

import (
	"errors"
	"fmt"
	"go-template/config"
	"go-template/internal/core/vaults"
	vt "go-template/internal/mocks/vaults_test"

	"go-template/models"
	"go-template/utils/donkeytest"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/stretchr/testify/assert"
)

func TestNewUserService(t *testing.T) {
	type args struct {
		userRepository vaults.User
		apiRepository  vaults.API
	}

	controller := gomock.NewController(t)
	defer controller.Finish()

	mockUserRepository := vt.NewMockUser(controller)
	mockApiRepository := vt.NewMockAPI(controller)

	tests := []struct {
		name string
		args args
		want UserImpl
	}{
		{
			name: "NewUser Test",
			args: args{userRepository: mockUserRepository, apiRepository: mockApiRepository},
			want: UserImpl{user: mockUserRepository, api: mockApiRepository},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewUser(tt.args.userRepository, tt.args.apiRepository), "NewUser(%v)", tt.args.userRepository)
		})
	}
}

func TestUserServiceImpl_GetUser(t *testing.T) {
	type fields struct {
		userRepository *vt.MockUser
		ctx            echo.Context
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockName := "mockName"
	mockEmail := "mockEmail"

	token := jwt.New()
	_ = token.Set("name", mockName)
	_ = token.Set("email", mockEmail)

	c.Set(config.UserKey, token)

	mockUserDB := models.User{
		UserDTO: models.UserDTO{DisplayName: mockName,
			Email: mockEmail},
	}

	tests := []struct {
		name    string
		prepare func(f *fields)
		want    *models.User
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "GetUser case: Found in db",
			prepare: func(f *fields) {
				f.userRepository.EXPECT().Stringify(gomock.Any()).Return(&mockName, &mockEmail, nil)
			},
			want:    &mockUserDB,
			wantErr: assert.NoError,
		},
		{
			name: "GetUser case: Found in db different displayname",
			prepare: func(f *fields) {
				f.userRepository.EXPECT().Stringify(gomock.Any()).Return(&mockName, &mockEmail, nil)
			},
			want:    &mockUserDB,
			wantErr: assert.NoError,
		},
		{
			name: "GetUser case: Found in db different displayname, and cannot save",
			prepare: func(f *fields) {
				f.userRepository.EXPECT().Stringify(gomock.Any()).Return(&mockName, &mockEmail, nil)
			},
			want:    &mockUserDB,
			wantErr: assert.NoError,
		},
		{
			name: "GetUser case: Not found in db",
			prepare: func(f *fields) {
				f.userRepository.EXPECT().Stringify(gomock.Any()).Return(&mockName, &mockEmail, nil)
			},
			want:    &mockUserDB,
			wantErr: assert.NoError,
		},
		{
			name: "GetUser case: cannot get token",
			prepare: func(f *fields) {
				e := echo.New()
				req := httptest.NewRequest(http.MethodPost, "/", nil)
				req.Header.Set("Content-Type", "application/json")
				rec := httptest.NewRecorder()
				context := e.NewContext(req, rec)

				f.ctx = context
			},
			want:    nil,
			wantErr: assert.Error,
		},
		{
			name: "GetUser case: cannot get name",
			prepare: func(f *fields) {
				e := echo.New()
				req := httptest.NewRequest(http.MethodPost, "/", nil)
				req.Header.Set("Content-Type", "application/json")
				rec := httptest.NewRecorder()
				context := e.NewContext(req, rec)

				f.ctx = context

				token := jwt.New()
				_ = token.Set("email", mockEmail)
			},
			want:    nil,
			wantErr: assert.Error,
		},
		{
			name: "GetUser case: cannot get email",
			prepare: func(f *fields) {
				e := echo.New()
				req := httptest.NewRequest(http.MethodPost, "/", nil)
				req.Header.Set("Content-Type", "application/json")
				rec := httptest.NewRecorder()
				context := e.NewContext(req, rec)

				f.ctx = context

				token := jwt.New()
				_ = token.Set("name", mockName)
			},
			want:    nil,
			wantErr: assert.Error,
		},
		{
			name: "GetUser case: cannot stringify token",
			prepare: func(f *fields) {
				f.userRepository.EXPECT().Stringify(gomock.Any()).Return(nil, nil, errors.New(""))
			},
			want:    nil,
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := donkeytest.StartTest(t)

			f := fields{
				userRepository: vt.NewMockUser(ctrl),
				ctx:            c,
			}
			if tt.prepare != nil {
				tt.prepare(&f)
			}

			u := UserImpl{
				user: f.userRepository,
			}

			got1, got2, err := u.ExtractUserInformation(f.ctx)
			if !tt.wantErr(t, err, fmt.Sprintf("GetUser(%v)", f.ctx)) {
				return
			}
			if tt.want == nil {
				return
			}
			assert.Equalf(t, tt.want.DisplayName, *got1, "GetUser(%v)", f.ctx)
			assert.Equalf(t, tt.want.Email, *got2, "GetUser(%v)", f.ctx)
		})
	}
}

func TestUserServiceImpl_GetUser1(t *testing.T) {
	type fields struct {
		userRepository *vt.MockUser
		ctx            echo.Context
	}

	type args struct {
		nameStr  *string
		emailStr *string
	}
	type testSetup struct {
		name    string
		prepare func(f *fields)
		args    args
		want    *models.User
		wantErr assert.ErrorAssertionFunc
	}

	mockName := "mockName"
	mockEmail := "mockEmail"
	user := models.User{
		UserDTO: models.UserDTO{DisplayName: mockName,
			Email: mockEmail},
	}

	tests := []testSetup{
		// TODO: Add test cases.
		{
			name: "Success",
			prepare: func(f *fields) {
				f.userRepository.EXPECT().FetchByEmail(gomock.Any()).Return(nil)
				f.userRepository.EXPECT().Save(gomock.Any()).Return(nil)
			},
			args: args{
				nameStr:  &mockName,
				emailStr: &mockEmail,
			},
			want:    &user,
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := donkeytest.StartTest(t)

			f := fields{
				userRepository: vt.NewMockUser(ctrl),
			}
			if tt.prepare != nil {
				tt.prepare(&f)
			}

			u := UserImpl{
				user: f.userRepository,
			}
			got, err := u.GetUser(tt.args.nameStr, tt.args.emailStr)
			if !tt.wantErr(t, err, fmt.Sprintf("GetUser(%v, %v)", tt.args.nameStr, tt.args.emailStr)) {
				return
			}
			assert.Equalf(t, tt.want, got, "GetUser(%v, %v)", tt.args.nameStr, tt.args.emailStr)
		})
	}
}

func TestUserServiceImpl_AllowAction(t *testing.T) {
	type fields struct {
		userRepository *vt.MockUser
		apiRepository  *vt.MockAPI
	}
	type args struct {
		permissions *models.Permissions
		action      models.TransactionAction
	}

	list := []models.Permission{{
		Action:    models.CreateTx,
		IsGranted: true,
	}}

	permissions := models.Permissions{List: &list}

	tests := []struct {
		name    string
		prepare func(*fields)
		args    args
		want    bool
	}{
		{
			name: "Success",
			prepare: func(fields *fields) {

			},
			args: args{
				permissions: &permissions,
				action:      models.CreateTx,
			},
			want: true,
		},
		{
			name: "No permission",
			prepare: func(fields *fields) {

			},
			args: args{
				permissions: &permissions,
				action:      models.UpdateTx,
			},
			want: false,
		},
		{
			name: "No permission given",
			prepare: func(fields *fields) {

			},
			args: args{
				permissions: nil,
				action:      models.UpdateTx,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := donkeytest.StartTest(t)

			f := fields{
				apiRepository:  vt.NewMockAPI(ctrl),
				userRepository: vt.NewMockUser(ctrl),
			}

			if tt.prepare != nil {
				tt.prepare(&f)
			}

			u := UserImpl{
				user: f.userRepository,
				api:  f.apiRepository,
			}

			assert.Equalf(t, tt.want, u.AllowAction(tt.args.permissions, tt.args.action), "AllowAction(%v, %v)", tt.args.permissions, tt.args.action)
		})
	}
}

func TestUserServiceImpl_CheckPermissions(t *testing.T) {
	type fields struct {
		userRepository vaults.User
		apiRepository  vaults.API
	}
	type args struct {
		roles        *models.Roles
		userInSystem *models.User
		user         *models.User
	}

	rolesA := []models.Role{models.AdminRole}
	rolesU := []models.Role{models.UserRole}

	deRolesA := models.Roles{Roles: &rolesA}
	deRolesU := models.Roles{Roles: &rolesU}

	user1 := models.User{
		BaseModel: models.BaseModel{
			ID: 2,
		},
	}
	user2 := models.User{
		BaseModel: models.BaseModel{
			ID: 2,
		},
	}

	userActions := []models.Permission{
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
			IsGranted: false,
		},
		{
			Action:    models.ReadTx,
			IsGranted: true,
		}}

	adminActions := []models.Permission{
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

	adminPermissions := models.Permissions{List: &adminActions}

	userPermissions := models.Permissions{List: &userActions}

	tests := []struct {
		name    string
		prepare func(*fields)
		args    args
		want    models.Permissions
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Success - admin",
			prepare: func(fields *fields) {

			},
			args: args{
				roles:        &deRolesA,
				userInSystem: &user1,
				user:         &user1,
			},
			want:    adminPermissions,
			wantErr: assert.NoError,
		},
		{
			name: "Success - user",
			prepare: func(fields *fields) {

			},
			args: args{
				roles:        &deRolesU,
				userInSystem: &user1,
				user:         &user2,
			},
			want:    userPermissions,
			wantErr: assert.NoError,
		},
		{
			name: "Error - no fields given",
			prepare: func(fields *fields) {

			},
			args: args{
				roles:        nil,
				userInSystem: nil,
				user:         nil,
			},
			want:    models.Permissions{},
			wantErr: assert.Error,
		},
		{
			name: "Error - no fields given",
			prepare: func(fields *fields) {

			},
			args: args{
				roles:        &deRolesA,
				userInSystem: nil,
				user:         nil,
			},
			want:    models.Permissions{},
			wantErr: assert.Error,
		},
		{
			name: "Error - no fields given",
			prepare: func(fields *fields) {
				var s []models.Role
				deRolesA.Roles = &s
			},
			args: args{
				roles:        &deRolesA,
				userInSystem: nil,
				user:         nil,
			},
			want:    models.Permissions{},
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := donkeytest.StartTest(t)

			f := fields{
				apiRepository:  vt.NewMockAPI(ctrl),
				userRepository: vt.NewMockUser(ctrl),
			}

			if tt.prepare != nil {
				tt.prepare(&f)
			}

			u := UserImpl{
				user: f.userRepository,
				api:  f.apiRepository,
			}
			got, _, err := u.CheckPermissions(tt.args.roles, tt.args.userInSystem, tt.args.user)
			if !tt.wantErr(t, err, fmt.Sprintf("CheckPermissions(%v, %v, %v)", tt.args.roles, tt.args.userInSystem, tt.args.user)) {
				return
			}
			assert.Equalf(t, tt.want, got, "CheckPermissions(%v, %v, %v)", tt.args.roles, tt.args.userInSystem, tt.args.user)
		})
	}
}
