package service

import "testing"

func TestCustomerService_GetAllCustomers(t *testing.T) {
	cs := CustomerService{}

	customers := cs.SearchCustomer("Dhyan Raj")

	if len(customers) < 1 {
		t.Error("Couldnt load customers")
	}
}



func TestAsciiSum(t *testing.T) {
	sum := AsciiSum("Dhyan Raj")

	if sum < 1 {
		t.Error("Wrong integer value of string")
	}
}