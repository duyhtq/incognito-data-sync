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
func (s *Server) Pdex24h(c *gin.Context) {
	transactions, err := s.transaction.Report24h()
	if err != nil {
		s.logger.Error("s.transaction.Report24h", zap.Error(err))
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

func (s *Server) PdexTradingV2(c *gin.Context) {
	rangeFilter := c.DefaultQuery("range", "day")
	token := c.DefaultQuery("token", "")

	transactions, err := s.transaction.PdexTradingV2(rangeFilter, token)
	if err != nil {
		s.logger.Error("s.transaction.PdexTradingV2", zap.Error(err))
		c.JSON(http.StatusInternalServerError, serializers.Resp{Error: err})
		return
	}

	c.JSON(http.StatusOK, serializers.Resp{
		Result: map[string]interface{}{
			"Data": transactions,
		},
	})
}
func (s *Server) Pdex24hV2(c *gin.Context) {
	transactions, err := s.transaction.Report24hV2()
	if err != nil {
		s.logger.Error("s.transaction.Pdex24hV2", zap.Error(err))
		c.JSON(http.StatusInternalServerError, serializers.Resp{Error: err})
		return
	}

	c.JSON(http.StatusOK, serializers.Resp{
		Result: map[string]interface{}{
			"Data": transactions,
		},
	})
}

func (s *Server) Shield(c *gin.Context) {
	transactions, err := s.transaction.Shield()
	if err != nil {
		s.logger.Error("s.transaction.Shield", zap.Error(err))
		c.JSON(http.StatusInternalServerError, serializers.Resp{Error: err})
		return
	}

	c.JSON(http.StatusOK, serializers.Resp{
		Result: map[string]interface{}{
			"Data": transactions,
		},
	})
}
func (s *Server) Unshield(c *gin.Context) {
	transactions, err := s.transaction.Unshield()
	if err != nil {
		s.logger.Error("s.transaction.Unshield", zap.Error(err))
		c.JSON(http.StatusInternalServerError, serializers.Resp{Error: err})
		return
	}

	c.JSON(http.StatusOK, serializers.Resp{
		Result: map[string]interface{}{
			"Data": transactions,
		},
	})
}

func (s *Server) Shield24h(c *gin.Context) {
	transactions, err := s.transaction.Shield24h()
	if err != nil {
		s.logger.Error("s.transaction.Shield24h", zap.Error(err))
		c.JSON(http.StatusInternalServerError, serializers.Resp{Error: err})
		return
	}

	c.JSON(http.StatusOK, serializers.Resp{
		Result: map[string]interface{}{
			"Data": transactions,
		},
	})
}
func (s *Server) Unshield24h(c *gin.Context) {
	transactions, err := s.transaction.Unshield24h()
	if err != nil {
		s.logger.Error("s.transaction.Unshield24h", zap.Error(err))
		c.JSON(http.StatusInternalServerError, serializers.Resp{Error: err})
		return
	}

	c.JSON(http.StatusOK, serializers.Resp{
		Result: map[string]interface{}{
			"Data": transactions,
		},
	})
}

func (s *Server) ShieldMonth(c *gin.Context) {
	transactions, err := s.transaction.ShieldMonth()
	if err != nil {
		s.logger.Error("s.transaction.ShieldMonth", zap.Error(err))
		c.JSON(http.StatusInternalServerError, serializers.Resp{Error: err})
		return
	}

	c.JSON(http.StatusOK, serializers.Resp{
		Result: map[string]interface{}{
			"Data": transactions,
		},
	})
}
func (s *Server) UnshieldMonth(c *gin.Context) {
	transactions, err := s.transaction.UnshieldMonth()
	if err != nil {
		s.logger.Error("s.transaction.UnshieldMonth", zap.Error(err))
		c.JSON(http.StatusInternalServerError, serializers.Resp{Error: err})
		return
	}

	c.JSON(http.StatusOK, serializers.Resp{
		Result: map[string]interface{}{
			"Data": transactions,
		},
	})
}
