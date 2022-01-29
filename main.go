package main

import (
	"TaskOne/main/controller"
	_ "TaskOne/main/docs"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger"
	"log"
	"net/http"
)

func init() {

	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}

}

func setupRoutes(){

	router := mux.NewRouter()
	// create group of routes "/api/v1"
	routeGroup := router.PathPrefix("/api/v1").Subrouter()

	//create routes
	//crud on order
	//post request to create a order
	routeGroup.HandleFunc("/orders", controller.Create).Methods("POST")

	//get request to fetch all orders
	routeGroup.HandleFunc("/orders", controller.GetAll).Methods("GET")

	//delete req to delete a one order with the id
	routeGroup.HandleFunc("/orders/id", controller.Delete).Methods("DELETE")

	//get req to get one order with the order id
	routeGroup.HandleFunc("/orders/id", controller.GetOne).Methods("GET")

	//put req to update the data
	routeGroup.HandleFunc("/orders/id", controller.Update).Methods("PUT")

	//upload file
	routeGroup.HandleFunc("/orders/upload/id", controller.UploadFile).Methods("POST")

	//download file
	routeGroup.HandleFunc("/orders/download", controller.DownloadFile).Methods("POST")

	// list all files
	routeGroup.HandleFunc("/orders/files", controller.GetAllFiles).Methods("GET")

	routeGroup.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)

	routeGroup.Use(mux.CORSMethodMiddleware(routeGroup))
	log.Fatal(http.ListenAndServe(":8080", routeGroup))

}

// @title Orders API
// @version 1.0
// @description This is a sample service for managing orders
// @termsOfService http://swagger.io/terms/
// @contact.name Nabil Salman
// @contact.email nabil.salman@symbyo.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath /api/v1
func main() {
	setupRoutes()
	//dynamoDb.CreateTable("ItemLists")

}
