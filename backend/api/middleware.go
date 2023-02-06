package api

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/mearaj/bhagad-house-booking/common/model"
	"github.com/mearaj/bhagad-house-booking/common/response"
	"github.com/mearaj/bhagad-house-booking/common/token"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"strings"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

// checkUserAuthorized httpStatus with value 0 also indicates status ok
func checkUserAuthorized(ctx *gin.Context, tokenMaker token.Maker) (payload *token.Payload, httpStatus int, err error) {
	authorizationHeader := ctx.GetHeader(authorizationHeaderKey)
	if len(authorizationHeader) == 0 {
		err = errors.New("authorization header is not provided")
		return nil, http.StatusUnauthorized, err
	}
	fields := strings.Fields(authorizationHeader)
	if len(fields) < 2 {
		return nil, http.StatusUnauthorized, err
	}
	authorizationType := strings.ToLower(fields[0])
	if authorizationType != authorizationTypeBearer {
		return nil, http.StatusUnauthorized, err
	}
	accessToken := fields[1]
	payload, err = tokenMaker.VerifyToken(accessToken)
	if err != nil {
		return nil, http.StatusUnauthorized, err
	}
	var user model.User
	err = usersCollection.FindOne(context.TODO(), bson.D{{Key: "email", Value: payload.UniqueStr}}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, http.StatusUnauthorized, errors.New("unauthorized")
		}
		return nil, http.StatusInternalServerError, err
	}
	var adminFound bool
	for _, role := range user.Roles {
		if role == model.UserRolesAdmin {
			adminFound = true
			break
		}
	}
	if !adminFound {
		return nil, http.StatusUnauthorized, errors.New("unauthorized")
	}
	return payload, http.StatusOK, nil
}

func authMiddleWare(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		payload, status, err := checkUserAuthorized(ctx, tokenMaker)
		if err != nil {
			ctx.AbortWithStatusJSON(status, response.AuthError{Error: err.Error()})
			return
		}
		ctx.Set(authorizationPayloadKey, payload)
		ctx.Next()
	}
}
