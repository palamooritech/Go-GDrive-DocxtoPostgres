package main

import (
	drive "31arthur/drive-editor/pkg/adapter/gdrive"
	storage "31arthur/drive-editor/pkg/adapter/postgres"
	"31arthur/drive-editor/pkg/domain"
	"log"
	"time"
)

func main() {

	store, err := storage.NewPgxStore()
	if err != nil {
		log.Fatal("In main", err)
	}
	if err := store.Init(); err != nil {
		log.Fatal("In main", err)
	}

	ticker := time.NewTicker(1 * time.Hour)

	server := domain.NewAPIServer(":3001", store)
	go drive.DriveAdapter(server)
	go func() {
		for range ticker.C {
			drive.DriveAdapter(server)
		}
	}()

	server.Run()
	// fmt.Println("Hello World!")
}
