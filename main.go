package main

import (
	"database/sql"
	"log"
	"simplebank/api"
	db "simplebank/db/model"
	"simplebank/factory"

	_ "github.com/lib/pq"
)

func main() {
	config, err := factory.LoadConfig(".")
	if err != nil {
		log.Fatal("Error! No s'ha pogut carregar el .env", err)
	}
	testDB, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Error! No et pots connectar a la base de dades: ", err)
	}

	store := db.NewStore(testDB)
	router := api.NewServer(store)

	err = router.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("Error! No es pot inicialitzar el server: ", err)
	}
}
