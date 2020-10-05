package api

import "github.com/gin-gonic/gin"

func (s *Server) Routes() {

	s.g.GET("/pdex-trading", s.ReportPdexTrading)

	s.g.GET("/pdex-24h", s.Pdex24h)

	s.g.GET("/v2/pdex-trading", s.PdexTradingV2)
	s.g.GET("/v2/pdex-24h", s.Pdex24hV2)

	s.g.GET("/shield", s.Shield)
	s.g.GET("/unshield", s.Unshield)

	s.g.GET("/pdex-volume", s.PdexVolume)

	s.g.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
}
