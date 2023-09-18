package api

import (
	"fmt"
	db "simplebank/db/sqlc"
	"simplebank/token"
	"simplebank/util"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	config       util.Config
	store        db.Store
	tokenManager token.TokenManager
	router       *gin.Engine
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	token, err := token.NewPasetoToken(config.TokenSymetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token manager: %w", err)
	}

	server := &Server{
		config:       config,
		store:        store,
		tokenManager: token,
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	server.setupRouter()

	return server, nil
}

func (server *Server) setupRouter() {

	router := gin.Default()

	router.POST("/user", server.createUser)
	router.POST("/user/login", server.loginUser)

	router.POST("/token/renew_access", server.renewAccessToken)

	authRoutes := router.Group("/").Use(authMiddleware(server.tokenManager))

	authRoutes.POST("/account", server.createAccount)
	authRoutes.GET("/account/:id", server.getAccount)
	authRoutes.GET("/account", server.listAccount)

	authRoutes.POST("/transfer", server.transfer)

	server.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
