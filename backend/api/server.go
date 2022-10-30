package api

import (
	"fmt"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/mearaj/bhagad-house-booking/common/db/sqlc"
	"github.com/mearaj/bhagad-house-booking/common/token"
	"github.com/mearaj/bhagad-house-booking/common/utils"
	log "github.com/sirupsen/logrus"
	"path/filepath"
	"runtime"
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
	_, p, _, ok := runtime.Caller(0) // provides path of this main file
	if !ok {
		log.Fatalln("error in runtime.Caller, cannot load path")
	}
	p = filepath.Join(p, filepath.FromSlash("../../dist/index.html"))
	router := gin.Default()
	router.Use(static.Serve("/", static.LocalFile("dist", false)))
	router.POST("/api/users", s.createUser)
	router.POST("/api/users/login", s.loginUser)
	authRoutes := router.Group("/api").Use(authMiddleWare(s.tokenMaker))

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
