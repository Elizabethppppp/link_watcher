package main

import (
	"link_watcher/config"
	"link_watcher/core/link_watcher"
	"link_watcher/db"
	"link_watcher/logger"
	"link_watcher/pgService"
	transport "link_watcher/transport/http"
	"net/http"
)

func main() {

	cfg, err := config.Load("config.yaml")
	if err != nil {
		panic(err)
	}

	if err := logger.Init(logger.Config{
		Level:  cfg.Logger.Level,
		Format: cfg.Logger.Format,
	}); err != nil {
		panic(err)
	}
	logger.Info("Config successfully parsed", "dbHost", cfg.DB.Host, "dbName", cfg.DB.DBName, "dbSchema", cfg.DB.Schema)

	logger.Debug("Database conection", "host", cfg.DB.Host, "port", cfg.DB.Port, "dbName", cfg.DB.DBName)

	dbConn, err := db.Connect(cfg.DB)
	if err != nil {
		logger.Fatal("Fail connect database", "error", err.Error())
	}

	defer func() {
		if err := dbConn.Close(); err != nil {
			logger.Error("Failed to close db", "error", err.Error())
		}
	}()

	logger.Info("Connection successfully established", "host", cfg.DB.Host, "port", cfg.DB.Port)

	pg := pgService.NewPgService(dbConn)
	svc := link_watcher.NewReduceService(pg)
	tp := transport.NewTransport(svc)

	handler := tp.Handler()

	if err := http.ListenAndServe(":8090", handler); err != nil {
		logger.Fatal("Server failed to start", "error", err.Error())
	}

}
