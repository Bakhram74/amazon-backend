package handler

import (
	db "github.com/Bakhram74/amazon-backend.git/db/sqlc"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) getProfile(ctx *gin.Context) {
	userId, err := getUserId(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	user, err := h.services.GetUserByID(ctx, userId)

	rsp := db.User{
		ID:         user.ID,
		Name:       user.Name,
		Email:      user.Email,
		Phone:      user.Phone,
		AvatarPath: user.AvatarPath,
		UpdatedAt:  user.UpdatedAt,
		CreatedAt:  user.CreatedAt,
	}
	ctx.JSON(http.StatusOK, rsp)
}
