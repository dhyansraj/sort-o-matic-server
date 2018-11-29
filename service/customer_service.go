package service

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"net/http"
	"strings"
)

var customerTable [][]string
var pointsTable [][]int
var rankTable [][]int
var customers []CustomerModel

var ranks [][]string


type CustomerModel struct {
	CustomerId int
	FirstName string
	LastName string
	Points int
}

type Customer struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Rank      int    `json:"rank"`
	Level     int    `json:"level"`
	Points    int    `json:"points"`
}

func (CustomerModel) TableName() string {
	return "customers"
}

type CustomerService struct {

}

//Load Customers from Database.
//Build Rank table
//Build Hash table
func init(){

	customers = loadAllCustomers()
	BuildRank()
	BuildTable()

}

//Allocate customers with same points in to an array.
//Last cell in the map with non black array will be the highest point
func BuildRank() {

	heighestPoint := 0

	for _, customer := range customers {
		if customer.Points > heighestPoint {
			heighestPoint = customer.Points
		}
	}
	ranks = make([][]string, heighestPoint + 1)
	for _, customer := range customers {
		name := customer.FirstName + " " + customer.LastName
		index := customer.Points

		names := ranks[index]

		ranks[index] = append(names, name)
	}
}

//Build hash table with ascii of the name mod size of customers be the hash
func BuildTable()  {



	customerTable = make([][]string, len(customers))
	pointsTable = make([][]int, len(customers))
	rankTable = make([][]int, len(customers))

	for _, customer := range customers {
		name := customer.FirstName + " " + customer.LastName
		index := SuperHash(name, len(customers))

		names := customerTable[index]

		customerTable[index] = append(names, name)

		points := pointsTable[index]

		pointsTable[index] = append(points, customer.Points)

		ranks := rankTable[index]

		rank := GetRank(name)

		rankTable[index] = append(ranks, rank)


	}
}

//Find the rank of the given customer
func GetRank(name string) int {
	rank := 0

	for point := len(ranks)-1; point >= 0; point-- {
		if len(ranks[point]) > 0 {
			rank++

			if found, _, _ := Search(ranks[point], name); found {
				return rank

			}

		}
	}

	return rank
}


//Use quick sort for sorting by name
func QuicksortByName(cs []Customer , left int, right int)  {
	if left < right {
		pivot := partitionByName(cs, left, right)
		QuicksortByName(cs, left, pivot - 1)
		QuicksortByName(cs, pivot + 1, right)
	}
}
func partitionByName(cs []Customer, left int, right int) int {
	pivot := cs[right].FirstName + " " + cs[right].LastName
	index := left - 1

	for pass := left; pass <= right-1; pass++ {
		if cs[pass].FirstName + " " + cs[pass].LastName <= pivot {
			index++

			if index < pass {
				cs[index], cs[pass] = cs[pass], cs[index]
			}
		}
	}

	index++

	cs[index], cs[right] = cs[right], cs[index]

	return index
}

//use Quick sort for sorting by Rank
func QuicksortByRank(cs []Customer , left int, right int)  {
	if left < right {
		pivot := partitionByRank(cs, left, right)
		QuicksortByRank(cs, left, pivot - 1)
		QuicksortByRank(cs, pivot + 1, right)
	}
}
func partitionByRank(cs []Customer, left int, right int) int {
	pivot := cs[right].Rank
	index := left - 1

	for pass := left; pass <= right-1; pass++ {
		if cs[pass].Rank <= pivot {
			index++

			if index < pass {
				cs[index], cs[pass] = cs[pass], cs[index]
			}
		}
	}

	index++

	cs[index], cs[right] = cs[right], cs[index]

	return index
}

//Search for a key in an array using linear search
func Search(values []string, key string) (bool, []string, []int) {
	result := make([]string, 0)
	indexes := make([]int, 0)
	found := false
	for index, value := range values {
		if value == key {
			found = true
			result = append(result, value)
			indexes = append(indexes, index)
		}
	}

	return found, result, indexes
}

//Find the hash of the value
func SuperHash(value string, dim int) int {
	hValue := AsciiSum(value) % dim

	return hValue
}

//Find the sum of Ascii characters in a string
func AsciiSum(value string) int {
	sum := 0

	for _, b := range []byte(value) {
		sum += int(b)
	}

	return sum
}

//Load all the customers from database
func loadAllCustomers()  []CustomerModel{

	db, err := gorm.Open("sqlite3", "sortomatic.db")

	defer db.Close()

	if err != nil {
		return nil
	}

	db.LogMode(true)

	var customers []CustomerModel

	customers = make([]CustomerModel, 5, 5)

	db.Find(&customers)

	return customers

}


//Search Customer by name
func (cs *CustomerService) SearchCustomer(name string)  []Customer{

	index := SuperHash(name, len(customers))

	found, _, indexes := Search(customerTable[index], name)

	if found {
		customers := make([]Customer, 0)

		for _, vIndex := range indexes {
			firstName := strings.Split(name, " ")[0]

			lastName := strings.Split(name, " ")[1]

			points := pointsTable[index][vIndex]

			rank := rankTable[index][vIndex]

			level := 1

			if points >= 5000 {
				level = 3
			} else if points >= 1000 {
				level = 2
			}

			customer := Customer{firstName, lastName, rank, level, points}
			customers = append(customers, customer)
		}

		return customers
	}


	return nil
	
}

//Sort all customers by name
func (service *CustomerService) SortByName(cms []CustomerModel) []Customer{
	cs := make([]Customer, 0)

	for _, cm := range cms {

		points := cm.Points
		rank := GetRank(cm.FirstName + " " + cm.LastName)

		level := 1

		if points >= 5000 {
			level = 3
		} else if points >= 1000 {
			level = 2
		}

		customer := Customer{cm.FirstName, cm.LastName, rank, level, points}

		cs = append(cs, customer)
	}

	QuicksortByName(cs, 0, len(cs) -1 )

	return cs
}

//Get the customer with highest point
func (service *CustomerService) GetCustomersByHighestPoint() Customer{
	for point := len(ranks)-1; point >= 0; point-- {
		if len(ranks[point]) > 0 {
			for _, customer := range service.SearchCustomer(ranks[point][0]){
				if customer.Points == point {
					return customer
				}
			}
		}
	}
	return Customer{}
}

//Get customers by the chosen level
func (service *CustomerService) GetCustomersByLevel(level string) []Customer{

	cs := make([]Customer, 0)

	for _, cm := range customers {

		points := cm.Points
		rank := GetRank(cm.FirstName + " " + cm.LastName)

		l := 1

		if level == "platinum" && points >= 5000 {
			l = 3
			customer := Customer{cm.FirstName, cm.LastName, rank, l, points}

			cs = append(cs, customer)
		} else if level == "gold" && points >= 1000 && points < 5000{
			l = 2
			customer := Customer{cm.FirstName, cm.LastName, rank, l, points}

			cs = append(cs, customer)
		}else if level == "silver" && points < 1000 {
			customer := Customer{cm.FirstName, cm.LastName, rank, l, points}

			cs = append(cs, customer)
		}



	}

	QuicksortByRank(cs, 0, len(cs) -1 )

	return cs

}

//Get the average loyalty point
func (service *CustomerService) GetAvgPoint() float32{
	sum := 0
	for _, customer := range customers {
		sum += customer.Points
	}
	return float32(sum) / float32(len(customers))
}

func (cs *CustomerService) GetAllCustomersHandler(w http.ResponseWriter, r *http.Request) {


	response := cs.SortByName(customers)

	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(response)

}

func (cs *CustomerService) SearchCustomersHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	response := cs.SearchCustomer(params["name"])

	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(response)

}

func (cs *CustomerService) GetCustomersByLevelHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	response := cs.GetCustomersByLevel(params["level"])

	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(response)

}

func (cs *CustomerService) GetCustomersByHighestPointHandler(w http.ResponseWriter, r *http.Request) {

	response := cs.GetCustomersByHighestPoint()

	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(response)

}

func (cs *CustomerService) GetAvgPointHandler(w http.ResponseWriter, r *http.Request) {

	response := cs.GetAvgPoint()

	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(response)

}