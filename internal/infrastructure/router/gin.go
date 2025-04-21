package router

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "api/docs" // Пакет, содержащий сгенерированный код Swagger

	"api/internal/adapters/api/middleware"
	"api/internal/adapters/logger"
	"api/internal/adapters/validator"
	"api/internal/config"
	db "api/internal/infrastructure/database"

	"github.com/gin-gonic/gin"
	"github.com/rs/cors"
	files "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/urfave/negroni"
)

type GinEngine struct {
	cfg        config.Config
	router     *gin.Engine
	log        logger.Logger
	validator  validator.Validator
	db         db.DBManager
	ctxTimeout time.Duration
}

func NewGinServer(
	cfg config.Config,
	log logger.Logger,
	validator validator.Validator,
	db db.DBManager,
	ctxTimeout time.Duration,
) *GinEngine {
	return &GinEngine{
		cfg:        cfg,
		router:     gin.New(),
		log:        log,
		validator:  validator,
		db:         db,
		ctxTimeout: ctxTimeout,
	}
}

func (g GinEngine) Listen() {
	gin.SetMode(gin.ReleaseMode)
	gin.Recovery()

	g.setupRoutes(g.router)

	n := negroni.New()
	c := cors.New(cors.Options{
		AllowedOrigins:      []string{"*"}, // TODO: Add allowed origins
		AllowPrivateNetwork: true,
		AllowedMethods:      []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodOptions, http.MethodDelete, http.MethodPatch},
		AllowCredentials:    true,
	})

	n.Use(c)
	n.Use(negroni.HandlerFunc(middleware.NewLogger(g.log).Execute))
	n.UseHandler(g.router)
	server := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 15 * time.Second,
		Addr:         fmt.Sprintf(":%d", g.cfg.AppPort),
		Handler:      n,
	}

	// Ctrl-C
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		g.log.WithFields(logger.Fields{"port": g.cfg.AppPort}).Infof("Starting HTTP server...")
		if err := server.ListenAndServe(); err != nil {
			g.log.WithError(err).Fatalln("Error starting HTTP server")
		}
	}()

	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	if err := server.Shutdown(ctx); err != nil {
		g.log.WithError(err).Fatalln("Server Shutdown Failed")
	}

	g.log.Infof("Service down")
}

func (g GinEngine) setupRoutes(router *gin.Engine) {

	router.GET("/swagger/*any", ginSwagger.WrapHandler(files.Handler, ginSwagger.URL("/docs/swagger.yaml")))
}

func buildParams(c *gin.Context, params ...string) {
	q := c.Request.URL.Query()

	for _, value := range params {
		switch value {
		case "page":
			if _, exists := q["page"]; !exists {
				q.Set("page", "1")
			}
		case "limit":
			if _, exists := q["limit"]; !exists {
				q.Set("limit", "10")
			}
		default:
			q.Add(value, c.Param(value))
		}
	}
	c.Request.URL.RawQuery = q.Encode()
}

type ErrorResponse struct { // временно errors response (мб поменять)
	Error string `json:"error" example:"error message"`
}
