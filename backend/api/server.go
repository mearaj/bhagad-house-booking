package api

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/mearaj/bhagad-house-booking/backend"
	"github.com/mearaj/bhagad-house-booking/common/db/sqlc"
	"github.com/mearaj/bhagad-house-booking/common/token"
)

type Server struct {
	config     backend.Config
	store      sqlc.Store
	router     *gin.Engine
	tokenMaker token.Maker
}

func NewServer(config backend.Config, store sqlc.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker:%w", err)
	}
	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}
	server.setupRouter()
	return server, nil
}

func (s *Server) setupRouter() {
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders: []string{"Content-Type,access-control-allow-origin, access-control-allow-headers,Authorization"},
	}))
	router.POST("/users/login", s.loginUser)
	router.GET("/bookings", s.getBookings)
	authRoutes := router.Group("/").Use(authMiddleWare(s.tokenMaker))
	authRoutes.POST("/bookings", s.createBooking)
	authRoutes.PUT("/bookings", s.updateBooking)
	authRoutes.DELETE("/bookings", s.deleteBooking)
	authRoutes.GET("/bookings/search", s.searchBookings)
	s.router = router
}

// Start runs the HTTP server on a specific address
func (s *Server) Start() error {
	return s.router.Run(fmt.Sprintf(":%s", s.config.ServerPort))
}
