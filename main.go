package main

import (
	"database/sql"
	"log"
	"simplebank/api"
	db "simplebank/db/model"

	_ "github.com/lib/pq"
)

const (
	dbDriver      = "postgres"
	dbSource      = "postgresql://root:Songoku13@localhost:5432/simple_bank?sslmode=disable"
	serverAddress = "0.0.0.0:8080"
)

func main() {
	testDB, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("Error! No et pots connectar a la base de dades: ", err)
	}

	store := db.NewStore(testDB)
	router := api.NewServer(store)

	err = router.Start(serverAddress)
	if err != nil {
		log.Fatal("Error! No es pot inicialitzar el server: ", err)
	}
}
