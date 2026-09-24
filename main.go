package main

import (
	"net/http"
)

type Link struct {
	DB *sql.DB
}

func NewLink(db *sql.DB) *Link {
	return &Link{
		DB: db,
	}

}

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

	mux := http.NewServeMux()

	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World"))
	})

	if err := http.ListenAndServe(":8090", mux); err != nil {
		logger.Fatal("Server failed to start", "error", err.Error())
	}

}
