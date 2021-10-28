package main

import (
	"flag"
	"log"
	"net/http"
	"sonar/pkg/csvstore"
	"sonar/pkg/restapi"
)

func main() {
	csvStore := flag.String("csv-store", "", "path to the CSV file.")

	flag.Parse()

	storage, err := csvstore.CreateMemoryStoreFromFile(*csvStore)
	if err != nil {
		log.Fatalf("failed to parse sonar data from CSV file: %s", err)
	}

	handler := restapi.Handler{ReadScanner: storage}
	http.Handle("/records", handler)

	log.Fatalln(http.ListenAndServe("0.0.0.0:6000", nil))
}
