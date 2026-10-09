package main

import (
	"context"
	"errors"
	"link_watcher/checker"
	"link_watcher/config"
	"link_watcher/core/link_watcher"
	"link_watcher/db"
	"link_watcher/logger"
	"link_watcher/pgService"
	transport "link_watcher/transport/http"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {

	cfg, err := config.Load("config.yaml")
	if err != nil {
		panic(err)
	}

	logger.Init(logger.Config{
		Level:  cfg.Logger.Level,
		Format: cfg.Logger.Format,
	})
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

	ctx, cancel := context.WithCancel(context.Background())
	//defer cancel()

	timeout := time.Duration(cfg.Checker.TimeoutSec) * time.Second
	chk := checker.NewChecker(pg, timeout)

	tasks := make(chan checker.Task, 1000)
	sc := checker.NewSchedule(pg, tasks, 1*time.Second)
	r := checker.NewRun(chk, tasks, 50)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		sc.Start(ctx)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		r.RunChecker(ctx)
	}()

	logger.Info("Background process started", "goroutines", 50, "timeout_sec", cfg.Checker.TimeoutSec)

	srv := &http.Server{
		Addr:    cfg.Server.Addr,
		Handler: handler,
	}

	go func() {
		logger.Info("HTTP server listening on " + cfg.Server.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal("Fail to start http server", "error", err.Error())
		}
	}()

	sigCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	<-sigCtx.Done()
	logger.Info("Shutdown signal received")

	cancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("Fail to shutdown http server", "error", err.Error())
	} else {
		logger.Info("Server shutdown successfully")
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		logger.Info("Background process finished")
	case <-shutdownCtx.Done():
		logger.Error("Background process timed out")
	}

	logger.Info("Shutdown complete")

}
