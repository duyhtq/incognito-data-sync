package api

import "github.com/gin-gonic/gin"

func (s *Server) Routes() {

	s.g.GET("/transactions", s.ListTransactionByPublicKey)

	s.g.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
}
