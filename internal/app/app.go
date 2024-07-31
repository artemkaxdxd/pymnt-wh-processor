package app

import (
	"backend/config"

	// handlers
	handlerOrder "backend/internal/controllers/http/v1/order"

	// entities
	entityOrder "backend/internal/entity/order"

	// services
	serviceOrder "backend/internal/service/order"

	// storages
	repoOrder "backend/internal/storage/mysql/order"

	"backend/pkg/db"
	"backend/pkg/httpserver"
	"backend/pkg/logger"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
)

func Run(cfg config.Config) {
	l := logger.NewZap("")

	// Connect to storage
	mySQL, err := db.NewMySQL(cfg.MySQL, l)
	if err != nil {
		l.Fatal("Unable to make mySQL connection: ", err)
	}

	// Auto migrations
	if err = mySQL.DB.AutoMigrate(
		&entityOrder.Order{},
		&entityOrder.OrderEvent{},
	); err != nil {
		l.Fatal("Auto migration failed: ", err)
	}

	// Storages
	orderRepo := repoOrder.NewRepo(mySQL)

	// Services
	orderSvc := serviceOrder.NewService(orderRepo)

	// HTTP server
	gin.SetMode(gin.ReleaseMode)
	g := gin.New()

	g.Use(
		gin.Logger(),
		gin.Recovery(),
	)

	// Handlers
	g.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, "pong")
	})

	handlerOrder.InitHandler(
		g, l,
		orderSvc,
	)

	server := httpserver.New(g, httpserver.Port(cfg.Server.Port))

	// Waiting signals
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Fatal("Signal interrupt error: ", s.String())
	case err := <-server.Notify():
		l.Fatal("Server notify err", err)
	}

	// Shutdown server
	err = server.Shutdown()
	if err != nil {
		l.Info("Server shutdown err: ", err)
	}
}
