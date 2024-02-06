package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"pet-bank/api"
	db "pet-bank/db/sqlc"
	"pet-bank/utils"
)

func main() {
	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal("Couldn't load config: ", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Couldn't connect to db: ", err)
	}

	store := db.NewStore(conn)

	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("Couldn't create server: ", err)
	}

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("Couldn't start server: ", err)
	}
}
