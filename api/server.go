package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/techschool/simplebank/db/sqlc"
)

// this Server serves all HTTP request for banking services
type Server struct {
	store db.Store
	// router will sent each http request
	// to the handler for processing
	router *gin.Engine
}

// NewServer creates a new HTTP sever and setup routing
func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccounts)
	// the create account need to be a method of server struct
	// because it needs to get access to the store obj
	// to save new account in db
	server.router = router
	return server
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
