package api

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"github.com/mearaj/bhagad-house-booking/common/db/sqlc"
	"github.com/mearaj/bhagad-house-booking/common/utils"
	"net/http"
)

type createUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

func (s *Server) createUser(ctx *gin.Context) {
	var req createUserRequest
	var resp sqlc.NewUserResponse
	if err := ctx.ShouldBindJSON(&req); err != nil {
		resp.Error = err.Error()
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		resp.Error = err.Error()
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}
	arg := sqlc.CreateUserParams{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPassword,
	}
	user, err := s.store.CreateUser(ctx, arg)
	if err != nil {
		resp.Error = err.Error()
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, resp)
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}
	user.Password = ""
	resp.User = user
	ctx.JSON(http.StatusOK, resp)
}

type getUserRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (s *Server) getUser(ctx *gin.Context) {
	var req getUserRequest
	var resp sqlc.UserResponse
	if err := ctx.ShouldBindUri(&req); err != nil {
		resp.Error = err.Error()
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}
	user, err := s.store.GetUserByID(ctx, req.ID)
	if err != nil {
		resp.Error = err.Error()
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, resp)
		}
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}
	user.Password = ""
	resp.User = user
	ctx.JSON(http.StatusOK, resp)
}

type listUsersRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=1"`
}

func (s *Server) listUsers(ctx *gin.Context) {
	var req listUsersRequest
	var resp sqlc.UsersResponse
	if err := ctx.ShouldBindQuery(&req); err != nil {
		resp.Error = err.Error()
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}
	arg := sqlc.ListUsersParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}
	users, err := s.store.ListUsers(ctx, arg)
	if err != nil {
		resp.Error = err.Error()
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}
	for _, user := range users {
		user.Password = ""
	}
	resp.Users = users
	ctx.JSON(http.StatusOK, resp)
}

type loginUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

func (s *Server) loginUser(ctx *gin.Context) {
	var req loginUserRequest
	var resp sqlc.LoginUserResponse
	if err := ctx.ShouldBindJSON(&req); err != nil {
		resp.Error = err.Error()
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}
	user, err := s.store.GetUserByEmail(ctx, req.Email)
	if err != nil {
		resp.Error = err.Error()
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, resp)
			return
		}
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}
	err = utils.CheckPassword(req.Password, user.Password)
	if err != nil {
		resp.Error = err.Error()
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}

	accessToken, err := s.tokenMaker.CreateToken(user.Email, s.config.AccessTokenDuration)
	if err != nil {
		resp.Error = err.Error()
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}
	user.Password = ""
	resp.User = user
	resp.AccessToken = accessToken
	ctx.JSON(http.StatusOK, resp)
}
