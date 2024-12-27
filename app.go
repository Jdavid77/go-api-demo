package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)


type App struct {
	Router *mux.Router
	DB     *sql.DB
}

func (app *App) Initialize(DbUsername, DbPassword, DbName string) error {
	connectionString := fmt.Sprintf("%v:%v@tcp(127.0.0.1:3306)/%v", DbUsername, DbPassword, DbName)
    var err error
	app.DB, err = sql.Open("mysql", connectionString)

	if err != nil {
		log.Fatal(err.Error())
		return err
	}

	app.Router = mux.NewRouter().StrictSlash(true)
	app.handleRoutes()
	return nil
}

func (app *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, app.Router))
}

func sendResponse(w http.ResponseWriter,status int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(response)
}

func sendError(w http.ResponseWriter,status int, err string) {
	error_message := map[string]string{"error": err}
	sendResponse(w, status, error_message)

}

func (app *App) getProducts(w http.ResponseWriter, r *http.Request) {
	products, err := getProducts(app.DB)
	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}
	sendResponse(w, http.StatusOK, products)
	
}

func (app *App) getProductById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var err error
	key, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendError(w, http.StatusBadRequest, "Invalid Product ID")
		return
	}
	product := Product{ID: key}
	err = product.getProductById(app.DB)
	if err != nil {
		switch err {
			case sql.ErrNoRows:
				sendError(w, http.StatusNotFound, "Product not found")
				
			default:
				sendError(w, http.StatusInternalServerError, err.Error())
				
		}
		return
	}

	sendResponse(w, http.StatusOK, product)
	
}

func (app *App) createProduct(w http.ResponseWriter, r *http.Request) {
	var product Product
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&product); err != nil {
		sendError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()
	if err := product.createProduct(app.DB); err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}
	sendResponse(w, http.StatusCreated, product)
	
}

func (app *App) updateProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendError(w, http.StatusBadRequest, "Invalid Product ID")
		return
	}
	var product Product
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&product); err != nil {
		sendError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()
	product.ID = key
	if err := product.updateProduct(app.DB); err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}
	sendResponse(w, http.StatusOK, product)
	
}

func (app *App) deleteProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendError(w, http.StatusBadRequest, "Invalid Product ID")
		return
	}
	product := Product{ID: key}
	if err := product.deleteProduct(app.DB); err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}
	sendResponse(w, http.StatusOK, map[string]string{"result": "success"})
	
}

func (app *App) handleRoutes() {
	app.Router.HandleFunc("/products", app.getProducts).Methods("GET")
	app.Router.HandleFunc("/product/{id}", app.getProductById).Methods("GET")
	app.Router.HandleFunc("/product", app.createProduct).Methods("POST")
	app.Router.HandleFunc("/product/{id}", app.updateProduct).Methods("PUT")
	app.Router.HandleFunc("/product/{id}", app.deleteProduct).Methods("DELETE")
}