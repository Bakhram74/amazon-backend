package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	db "github.com/Bakhram74/amazon-backend.git/db/sqlc"
	"github.com/Bakhram74/amazon-backend.git/internal/config"
	"github.com/Bakhram74/amazon-backend.git/internal/service"
	mock_service "github.com/Bakhram74/amazon-backend.git/internal/service/mocks"
	"github.com/Bakhram74/amazon-backend.git/pkg/utils"

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
				"name":     user.Name,
				"email":    user.Email,
				"password": password,
			},
			buildStubs: func(r *mock_service.MockAuthorization) {
				arg := db.CreateUserParams{
					Name:           user.Name,
					Email:          user.Email,
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
		//
		//{
		//	name: "InternalError",
		//	body: gin.H{
		//		"name":     user.Name,
		//		"email":    user.Email,
		//		"password": password,
		//	},
		//	buildStubs: func(r *mock_service.MockAuthorization) {
		//		r.EXPECT().
		//			CreateUser(gomock.Any(), gomock.Any()).
		//			Times(1).
		//			Return(db.User{}, sql.ErrConnDone)
		//	},
		//	checkResponse: func(recorder *httptest.ResponseRecorder) {
		//		require.Equal(t, http.StatusInternalServerError, recorder.Code)
		//	},
		//},
		//{
		//	name: "InvalidUsername",
		//	body: gin.H{
		//		"name":     "invalid-user#1",
		//		"email":    user.Email,
		//		"password": password,
		//	},
		//	buildStubs: func(r *mock_service.MockAuthorization) {
		//		r.EXPECT().
		//			CreateUser(gomock.Any(), gomock.Any()).
		//			Times(0)
		//	},
		//	checkResponse: func(recorder *httptest.ResponseRecorder) {
		//		require.Equal(t, http.StatusBadRequest, recorder.Code)
		//	},
		//},
		//{
		//	name: "TooShortPassword",
		//	body: gin.H{
		//		"name":     user.Name,
		//		"email":    user.Email,
		//		"password": "123",
		//	},
		//	buildStubs: func(r *mock_service.MockAuthorization) {
		//		r.EXPECT().
		//			CreateUser(gomock.Any(), gomock.Any()).
		//			Times(0)
		//	},
		//	checkResponse: func(recorder *httptest.ResponseRecorder) {
		//		require.Equal(t, http.StatusBadRequest, recorder.Code)
		//	},
		//},
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

func randomUser(t *testing.T) (user db.User, password string) {
	password = utils.RandomString(6)
	hashedPassword, err := utils.HashPassword(password)
	require.NoError(t, err)

	user = db.User{
		Name:           utils.RandomString(6),
		Email:          utils.RandomEmail(),
		Phone:          "+7(###) ###-###-###",
		HashedPassword: hashedPassword,
		AvatarPath:     "https://cojo.ru/wp-content/uploads/2022/12/avatarki-dlia-vatsapa-3.webp",
	}
	return
}

func requireBodyMatchUser(t *testing.T, body *bytes.Buffer, user db.User) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotUser db.User
	err = json.Unmarshal(data, &gotUser)

	require.NoError(t, err)
	require.Equal(t, user.Name, gotUser.Name)
	require.Equal(t, user.Phone, gotUser.Phone)
	require.Equal(t, user.Email, gotUser.Email)
	require.Empty(t, gotUser.HashedPassword)
}
