package main

import (
	"set-game/src/service"
	"log"
)

func main() {
	srv, err := service.NewService(&service.Config{
		Port:      "8080",
		DBAddress: "./local/db.sqlite3",
	})
	if err != nil {
		log.Panic(err.Error())
	}
	if err = srv.Run(); err != nil {
		log.Panic(err.Error())
	}
}
