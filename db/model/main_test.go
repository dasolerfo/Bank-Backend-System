package db

import (
	"database/sql"
	"log"
	"os"
	"simplebank/factory"
	"testing"

	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	var err error

	config, err := factory.LoadConfig("../..")
	if err != nil {
		log.Fatal("Error! No s'ha pogut carregar el .env", err)
	}
	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Error! No et pots connectar a la base de dades: ", err)
	}

	testQueries = New(testDB)
	os.Exit(m.Run())

}
