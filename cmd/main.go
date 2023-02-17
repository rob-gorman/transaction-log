package main

import (
	"auditlog/cmd/api"
	"auditlog/internal/config"
	"auditlog/internal/data/cassandra"
	"auditlog/utils"
	"context"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	
	config.InitConfig()

	logger := utils.NewLogger()
	db := cassandra.New(ctx, logger)
	auth := db
	app := api.New(logger, db, auth)

	app.Run(ctx)
}
