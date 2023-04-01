package main

import (
	"fmt"
	"log"

	"example.com/api"
	"example.com/storage"
	"example.com/types"
)

func main() {
	var app types.App
	var store *storage.DatabaseStore
	var err error

	app.Init()
	store, err = storage.NewDatabaseStore(app.IsSqlite)
	handleError(err)

	err = seedDatabase(store, app.IsSeeding)
	handleError(err)

	err = store.Init()
	handleError(err)

	server := api.NewApiServer(app.Port, store)
	server.Run()
}

func seedDatabase(store *storage.DatabaseStore, isSeeding bool) error {
	if !isSeeding {
		return nil
	}

	fmt.Println("Dropping current data and seeding with new user")

	account, err := store.SeedDatabase()
	if err != nil {
		return err
	}
	fmt.Printf("\nAccount %d created\nEmail: %s\n\n", account.AccountNumber, account.Email)

	return nil
}

func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
