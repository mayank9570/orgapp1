package controllers

import (
	"context"
	"encoding/json"
	"organization-api/configs"
	"organization-api/models"
	"organization-api/responses"

	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var organizationCollection *mongo.Collection = configs.GetCollection(configs.DB, "organizations")
var validate = validator.New()

func CreateOrganization() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var organization models.Organization
		defer cancel()

		if err := json.NewDecoder(r.Body).Decode(&organization); err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			response := responses.OrganizationResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		if validationErr := validate.Struct(&organization); validationErr != nil {
			rw.WriteHeader(http.StatusBadRequest)
			response := responses.OrganizationResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		newOrganization := models.Organization{
			Id:       primitive.NewObjectID(),
			Name:     organization.Name,
			Location: organization.Location,
			City:     organization.City,
			State:    organization.State,
			Email:    organization.Email,
		}
		result, err := organizationCollection.InsertOne(ctx, newOrganization)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			response := responses.OrganizationResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		rw.WriteHeader(http.StatusCreated)
		response := responses.OrganizationResponse{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"data": result}}
		json.NewEncoder(rw).Encode(response)
	}
}
func GetAOrganization() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		params := mux.Vars(r)
		organizationId := params["organizationId"]
		var organization models.Organization
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(organizationId)

		err := organizationCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&organization)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			response := responses.OrganizationResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		rw.WriteHeader(http.StatusOK)
		response := responses.OrganizationResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": organization}}
		json.NewEncoder(rw).Encode(response)
	}
}

func EditAOrganization() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		params := mux.Vars(r)
		organizationId := params["organizationId"]
		var organization models.Organization
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(organizationId)

		if err := json.NewDecoder(r.Body).Decode(&organization); err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			response := responses.OrganizationResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		if validationErr := validate.Struct(&organization); validationErr != nil {
			rw.WriteHeader(http.StatusBadRequest)
			response := responses.OrganizationResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		update := bson.M{"name": organization.Name, "location": organization.Location, "city": organization.City, "state": organization.State}
		result, err := organizationCollection.UpdateOne(ctx, bson.M{"id": objId}, bson.M{"$set": update})
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			response := responses.OrganizationResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		var updatedOrganization models.Organization
		if result.MatchedCount == 1 {
			err := organizationCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&updatedOrganization)
			if err != nil {
				rw.WriteHeader(http.StatusInternalServerError)
				response := responses.OrganizationResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
				json.NewEncoder(rw).Encode(response)
				return
			}
		}

		rw.WriteHeader(http.StatusOK)
		response := responses.OrganizationResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": updatedOrganization}}
		json.NewEncoder(rw).Encode(response)
	}
}

func DeleteAOrganization() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		params := mux.Vars(r)
		organizationId := params["organizationId"]
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(organizationId)

		result, err := organizationCollection.DeleteOne(ctx, bson.M{"id": objId})

		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			response := responses.OrganizationResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		if result.DeletedCount < 1 {
			rw.WriteHeader(http.StatusNotFound)
			response := responses.OrganizationResponse{Status: http.StatusNotFound, Message: "error", Data: map[string]interface{}{"data": "Organization with specified ID not found!"}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		rw.WriteHeader(http.StatusOK)
		response := responses.OrganizationResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": "Organization successfully deleted!"}}
		json.NewEncoder(rw).Encode(response)
	}
}

func GetAllOrganization() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var organizations []models.Organization
		defer cancel()

		results, err := organizationCollection.Find(ctx, bson.M{})

		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			response := responses.OrganizationResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		defer results.Close(ctx)
		for results.Next(ctx) {
			var singleOrganization models.Organization
			if err = results.Decode(&singleOrganization); err != nil {
				rw.WriteHeader(http.StatusInternalServerError)
				response := responses.OrganizationResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
				json.NewEncoder(rw).Encode(response)
			}
			organizations = append(organizations, singleOrganization)
		}

		rw.WriteHeader(http.StatusOK)
		response := responses.OrganizationResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": organizations}}
		json.NewEncoder(rw).Encode(response)
	}
}
