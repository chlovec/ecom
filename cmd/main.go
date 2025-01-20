package main

import (
	"ecom/cmd/api"
	"ecom/config"
	"ecom/db"
	"log"
)

func main() {
	db, err := db.InitDB(config.GetDBConfig())
	if err != nil {
		log.Fatalf("Unable to open db connection: %v\n", err)
	}
	log.Print("Successfully connected to db!")

	err = api.InitServer(config.Envs.Port, db)
	if err != nil {
		log.Fatalf("Unable to start api server: %v\n", err)
	}
}
