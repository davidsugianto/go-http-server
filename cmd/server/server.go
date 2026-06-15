package main

import (
	"net/http"

	httpHandler "github.com/davidsugianto/go-http-server/internal/handler/http"
	"github.com/davidsugianto/go-http-server/internal/handler/middleware"
	"github.com/davidsugianto/go-http-server/internal/pkg/config"
	"github.com/davidsugianto/go-pkgs/grace"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Server struct {
	*http.Server
	handler *httpHandler.Handler
	config  *config.Config
}

type Dependencies struct {
	Config *config.Config
}

func New(deps Dependencies) *Server {
	return &Server{
		Server:  &http.Server{},
		handler: httpHandler.New(httpHandler.Dependencies{}),
		config:  deps.Config,
	}
}

func (s *Server) routes(r *gin.Engine) {
	v1 := r.Group("/v1")
	v1.Use(gin.Recovery(), middleware.RequestID())

	v1.GET("/ping", s.handler.Ping)
}

func (s *Server) Run(port string) error {
	r := gin.Default()

	corsConfig := cors.Config{
		AllowOrigins:     s.config.CORS.AllowedOrigins,
		AllowMethods:     s.config.CORS.AllowedMethods,
		AllowHeaders:     s.config.CORS.AllowedHeaders,
		AllowCredentials: s.config.CORS.AllowCredentials,
	}
	r.Use(cors.New(corsConfig))

	s.routes(r)

	s.Addr = port
	s.Handler = r

	return grace.ServeHTTP(s.Addr, s.Handler)
}
