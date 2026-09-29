package main

import (
	"flag"
	"link_watcher/config"
	"link_watcher/db"
	"link_watcher/db/migrations"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

func main() {
	configPath := flag.String("config", "config.yaml", "путь к config.yaml")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	dbConn, err := db.Connect(cfg.DB)
	if err != nil {
		log.Fatalf("db connect: %v", err)
	}
	defer func() {
		if err := dbConn.Close(); err != nil {
			log.Printf("close db: %v", err)
		}
	}()

	driver, err := postgres.WithInstance(dbConn, &postgres.Config{
		MigrationsTable: "schema_migrations",
		DatabaseName:    cfg.DB.DBName,
	})
	if err != nil {
		log.Fatalf("migrate driver: %v", err)
	}

	source, err := iofs.New(migrations.MigrateFS, ".")
	if err != nil {
		log.Fatalf("migrate iofs: %v", err)
	}

	m, err := migrate.NewWithInstance("file", source, "postgres", driver)
	if err != nil {
		log.Fatalf("migrate: %v", err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Fatalf("migrate: %v", err)
	}

	log.Println("migrations applied successfully")
}
