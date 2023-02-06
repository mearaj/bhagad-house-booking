package api

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/mearaj/bhagad-house-booking/common/model"
	"github.com/mearaj/bhagad-house-booking/common/request"
	"github.com/mearaj/bhagad-house-booking/common/response"
	"github.com/mearaj/bhagad-house-booking/common/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

func (s *Server) createUser(ctx *gin.Context) {
	var rq request.CreateUser
	var rsp response.NewUser
	if err := ctx.ShouldBindJSON(&rq); err != nil {
		rsp.Error = err.Error()
		ctx.JSON(http.StatusBadRequest, rsp)
		return
	}
	hashedPassword, err := utils.HashPassword(rq.Password)
	if err != nil {
		rsp.Error = err.Error()
		ctx.JSON(http.StatusInternalServerError, rsp)
		return
	}
	rq.Password = hashedPassword
	result, err := usersCollection.InsertOne(context.TODO(), rq)
	if err != nil {
		rsp.Error = err.Error()
		ctx.JSON(http.StatusInternalServerError, rsp)
		return
	}
	id := result.InsertedID.(primitive.ObjectID)
	err = usersCollection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&rsp.User)
	if err != nil {
		rsp.Error = err.Error()
		ctx.JSON(http.StatusInternalServerError, rsp)
		return
	}
	ctx.JSON(http.StatusOK, rsp)
}

func (s *Server) loginUser(ctx *gin.Context) {
	var rq request.LoginUser
	var rsp response.LoginUser
	var user model.User
	if err := ctx.ShouldBindJSON(&rq); err != nil {
		rsp.Error = err.Error()
		ctx.JSON(http.StatusBadRequest, rsp)
		return
	}

	err := usersCollection.FindOne(context.TODO(), bson.M{"email": rq.Email}).Decode(&user)
	if err != nil {
		rsp.Error = err.Error()
		ctx.JSON(http.StatusInternalServerError, rsp)
		return
	}
	err = utils.CheckPassword(rq.Password, user.Password)
	if err != nil {
		rsp.Error = err.Error()
		ctx.JSON(http.StatusInternalServerError, rsp)
		return
	}

	accessToken, err := s.tokenMaker.CreateToken(user.Email, s.config.AccessTokenDuration)
	if err != nil {
		rsp.Error = err.Error()
		ctx.JSON(http.StatusInternalServerError, rsp)
		return
	}
	payload, err := s.tokenMaker.VerifyToken(accessToken)
	if err != nil {
		rsp.Error = err.Error()
		ctx.JSON(http.StatusInternalServerError, rsp)
		return
	}
	rsp.User = response.User{Name: user.Name, Email: user.Email, ID: user.ID, Roles: user.Roles}
	rsp.AccessToken = accessToken
	rsp.ExpiresAt = payload.ExpiredAt
	ctx.JSON(http.StatusOK, rsp)
}
