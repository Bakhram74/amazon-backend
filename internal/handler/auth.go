package handler

import (
	"errors"
	db "github.com/Bakhram74/amazon-backend.git/db/sqlc"
	"github.com/Bakhram74/amazon-backend.git/pkg/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type createUserRequest struct {
	Name     string `json:"name" binding:"required,alphanum"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

func (h *Handler) signUp(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	arg := db.CreateUserParams{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}

	user, err := h.services.CreateUser(ctx, arg)
	if err != nil {
		if db.ErrorCode(err) == db.UniqueViolation {
			ctx.JSON(http.StatusForbidden, err)
			return
		}
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	accessToken, _, err := h.tokenMaker.CreateToken(
		user.ID,
		h.config.AccessTokenDuration,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	refreshToken, _, err := h.tokenMaker.CreateToken(
		user.ID,
		h.config.RefreshTokenDuration,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	rsp := userData{
		User:         newUserResponse(user),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	ctx.JSON(http.StatusOK, rsp)
}

type loginUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

func (h *Handler) signIn(ctx *gin.Context) {
	var req loginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.services.GetUser(ctx, req.Email)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, err.Error())
			return
		}
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	err = utils.CheckPassword(req.Password, user.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, err.Error())
		return
	}

	accessToken, _, err := h.tokenMaker.CreateToken(
		user.ID,
		h.config.AccessTokenDuration,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	refreshToken, _, err := h.tokenMaker.CreateToken(
		user.ID,
		h.config.RefreshTokenDuration,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	rsp := userData{
		User:         newUserResponse(user),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	ctx.JSON(http.StatusOK, rsp)
}

//type renewAccessTokenRequest struct {
//	RefreshToken string `json:"refresh_token" binding:"required"`
//}
//
//type renewAccessTokenResponse struct {
//	AccessToken          string    `json:"access_token"`
//	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
//}

//func (h *Handler) renewAccessToken(ctx *gin.Context) {
//	var req renewAccessTokenRequest
//	if err := ctx.ShouldBindJSON(&req); err != nil {
//		ctx.JSON(http.StatusBadRequest, errorResponse("Invalid json provided", err))
//		return
//	}
//
//	refreshPayload, err := h.tokenMaker.VerifyToken(req.RefreshToken)
//	if err != nil {
//		ctx.JSON(http.StatusUnauthorized, errorResponse("Invalid refresh token", err))
//		return
//	}
//
//	session, err := h.services.GetSession(ctx, refreshPayload.ID)
//	if err != nil {
//		if errors.Is(err, db.ErrRecordNotFound) {
//			ctx.JSON(http.StatusNotFound, errorResponse("Incorrect id", err))
//			return
//		}
//		ctx.JSON(http.StatusInternalServerError, errorResponse("Server error", err))
//		return
//	}
//
//	if session.IsBlocked {
//		err := fmt.Errorf("blocked session")
//		ctx.JSON(http.StatusUnauthorized, errorResponse("", err))
//		return
//	}
//
//	if session.Userid != refreshPayload.Userid {
//		err := fmt.Errorf("incorrect session user")
//		ctx.JSON(http.StatusUnauthorized, errorResponse("", err))
//		return
//	}
//
//	if session.RefreshToken != req.RefreshToken {
//		err := fmt.Errorf("mismatched session token")
//		ctx.JSON(http.StatusUnauthorized, errorResponse("", err))
//		return
//	}
//
//	if time.Now().After(session.ExpiresAt) {
//		err := fmt.Errorf("expired session")
//		ctx.JSON(http.StatusUnauthorized, errorResponse("", err))
//		return
//	}
//
//	accessToken, accessPayload, err := h.tokenMaker.CreateToken(
//		refreshPayload.Userid,
//		h.config.AccessTokenDuration,
//	)
//	if err != nil {
//		ctx.JSON(http.StatusInternalServerError, errorResponse("Server error", err))
//		return
//	}
//
//	rsp := renewAccessTokenResponse{
//		AccessToken:          accessToken,
//		AccessTokenExpiresAt: accessPayload.ExpiredAt,
//	}
//	ctx.JSON(http.StatusOK, rsp)
//}

type userResponse struct {
	ID         int64     `json:"id"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	Phone      string    `json:"phone"`
	AvatarPath string    `json:"avatar_path"`
	UpdatedAt  time.Time `json:"updated_at"`
	CreatedAt  time.Time `json:"created_at"`
}

func newUserResponse(user db.User) userResponse {
	return userResponse{
		ID:         user.ID,
		Name:       user.Name,
		Email:      user.Email,
		Phone:      user.Phone,
		AvatarPath: user.AvatarPath,
		UpdatedAt:  user.UpdatedAt,
		CreatedAt:  user.CreatedAt,
	}
}

type userData struct {
	User         userResponse `json:"user"`
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
}
