package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {
	address := flag.String("address", "localhost:3000", "Address of the server")
	folder := flag.String("folder", ".", "Folder that should be served")
	flag.Parse()

	fs := http.FileServer(http.Dir(*folder))
	http.Handle("/", fs)

	log.Print("Listening on " + *address + " to serve " + *folder + "...")
	err := http.ListenAndServe("localhost:3000", nil)
	if err != nil {
		log.Fatal(err)
	}
}
