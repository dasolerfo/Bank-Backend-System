package api

import (
	"errors"
	"net/http"
	token "simplebank/token"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	authorizationKey        = "authorization"
	authTypeBearer          = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

func AuthMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		authHeader := ctx.GetHeader(authorizationKey)
		if authHeader == "" {
			//ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			//ctx.Abort()
			err := errors.New("authorization header is required")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		// Validate the token (this is a placeholder, implement your own logic)

		//token := authHeader[len("Bearer "):]

		fields := strings.Fields(authHeader)
		if len(fields) < 2 {
			err := errors.New("invalid authorization header format")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		authType := strings.ToLower(fields[0])
		if authTypeBearer != authType {
			err := errors.New("invalid authorization header type")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		ctx.Set(authorizationPayloadKey, payload)
		ctx.Next()
	}
}
