package handler

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	db "github.com/Bakhram74/amazon.git/db/sqlc"
	"github.com/Bakhram74/amazon.git/internal/config"
	"github.com/Bakhram74/amazon.git/internal/service"
	mock_service "github.com/Bakhram74/amazon.git/internal/service/mocks"
	"github.com/Bakhram74/amazon.git/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

type eqCreateUserParamsMatcher struct {
	arg      db.CreateUserParams
	password string
}

func (e eqCreateUserParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(db.CreateUserParams)
	if !ok {
		return false
	}

	err := utils.CheckPassword(e.password, arg.HashedPassword)
	if err != nil {
		return false
	}

	e.arg.HashedPassword = arg.HashedPassword
	return reflect.DeepEqual(e.arg, arg)
}

func (e eqCreateUserParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %v", e.arg, e.password)
}

func EqCreateUserParams(arg db.CreateUserParams, password string) gomock.Matcher {
	return eqCreateUserParamsMatcher{arg, password}
}

func TestCreateUserAPI(t *testing.T) {
	user, password := randomUser(t)

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(r *mock_service.MockAuthorization)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"username":     user.Username,
				"password":     password,
				"phone_number": user.PhoneNumber,
			},
			buildStubs: func(r *mock_service.MockAuthorization) {
				arg := db.CreateUserParams{
					Username:       user.Username,
					PhoneNumber:    user.PhoneNumber,
					HashedPassword: password,
				}
				r.EXPECT().
					CreateUser(gomock.Any(), EqCreateUserParams(arg, password)).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchUser(t, recorder.Body, user)
			},
		},

		{
			name: "InternalError",
			body: gin.H{
				"username":     user.Username,
				"password":     password,
				"phone_number": user.PhoneNumber,
			},
			buildStubs: func(r *mock_service.MockAuthorization) {
				r.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InvalidUsername",
			body: gin.H{
				"username":     "invalid-user#1",
				"password":     password,
				"phone_number": user.PhoneNumber,
			},
			buildStubs: func(r *mock_service.MockAuthorization) {
				r.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "TooShortPassword",
			body: gin.H{
				"username":     user.Username,
				"password":     "123",
				"phone_number": user.PhoneNumber,
			},
			buildStubs: func(r *mock_service.MockAuthorization) {
				r.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mock_service.NewMockAuthorization(ctrl)
			tc.buildStubs(repo)

			s := &service.Service{Authorization: repo}
			config := config.Config{
				TokenSymmetricKey:   utils.RandomString(32),
				AccessTokenDuration: time.Minute,
			}
			handler, _ := NewHandler(config, s)
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)
			r := gin.New()
			r.POST("/sign-up", handler.signUp)
			req := httptest.NewRequest("POST", "/sign-up",
				bytes.NewReader(data))
			require.NoError(t, err)

			r.ServeHTTP(recorder, req)
			tc.checkResponse(recorder)
		})
	}
}

func TestGetUserAPI(t *testing.T) {
	user, password := randomUser(t)

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(r *mock_service.MockAuthorization)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"phone_number": user.PhoneNumber,
				"password":     password,
			},
			buildStubs: func(r *mock_service.MockAuthorization) {
				r.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(user.PhoneNumber)).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},

		{
			name: "InternalError",
			body: gin.H{
				"password":     password,
				"phone_number": user.PhoneNumber,
			},
			buildStubs: func(r *mock_service.MockAuthorization) {
				r.EXPECT().
					GetUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "PhoneNumberNotFound",
			body: gin.H{
				"phone_number": "12345",
				"password":     password,
			},
			buildStubs: func(r *mock_service.MockAuthorization) {
				r.EXPECT().
					GetUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, db.ErrRecordNotFound)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "TooShortPassword",
			body: gin.H{
				"password":     "123",
				"phone_number": user.PhoneNumber,
			},
			buildStubs: func(r *mock_service.MockAuthorization) {
				r.EXPECT().
					GetUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "IncorrectPassword",
			body: gin.H{
				"password":     "incorrect",
				"phone_number": user.PhoneNumber,
			},
			buildStubs: func(r *mock_service.MockAuthorization) {
				r.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(user.PhoneNumber)).
					Times(1).Return(db.User{}, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mock_service.NewMockAuthorization(ctrl)
			tc.buildStubs(repo)

			s := &service.Service{Authorization: repo}
			cfg := config.Config{
				TokenSymmetricKey:   utils.RandomString(32),
				AccessTokenDuration: time.Minute,
			}
			handler, _ := NewHandler(cfg, s)
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)
			r := gin.New()
			r.POST("/sign-in", handler.signIn)
			req := httptest.NewRequest("POST", "/sign-in",
				bytes.NewReader(data))
			require.NoError(t, err)

			r.ServeHTTP(recorder, req)
			tc.checkResponse(recorder)
		})
	}
}

func randomUser(t *testing.T) (user db.User, password string) {
	password = utils.RandomString(6)
	hashedPassword, err := utils.HashPassword(password)
	require.NoError(t, err)

	user = db.User{
		Username:       utils.RandomString(6),
		HashedPassword: hashedPassword,
		PhoneNumber:    utils.RandomNumbers(9),
	}
	return
}

func requireBodyMatchUser(t *testing.T, body *bytes.Buffer, user db.User) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotUser db.User
	err = json.Unmarshal(data, &gotUser)

	require.NoError(t, err)
	require.Equal(t, user.Username, gotUser.Username)
	require.Equal(t, user.PhoneNumber, gotUser.PhoneNumber)
	require.Empty(t, gotUser.HashedPassword)
}
