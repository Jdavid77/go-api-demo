package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

var a App

func createTable () {
	query := `CREATE TABLE IF NOT EXISTS products (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(50) NOT NULL,
		quantity INT NOT NULL,
		price FLOAT NOT NULL
	)`
	_, err := a.DB.Exec(query)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func TestMain(m *testing.M) {
	err := a.Initialize(DbUsernameTest, DbPasswordTest, "test")
	if err != nil {
		log.Fatal(err.Error())
	}
	createTable()
	m.Run()
}

func clearTable() {
	a.DB.Exec("DELETE FROM products")
	a.DB.Exec("ALTER TABLE products AUTO_INCREMENT = 1")
	log.Printf("Table cleared")
}

func addProduct(name string, quantity int, price float64) {
	query := fmt.Sprintf("INSERT INTO products(name, quantity, price) VALUES('%v', %v, %v)", name, quantity, price)
	_, err := a.DB.Exec(query)
	if err != nil {
		log.Println(err.Error())
	}
} 

func TestGetProduct(t *testing.T) {
	clearTable()
	addProduct("test-product", 100, 50.5)
	request, _ := http.NewRequest("GET", "/product/1", nil)
	response := sendRequest(request)
	checkStatusCode(t, http.StatusOK, response.Code)
}

func TestCreateProduct(t *testing.T) {
	clearTable()
	payload := []byte(`{"name":"TV","quantity":100,"price":200}`)
	request, _ := http.NewRequest("POST", "/product", bytes.NewBuffer(payload))
	request.Header.Set("Content-Type", "application/json")
	response := sendRequest(request)
	checkStatusCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["name"] != "TV" {
		t.Errorf("Expected product name to be 'TV'. Got '%v'", m["name"])
	}
	if m["quantity"] != 100.0 {
		t.Errorf("Expected product quantity to be 100. Got '%v'", m["quantity"])
	}
	if m["price"] != 200.0 {	
		t.Errorf("Expected product price to be 200. Got '%v'", m["price"])
	}

}

func TestDeleteProduct(t *testing.T) {
	clearTable()
	addProduct("table", 100, 100)
	request, _ := http.NewRequest("GET", "/product/1", nil)
	response := sendRequest(request)
	checkStatusCode(t, http.StatusOK, response.Code)
	
	request, _ = http.NewRequest("DELETE", "/product/1", nil)
	response = sendRequest(request)
	checkStatusCode(t, http.StatusOK, response.Code)

	request, _ = http.NewRequest("GET", "/product/1", nil)
	response = sendRequest(request)
	checkStatusCode(t, http.StatusNotFound, response.Code)

}

func TestUpdateProduct(t *testing.T) {
	clearTable()
	addProduct("table", 100, 100)
	req, _ := http.NewRequest("GET", "/product/1", nil)
	response := sendRequest(req)

	var oldValues map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &oldValues)

	var payload = []byte(`{"name":"Router","quantity":50,"price":20}`)
	req, _ = http.NewRequest("PUT", "/product/1", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")

	response = sendRequest(req)
	var newValues map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &newValues)

	if newValues["id"] != oldValues["id"] {
		t.Errorf("Expected id: %v,  Got %v", newValues["id"], oldValues["id"])
	}

	if newValues["price"] == oldValues["price"] {
		t.Errorf("Expected price: %v,  Got %v", newValues["price"], oldValues["price"])
	}

	if newValues["quantity"] == oldValues["quantity"] {
		t.Errorf("Expected quantity: %v,  Got %v", newValues["quantity"], oldValues["quantity"])
	}

	if newValues["name"] == oldValues["name"] {
		t.Errorf("Expected name to be different. Got %v", newValues["name"])
	}


}

func checkStatusCode(t *testing.T, expectedStatusCode, actualStatusCode int) {
	if expectedStatusCode != actualStatusCode {
		t.Errorf("Expected status code %v but got %v", expectedStatusCode, actualStatusCode)
	}
}

func sendRequest(request *http.Request) *httptest.ResponseRecorder{
	recorder := httptest.NewRecorder()
	a.Router.ServeHTTP(recorder, request)
	return recorder
}