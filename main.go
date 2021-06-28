package main

import (
	"database/sql"
	"log"
)

func init() {
	for _, driverName := range sql.Drivers() {
		log.Printf("Load %s", driverName)
	}
}
