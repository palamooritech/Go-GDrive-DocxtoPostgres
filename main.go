package main

import (
	drive "31arthur/drive-editor/pkg/adapter/gdrive"
	storage "31arthur/drive-editor/pkg/adapter/postgres"
	"31arthur/drive-editor/pkg/domain"
	"log"
	"time"
)

func main() {

	// creating an instance of NewPgxStore (for Postgres Connection),
	// which can be accessed as immutable thoughout the application
	store, err := storage.NewPgxStore()
	if err != nil {
		log.Fatal("In main", err)
	}
	if err := store.Init(); err != nil {
		log.Fatal("In main", err)
	}

	ticker := time.NewTicker(1 * time.Hour)

	//creating an instance of NewAPIServer to store the listening address
	// and the PGXStore instance for accessing them immutablely throughout
	// the application.
	server := domain.NewAPIServer(":3001", store)
	go drive.DriveAdapter(server)

	// for making the drive adapter run every one hour
	go func() {
		for range ticker.C {
			drive.DriveAdapter(server)
		}
	}()

	//server is of instance APIServer, which is used to run/start the web server
	server.Run()
	// fmt.Println("Hello World!")
}
