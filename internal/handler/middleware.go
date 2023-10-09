package handler

import (
	"errors"
	"github.com/Bakhram74/amazon-backend.git/pkg/token"

	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

const (
	authorizationHeaderKey  = "Authorization"
	authorizationTypeBearer = "Bearer"
	userCtx                 = "userId"
	authorizationPayloadKey = "authorization_payload"
)

func authMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(authorizationHeaderKey)

		if len(authorizationHeader) == 0 {
			err := errors.New("authorization header is not provided")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, err)
			return
		}

		fields := strings.Split(authorizationHeader, " ")
		if len(fields) != 2 || fields[0] != "Bearer" {
			err := errors.New("invalid authorization header format")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, err)
			return
		}

		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, err.Error())
			return
		}

		ctx.Set(userCtx, payload.Userid)
		ctx.Next()
	}
}

func getUserId(c *gin.Context) (int64, error) {
	id, ok := c.Get(userCtx)
	if !ok {
		return 0, errors.New("user id not found")
	}

	idInt64, ok := id.(int64)
	if !ok {
		return 0, errors.New("user id is of invalid type")
	}

	return idInt64, nil
}
