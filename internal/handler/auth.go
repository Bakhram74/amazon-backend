package handler

import (
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
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,

		Phone:      user.Phone,
		AvatarPath: user.AvatarPath,
		UpdatedAt:  user.UpdatedAt,
		CreatedAt:  user.CreatedAt,
	}
}
func (h *Handler) signUp(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse("Invalid json provided", err))
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse("Server error", err))
		return
	}

	arg := db.CreateUserParams{
		Name:           req.Name,
		Email:          req.Email,
		HashedPassword: hashedPassword,
	}

	user, err := h.services.CreateUser(ctx, arg)
	if err != nil {
		if db.ErrorCode(err) == db.UniqueViolation {
			ctx.JSON(http.StatusForbidden, err)
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse("Server error", err))
		return
	}

	rsp := newUserResponse(user)
	ctx.JSON(http.StatusOK, rsp)
}

//type loginUserRequest struct {
//	PhoneNumber string `json:"phone_number" binding:"required"`
//	Password    string `json:"password" binding:"required,min=6"`
//}
//
//type loginUserResponse struct {
//	SessionID             uuid.UUID    `json:"session_id"`
//	AccessToken           string       `json:"access_token"`
//	AccessTokenExpiresAt  time.Time    `json:"access_token_expires_at"`
//	RefreshToken          string       `json:"refresh_token"`
//	RefreshTokenExpiresAt time.Time    `json:"refresh_token_expires_at"`
//	User                  userResponse `json:"user"`
//}

//func (h *Handler) signIn(ctx *gin.Context) {
//	var req loginUserRequest
//	if err := ctx.ShouldBindJSON(&req); err != nil {
//		ctx.JSON(http.StatusBadRequest, errorResponse("Invalid json provided", err))
//		return
//	}
//
//	user, err := h.services.GetUser(ctx, req.PhoneNumber)
//	if err != nil {
//		if errors.Is(err, db.ErrRecordNotFound) {
//			ctx.JSON(http.StatusNotFound, errorResponse("User not found", err))
//			return
//		}
//		ctx.JSON(http.StatusInternalServerError, errorResponse("Server error", err))
//		return
//	}
//
//	err = utils.CheckPassword(req.Password, user.HashedPassword)
//	if err != nil {
//		ctx.JSON(http.StatusUnauthorized, errorResponse("Please provide valid login details", err))
//		return
//	}
//	accessToken, accessPayload, err := h.tokenMaker.CreateToken(
//		user.ID,
//		h.config.AccessTokenDuration,
//	)
//	if err != nil {
//		ctx.JSON(http.StatusInternalServerError, errorResponse("Server error", err))
//		return
//	}
//	refreshToken, refreshPayload, err := h.tokenMaker.CreateToken(
//		user.ID,
//		h.config.RefreshTokenDuration,
//	)
//	if err != nil {
//		ctx.JSON(http.StatusInternalServerError, errorResponse("Server error", err))
//		return
//	}
//	session, err := h.services.CreateSession(ctx, db.CreateSessionParams{
//		ID:           refreshPayload.ID,
//		Userid:       user.ID,
//		RefreshToken: refreshToken,
//		UserAgent:    ctx.Request.UserAgent(),
//		ClientIp:     ctx.ClientIP(),
//		IsBlocked:    false,
//		ExpiresAt:    refreshPayload.ExpiredAt,
//	})
//	rsp := loginUserResponse{
//		SessionID:             session.ID,
//		AccessToken:           accessToken,
//		AccessTokenExpiresAt:  accessPayload.ExpiredAt,
//		RefreshToken:          refreshToken,
//		RefreshTokenExpiresAt: refreshPayload.ExpiredAt,
//		User:                  newUserResponse(user),
//	}
//	ctx.JSON(http.StatusOK, rsp)
//}

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
