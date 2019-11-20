package main

import (
	"log"
	"os"

	"github.com/lffranca/trophy-count-model/database"
	mygrpc "github.com/lffranca/trophy-count-model/grpc"
)

func main() {
	db, errDB := database.InitDatabase()
	if errDB != nil {
		log.Println(errDB)
		os.Exit(1)
	}

	if err := mygrpc.InitGRPC(db); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
