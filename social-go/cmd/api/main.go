package main

import (
	"log"
	"os"
	"social/internal/db"
	"social/internal/env"
	store2 "social/internal/store"

	"github.com/joho/godotenv"
)

const version = "0.0.1"

// seguir en la seccion 7 del curso de udemy en el capitulo 43
func main() {
	if os.Getenv("APP_ENV") != "production" {
		if err := godotenv.Load(); err != nil {
			log.Println("warning: .env not loaded:", err)
		}
	}

	cfg := Config{
		addr: env.GetString("SERVER_ADDR", ":8080"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost:5432/social?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 10),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 5),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "12m"),
		},
		env:     env.GetString("ENV", "development"),
		version: version,
	}

	dbConnection, err := db.New(cfg.db.addr, cfg.db.maxOpenConns, cfg.db.maxIdleConns, cfg.db.maxIdleTime)

	if err != nil {
		log.Panic(err)
	}

	defer dbConnection.Close()
	log.Println("DB connection established")

	store := store2.NewStorage(dbConnection)

	app := &Application{
		config: cfg,
		store:  store,
	}

	mux := app.mount()

	log.Fatal(app.run(mux))
}
