package main

import (
	"fmt"

	"github.com/Dsmit05/metida/internal/api"
	"github.com/Dsmit05/metida/internal/config"
	"github.com/Dsmit05/metida/internal/cryptography"
	"github.com/Dsmit05/metida/internal/logger"
	"github.com/Dsmit05/metida/internal/repositories"
)

// @title metida
// @version 0.0.1
// @description API template server

// @host localhost:8080
// @basePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorizations
func main() {
	// Init settings from cmd flag
	flagCmd, err := config.NewCommandLine()
	if err != nil {
		fmt.Println(err)
		return
	}

	// Init logger
	if err := logger.InitLogger(flagCmd); err != nil {
		panic(err)
	}

	// Init config api
	cfg, err := config.NewConfig(flagCmd)
	if err != nil {
		logger.L.Error(err)
		return
	}
	logger.L.Info(cfg)

	// Init connect to db
	db, err := repositories.NewPostgresRepository(cfg)
	if err != nil {
		logger.L.Error(err)
		return
	}
	defer db.Close()

	managerToken := cryptography.NewManagerToken("kita")

	apiV1 := api.V1(db, managerToken, cfg)
	if err := apiV1.Start(); err != nil {
		logger.L.Error(err)
	}

}
