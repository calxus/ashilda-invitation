package main

import (
	"log"
	"net/http"
)

func main() {
	router := newRouter("/api/v1")
	router.populateRoutes()
	log.Fatal(http.ListenAndServe(":8082", router.mux))
}
