package db

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"os"
	"pet-bank/utils"
	"testing"
)

var testStore Store

func TestMain(m *testing.M) {
	config, err := utils.LoadConfig("../..")
	if err != nil {
		log.Fatal("Couldn't load config: ", err)
	}

	connPool, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal("Couldn't connect to db: ", err)
	}

	testStore = NewStore(connPool)

	os.Exit(m.Run())
}
