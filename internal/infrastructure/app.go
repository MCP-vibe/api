package infrastructure

import (
	"api/internal/adapters/logger"
	"api/internal/adapters/validator"
	"api/internal/config"
	db "api/internal/infrastructure/database"
	"api/internal/infrastructure/log"
	"api/internal/infrastructure/router"
	"api/internal/infrastructure/validation"
	"time"
)

type app struct {
	cfg        config.Config
	logger     logger.Logger
	validator  validator.Validator
	dbManager  db.DBManager
	ctxTimeout time.Duration
	webServer  *router.GinEngine
}

func NewConfig(config config.Config) *app {
	return &app{
		cfg: config,
	}
}

func (a *app) ContextTimeout(t time.Duration) *app {
	a.ctxTimeout = t
	return a
}

func (a *app) Logger() *app {
	log := log.NewLogrusLogger()
	a.logger = log
	a.logger.Infof("Success log configured")
	return a
}

func (a *app) Database() *app {
	gormDB, err := db.NewPostgresConnection(a.cfg)
	if err != nil {
		a.logger.Fatalln("Failed to configure database:", err)
	}
	a.dbManager = *db.NewDBManager(gormDB)
	a.logger.Infof("Success database configured")
	return a
}

func (a *app) Validator() *app {
	v, err := validation.NewGoPlayground()
	if err != nil {
		a.logger.Fatalln(err)
	}

	a.logger.Infof("Success validator configured")

	a.validator = v

	return a
}

func (a *app) WebServer() *app {
	s := router.NewGinServer(
		a.cfg,
		a.logger,
		a.validator,
		a.dbManager,
		a.ctxTimeout,
	)

	a.logger.Infof("Success router server configured")

	a.webServer = s
	return a
}

func (a *app) Start() {
	a.webServer.Listen()
}
