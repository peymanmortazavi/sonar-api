package main

import (
	"flag"
	"log"
	"sonar/pkg/csvstore"
)

func main() {
	csvStore := flag.String("csv-store", "", "path to the CSV file.")

	flag.Parse()

	_, err := csvstore.CreateMemoryStoreFromFile(*csvStore)
	if err != nil {
		log.Fatalf("failed to parse sonar data from CSV file: %s", err)
	}
}
