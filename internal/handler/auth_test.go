package handler

import (
	"bytes"
	"database/sql"
	"encoding/json"
	db "github.com/Bakhram74/amazon-backend.git/db/sqlc"
	"github.com/Bakhram74/amazon-backend.git/internal/config"
	"github.com/Bakhram74/amazon-backend.git/internal/service"
	mock_service "github.com/Bakhram74/amazon-backend.git/internal/service/mocks"
	"github.com/Bakhram74/amazon-backend.git/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

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
				"name":     user.Name,
				"email":    user.Email,
				"password": password,
			},
			buildStubs: func(r *mock_service.MockAuthorization) {
				arg := db.CreateUserParams{
					Name:     user.Name,
					Email:    user.Email,
					Password: password,
				}

				r.EXPECT().
					CreateUser(gomock.Any(), gomock.Eq(arg)).
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
				"name":     user.Name,
				"email":    user.Email,
				"password": password,
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
				"name":     "invalid-user#1",
				"email":    user.Email,
				"password": password,
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
				"name":     user.Name,
				"email":    user.Email,
				"password": "123",
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

func TestLoginUserAPI(t *testing.T) {
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
				"email":    user.Email,
				"password": password,
			},
			buildStubs: func(r *mock_service.MockAuthorization) {
				r.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(user.Email)).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "UserNotFound",
			body: gin.H{
				"email":    "NotFound@list.ru",
				"password": password,
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
			name: "InternalError",
			body: gin.H{
				"email":    user.Email,
				"password": password,
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
			name: "InvalidEmail",
			body: gin.H{
				"email":    "invalid-email",
				"password": password,
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
				"password": "incorrect",
				"email":    user.Email,
			},
			buildStubs: func(r *mock_service.MockAuthorization) {
				r.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(user.Email)).
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
			config := config.Config{
				TokenSymmetricKey:   utils.RandomString(32),
				AccessTokenDuration: time.Minute,
			}
			handler, _ := NewHandler(config, s)

			recorder := httptest.NewRecorder()

			// Marshal body data to JSON
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			r := gin.New()
			r.POST("/sign-in", handler.signIn)
			//
			req, err := http.NewRequest(http.MethodPost, "/sign-in", bytes.NewReader(data))
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
		Name:       utils.RandomString(6),
		Email:      utils.RandomEmail(),
		Phone:      "+7(###) ###-###-###",
		Password:   hashedPassword,
		AvatarPath: "https://cojo.ru/wp-content/uploads/2022/12/avatarki-dlia-vatsapa-3.webp",
	}
	return
}
