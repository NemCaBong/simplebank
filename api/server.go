package api

import (
	"fmt"

	"github.com/techschool/simplebank/token"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/techschool/simplebank/db/sqlc"
	"github.com/techschool/simplebank/db/util"
)

// this Server serves all HTTP request for banking services
type Server struct {
	config     util.Config
	tokenMaker token.Maker
	store      db.Store
	// router will sent each http request
	// to the handler for processing
	router *gin.Engine
}

// NewServer creates a new HTTP sever and setup routing
func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}
	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}
	server.setupRouter()

	// register the new validator with Gin
	// get the validator engine that Gin is using
	// Underlying the binding.Validator is a pointer to the StructValidator in "github.com/go-playground/validator/v10"
	if valid, ok := binding.Validator.Engine().(*validator.Validate); ok {
		valid.RegisterValidation("currency", validCurrency)
	}
	return server, nil
}

// run all the router method
func (server *Server) setupRouter() {
	router := gin.Default()

	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)

	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))
	authRoutes.POST("/accounts", server.createAccount)
	authRoutes.GET("/accounts/:id", server.getAccount)
	authRoutes.GET("/accounts", server.listAccounts)
	authRoutes.POST("/transfers", server.createTransferTx)

	server.router = router
}
func errorResponse(err error) gin.H {
	// gin.H is a map[string]interface{} <=> return every key-value as we like
	return gin.H{"error": err.Error()}
}

func (server *Server) Start(address string) error {
	// router field is private => cannot be access out api package
	return server.router.Run(address)
	// so that we have this public start function
	// because Start is viết hoa
}
