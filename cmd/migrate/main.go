package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"sass.com/configsvc/internal/migrations"
)

func main() {
	reset := flag.Bool("reset", false, "delete existing DB file before migration")
	seed := flag.Bool("seed", false, "insert default records after migration")
	flag.Parse()

	dbPath, _ := filepath.Abs("./data/config.db")
	if *reset {
		if err := os.Remove(dbPath); err == nil {
			log.Println("removed old DB file:", dbPath)
		} else {
			log.Println("failed to remove DB file:", err)
		}
	}

	if err := migrations.RunMigrations(dbPath, *seed); err != nil {
		log.Fatal(err)
	}
}
