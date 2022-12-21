package routes

import (
	"organization-api/controllers"

	"github.com/gorilla/mux"
)

func OrganizationRoute(router *mux.Router) {
	router.HandleFunc("/organization", controllers.CreateOrganization()).Methods("POST")
	router.HandleFunc("/organization/{organizationId}", controllers.GetAOrganization()).Methods("GET")
	router.HandleFunc("/organization/{organizationId}", controllers.EditAOrganization()).Methods("PUT")
	router.HandleFunc("/organization/{organizationId}", controllers.DeleteAOrganization()).Methods("DELETE")
	router.HandleFunc("/organizations", controllers.GetAllOrganization()).Methods("GET")

}
