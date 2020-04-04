package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/duyhtq/incognito-data-sync/serializers"
)

func (s *Server) ListTransactionByPublicKey(c *gin.Context) {
	fmt.Println("ListTransactionByPublicKey")
	// page, limit := s.pagingFromContext(c)
	publicKeys := c.DefaultQuery("publicKeys", "")
	id := c.DefaultQuery("id", "0")

	// filter := map[string]interface{}{}
	// if memo != "" {
	// 	filter["memo"] = memo
	// }
	idInt, err := strconv.Atoi(id)
	if err != nil {
		s.logger.Error("s.transaction", zap.Error(err))
		c.JSON(http.StatusInternalServerError, serializers.Resp{Error: err})
		return
	}

	transactions, err := s.transaction.ListTransaction(publicKeys, idInt)
	if err != nil {
		s.logger.Error("s.transaction", zap.Error(err))
		c.JSON(http.StatusInternalServerError, serializers.Resp{Error: err})
		return
	}

	c.JSON(http.StatusOK, serializers.Resp{
		Result: map[string]interface{}{
			"Transactions": transactions,
		},
	})
}
