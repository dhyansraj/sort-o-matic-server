package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"uol.com/sort-o-matic-server/service"
)

func main() {

	customerService := service.CustomerService{}

	router := mux.NewRouter()

	router.HandleFunc("/sortomatic/customers", customerService.GetAllCustomersHandler).Methods("GET")
	router.HandleFunc("/sortomatic/search/{name}", customerService.SearchCustomersHandler).Methods("GET")
	router.HandleFunc("/sortomatic/customers/{level}", customerService.GetCustomersByLevelHandler).Methods("GET")
	router.HandleFunc("/sortomatic/highest-customer", customerService.GetCustomersByHighestPointHandler).Methods("GET")
	router.HandleFunc("/sortomatic/average-point", customerService.GetAvgPointHandler).Methods("GET")

	log.Fatal(http.ListenAndServe(":8000", router))

}
