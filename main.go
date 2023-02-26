package main

import (
	"database/sql"
	"log"

	_ "github.com/golang/mock/mockgen/model"
	_ "github.com/lib/pq"
	"github.com/prepStation/simple_bank/api"
	db "github.com/prepStation/simple_bank/db/sqlc"
	"github.com/prepStation/simple_bank/utils"
)

func main() {
	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal("Cannot read config: ", err)
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatalf("Cannot connect to database %v\n", err)
	}
	store := db.NewStore(conn)
	server := api.NewServer(store)
	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatalf("Cannot start server : %v\n", err)
	}
}
