package api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/techschool/simplebank/token"

	"github.com/gin-gonic/gin"
)

const (
	authorizationHeaderKey = "authorization"
	// assume that now only support Bearer token
	authorizationTypeBearer = "bearer"
	// key of the payload
	authorizationPayloadKey = "authorization_payload"
)

// this func is not the middleware
// return the Authentication Middleware Function
func authMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	// gin.HandlerFunc is in fact a *gin.Context

	// this is the authentication middleware func
	return func(ctx *gin.Context) {
		// get the header from the context
		authorizationHeader := ctx.GetHeader(authorizationHeaderKey)
		if len(authorizationHeader) == 0 {
			err := errors.New("authorization header is not provided")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		// The bearer prefix make the server know the type of authorization
		// in reality server support many kinds of authorization schema:
		// OAuth, AWS signature, Digest, ...

		// return a []string with each word of the header
		fields := strings.Fields(authorizationHeader)

		// at least have 2 elems Bearer and token
		if len(fields) < 2 {
			err := errors.New("invalid authorization header format")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		// 1st elems should be the authorization Type
		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationTypeBearer {
			err := fmt.Errorf("unsupported authorization type: %s", authorizationType)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		// 2nd elems is the token
		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		ctx.Set(authorizationPayloadKey, payload)
		// handle the ctx to the next handler
		ctx.Next()
	}
}
