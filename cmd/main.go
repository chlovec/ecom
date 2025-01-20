package main

import (
	"database/sql"
	"ecom/cmd/api"
	"ecom/config"
	"ecom/db"
	"log"
)

func main() {
	db, err := db.NewMySQLStorage(config.GetDBConfig())
	if err != nil {
		log.Fatalf("Unable to open db connection: %v\n", err)
	}

	initStorage(db)

	server := api.NewAPIServer(config.Envs.Port, db)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}

func initStorage(db *sql.DB) {
	err := db.Ping()
	if err != nil {
		log.Fatalf("Unable to connect to db: %v\n", err)
	}

	log.Print("Successfully connected to db!")
}
