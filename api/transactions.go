package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/duyhtq/incognito-data-sync/serializers"
)

func (s *Server) ReportPdexTrading(c *gin.Context) {
	rangeFilter := c.DefaultQuery("range", "day")
	token := c.DefaultQuery("token", "")

	transactions, err := s.transaction.ReportPdexTrading(rangeFilter, token)
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

func (s *Server) PdexVolume(c *gin.Context) {
	token1str := c.DefaultQuery("token1str", "")
	token2str := c.DefaultQuery("token2str", "")

	volume, err := s.transaction.PdexVolume(token1str, token2str)
	fmt.Println(volume)
	if err != nil {
		s.logger.Error("s.transaction.PdexVolume", zap.Error(err))
		c.JSON(http.StatusInternalServerError, serializers.Resp{Error: err})
		return
	}

	c.JSON(http.StatusOK, serializers.Resp{
		Result: volume,
	})
}
