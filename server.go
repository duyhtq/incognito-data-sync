package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/duyhtq/incognito-data-sync/api"
	config "github.com/duyhtq/incognito-data-sync/config"
	"github.com/duyhtq/incognito-data-sync/databases/postgresql"
	"github.com/duyhtq/incognito-data-sync/service"
)

func main() {
	conf := config.GetConfig()

	logger := service.NewLogger(conf)

	db, err := postgresql.Init(conf)

	if err != nil {
		log.Println("error:", err)
	}

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "POST", "PUT", "HEAD", "OPTIONS", "DELETE"},
		AllowHeaders:    []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		MaxAge:          12 * time.Hour,
	}))

	databaseConnect, err := postgresql.NewTransactionsStore(db)
	if err != nil {
		log.Println("error:", err)
	}

	var transaction = service.NewTransactionService(conf, databaseConnect, logger.With(zap.String("service", "transaction")))

	svr := api.NewServer(r, transaction, logger.With(zap.String("module", "api")))

	svr.Routes()

	if err := r.Run(fmt.Sprintf(":%d", conf.Port)); err != nil {
		logger.Fatal("router.Run", zap.Error(err))
	}
}
