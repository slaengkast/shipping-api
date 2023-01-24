package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type bookingHandler interface {
	BookShipping(c *gin.Context)
	GetBooking(c *gin.Context)
}

type server struct {
	router         *gin.Engine
	port           int
	bookingHandler bookingHandler
}

func New(bookingHandler bookingHandler, port int) *server {
	router := gin.New()
	return &server{
		router:         router,
		port:           port,
		bookingHandler: bookingHandler,
	}
}

func (s *server) Run() error {
	log.Info().Msgf("starting http server on port %d", s.port)

	s.router.Use(logMiddleware())
	s.router.Use(gin.Recovery())
	s.setupRoutes()

	return s.router.Run(fmt.Sprintf(":%d", s.port))
}

func (s *server) handleHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "healthy"})
}

func (s *server) setupRoutes() {
	s.router.GET("/health", s.handleHealth)

	apiRouter := s.router.Group("api")
	{
		apiRouter.GET("/shipping/:id", s.bookingHandler.GetBooking)
		apiRouter.POST("/shipping", s.bookingHandler.BookShipping)
	}
}

func logMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		status := c.Writer.Status()
		method := c.Request.Method
		path := c.Request.URL.Path
		switch status {
		case http.StatusInternalServerError:
			log.Error().Int("status", status).Str("method", method).Str("path", path).Msg("")
		case http.StatusBadRequest:
			log.Info().Int("status", status).Str("method", method).Str("path", path).Msg("")
		case http.StatusNotFound:
			log.Info().Int("status", status).Str("method", method).Str("path", path).Msg("")
		case http.StatusConflict:
			log.Info().Int("status", status).Str("method", method).Str("path", path).Msg("")
		}
	}
}
