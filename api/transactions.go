package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/duyhtq/incognito-data-sync/serializers"
)

func (s *Server) ReportPdexTrading(c *gin.Context) {

	transactions, err := s.transaction.ReportPdexTrading()
	if err != nil {
		s.logger.Error("s.transaction.ReportPdexTrading", zap.Error(err))
		c.JSON(http.StatusInternalServerError, serializers.Resp{Error: err})
		return
	}

	c.JSON(http.StatusOK, serializers.Resp{
		Result: map[string]interface{}{
			"Data": transactions,
		},
	})
}
