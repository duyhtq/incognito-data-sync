package api

import (
	"strconv"

	"github.com/duyhtq/incognito-data-sync/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Server struct {
	g           *gin.Engine
	transaction *service.Transaction
	logger      *zap.Logger
}

func (s *Server) pagingFromContext(c *gin.Context) (int, int) {
	var (
		pageS  = c.DefaultQuery("page", "1")
		limitS = c.DefaultQuery("limit", "10")
		page   int
		limit  int
		err    error
	)

	page, err = strconv.Atoi(pageS)
	if err != nil {
		page = 1
	}

	limit, err = strconv.Atoi(limitS)
	if err != nil {
		limit = 10
	}

	return page, limit
}

func NewServer(g *gin.Engine, transaction *service.Transaction, logger *zap.Logger) *Server {
	return &Server{
		g:           g,
		transaction: transaction,
		logger:      logger,
	}
}
