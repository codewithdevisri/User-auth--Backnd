package main

import (
	"github.com/Ayan25844/netflix/controller"
	"github.com/Ayan25844/netflix/cors"
	"log"
	"net/http"
)

func main() {
	r := controller.Router()
	err := http.ListenAndServe(":4000", cors.EnableCors(r))
	if err != nil {
		log.Fatal(err)
	}
}
