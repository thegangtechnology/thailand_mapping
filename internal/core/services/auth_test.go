package services

import (
	"errors"
	"fmt"
	vt "go-template/internal/mocks/vaults_test"
	"go-template/utils/donkeytest"
	"net/http"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/stretchr/testify/assert"
)

func TestNewJwtAuthenticatorService(t *testing.T) {
	type args struct {
		auth *vt.MockAuthenticator
	}
	type testSetup struct {
		name string
		args args
		want Authenticator
	}

	controller := donkeytest.StartTest(t)
	defer controller.Finish()

	authRepo := vt.NewMockAuthenticator(controller)

	tests := []testSetup{
		{
			name: "Case #1 service created",
			args: args{authRepo},
			want: AuthenticatorServiceImpl{authRepo},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewAuthenticatorService(tt.args.auth), "NewAuthenticatorService(%v)", tt.args.auth)
		})
	}
}

func TestAuthenticatorServiceImpl_ParseJWT(t *testing.T) {
	type fields struct {
		AuthenticatorRepository *vt.MockAuthenticator
	}
	type args struct {
		req *http.Request
	}
	type testSetup struct {
		name    string
		prepare func(f *fields)
		args    args
		want    jwt.Token
		wantErr assert.ErrorAssertionFunc
	}
	tests := []testSetup{
		{
			name: "Case #1 no error",
			prepare: func(f *fields) {
				f.AuthenticatorRepository.EXPECT().GetJWSFromRequest(gomock.Any()).Return("", nil)
				f.AuthenticatorRepository.EXPECT().ParseWithKeySet(gomock.Any()).Return(nil, nil)
				f.AuthenticatorRepository.EXPECT().ValidateToken(gomock.Any(), gomock.Any()).Return(nil)
			},
			args:    args{},
			want:    nil,
			wantErr: assert.NoError,
		},
		{
			name: "Case #2 no error when verification is disabled",
			prepare: func(f *fields) {
				os.Setenv("JWT_VERIFY", "false")
				f.AuthenticatorRepository.EXPECT().GetJWSFromRequest(gomock.Any()).Return("", nil)
				f.AuthenticatorRepository.EXPECT().ParseWithoutVerify(gomock.Any()).Return(nil, nil)
			},
			args:    args{},
			want:    nil,
			wantErr: assert.NoError,
		},
		{
			name: "Case #3 error when verification is disabled due to parse error",
			prepare: func(f *fields) {
				os.Setenv("JWT_VERIFY", "false")
				f.AuthenticatorRepository.EXPECT().GetJWSFromRequest(gomock.Any()).Return("", nil)
				f.AuthenticatorRepository.EXPECT().ParseWithoutVerify(gomock.Any()).Return(nil, errors.New(""))
			},
			args:    args{},
			want:    nil,
			wantErr: assert.Error,
		},
		{
			name: "Case #3 error due to validation error",
			prepare: func(f *fields) {
				os.Setenv("JWT_VERIFY", "true")
				f.AuthenticatorRepository.EXPECT().GetJWSFromRequest(gomock.Any()).Return("", nil)
				f.AuthenticatorRepository.EXPECT().ParseWithKeySet(gomock.Any()).Return(nil, nil)
				f.AuthenticatorRepository.EXPECT().ValidateToken(gomock.Any(), gomock.Any()).Return(errors.New(""))
			},
			args:    args{},
			want:    nil,
			wantErr: assert.Error,
		},
		{
			name: "Case #4 error due to parse error",
			prepare: func(f *fields) {
				os.Setenv("JWT_VERIFY", "true")
				f.AuthenticatorRepository.EXPECT().GetJWSFromRequest(gomock.Any()).Return("", nil)
				f.AuthenticatorRepository.EXPECT().ParseWithKeySet(gomock.Any()).Return(nil, errors.New(""))
			},
			args:    args{},
			want:    nil,
			wantErr: assert.Error,
		},
		{
			name: "Case #5 error due to wrong request",
			prepare: func(f *fields) {
				f.AuthenticatorRepository.EXPECT().GetJWSFromRequest(gomock.Any()).Return("", errors.New(""))
			},
			args:    args{},
			want:    nil,
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := donkeytest.StartTest(t)
			defer controller.Finish()

			f := fields{
				AuthenticatorRepository: vt.NewMockAuthenticator(controller),
			}
			if tt.prepare != nil {
				tt.prepare(&f)
			}
			j := AuthenticatorServiceImpl{
				authenticator: f.AuthenticatorRepository,
			}

			got, err := j.ParseJWT(tt.args.req)
			if !tt.wantErr(t, err, fmt.Sprintf("ParseJWT(%v)", tt.args.req)) {
				return
			}
			assert.Equalf(t, tt.want, got, "ParseJWT(%v)", tt.args.req)
		})
	}
}
