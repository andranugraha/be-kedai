package main

import (
	"kedai/backend/be-kedai/connection"
	"log"
)

func main() {
	err := connection.ConnectDB()
	if err != nil {
		log.Fatal("couldn't connect to DB:", err.Error())
	}

	err = connection.ConnectCache()
	if err != nil {
		log.Fatal("couldn't connect to Cache:", err.Error())
	}
}
