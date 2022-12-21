package main

import (
	"encoding/json"
	"log"
	"net/http"

	"organization-api/configs"
	"organization-api/routes"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "application/json")

		json.NewEncoder(rw).Encode(map[string]string{"data": "Welcome To our Organization"})

	}).Methods("GET")

	configs.ConnectDB()

	routes.OrganizationRoute(router)

	log.Fatal(http.ListenAndServe(":5050", router))

}
