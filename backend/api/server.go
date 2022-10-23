package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mearaj/bhagad-house-booking/common/db/sqlc"
	"github.com/mearaj/bhagad-house-booking/common/token"
	"github.com/mearaj/bhagad-house-booking/common/utils"
)

type Server struct {
	config     utils.Config
	store      sqlc.Store
	router     *gin.Engine
	tokenMaker token.Maker
}

func NewServer(config utils.Config, store sqlc.Store) (*Server, error) {
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
	router.POST("/users", s.createUser)
	router.POST("/users/login", s.loginUser)
	authRoutes := router.Group("/").Use(authMiddleWare(s.tokenMaker))

	authRoutes.GET("/users/:id", s.getUser)

	authRoutes.POST("/customers", s.createCustomer)
	authRoutes.GET("/customers/:id", s.getCustomer)
	authRoutes.GET("/customers", s.listCustomers)
	s.router = router
}

// Start runs the HTTP server on a specific address
func (s *Server) Start(address string) error {
	return s.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
