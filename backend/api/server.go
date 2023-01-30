package api

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/mearaj/bhagad-house-booking/backend"
	"github.com/mearaj/bhagad-house-booking/common/alog"
	"github.com/mearaj/bhagad-house-booking/common/model"
	"github.com/mearaj/bhagad-house-booking/common/request"
	"github.com/mearaj/bhagad-house-booking/common/response"
	"github.com/mearaj/bhagad-house-booking/common/token"
	"github.com/mearaj/bhagad-house-booking/common/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Server struct {
	config     backend.Config
	client     *mongo.Client
	router     *gin.Engine
	tokenMaker token.Maker
}

func NewServer(config backend.Config, client *mongo.Client) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker:%w", err)
	}
	if client == nil {
		return nil, errors.New("mongo client is nil")
	}
	server := &Server{
		config:     config,
		client:     client,
		tokenMaker: tokenMaker,
	}
	server.setupCollections()
	server.setupAdmin()
	server.setupRouter()
	return server, nil
}

func (s *Server) setupCollections() {
	database := s.client.Database("bhagad_house_booking")
	usersCollection = database.Collection("users")
	bookingsCollection = database.Collection("bookings")
	transactionsCollection = database.Collection("transactions")
}

// setupAdmin for security reasons, password is removed after the function is returned
func (s *Server) setupAdmin() {
	rsp := response.NewUser{}
	defer func() {
		s.config.AdminPassword = ""
		rsp = response.NewUser{}
	}()
	if s.config.AdminEmail == "" {
		return
	}
	if len(s.config.AdminPassword) < 6 {
		return
	}
	// check if user exists, error is expected if it doesn't
	err := usersCollection.FindOne(context.TODO(), bson.M{"email": s.config.AdminEmail}).Decode(&rsp.User)
	if err != nil {
		alog.Logger().Errorln(err)
	}
	hashedPassword, err := utils.HashPassword(s.config.AdminPassword)
	if err != nil {
		alog.Logger().Errorln(err)
		return
	}
	rq := request.CreateUser{
		Name:     s.config.AdminName,
		Email:    s.config.AdminEmail,
		Password: hashedPassword,
		Roles:    []model.UserRoles{model.UserRolesUser, model.UserRolesAdmin},
	}
	_, err = usersCollection.InsertOne(context.TODO(), rq)
	if err != nil {
		alog.Logger().Errorln(err)
		return
	}
	err = usersCollection.FindOne(context.TODO(), bson.M{"email": s.config.AdminEmail}).Decode(&rsp.User)
	if err != nil {
		alog.Logger().Errorln(err)
		return
	}
	alog.Logger().Println("Successfully created admin user")
}

func (s *Server) setupRouter() {
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders: []string{"Content-Type,access-control-allow-origin, access-control-allow-headers,Authorization"},
	}))
	router.POST("/users/login", s.loginUser)
	router.GET("/bookings", s.getBookings)
	authRoutes := router.Group("/").Use(authMiddleWare(s.tokenMaker))
	authRoutes.POST("/bookings", s.createBooking)
	authRoutes.PUT("/bookings", s.updateBooking)
	authRoutes.DELETE("/bookings", s.deleteBooking)
	authRoutes.GET("/bookings/search", s.searchBookings)
	authRoutes.GET("/bookings/:booking_id/transactions", s.getTransactions)
	authRoutes.POST("/transactions", s.addUpdateTransaction)
	authRoutes.DELETE("/transactions", s.deleteTransaction)
	s.router = router
}

// Start runs the HTTP server on a specific address
func (s *Server) Start() error {
	return s.router.Run(fmt.Sprintf(":%s", s.config.ServerPort))
}
