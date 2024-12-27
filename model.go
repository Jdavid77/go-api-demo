package main

import (
	"database/sql"
	"fmt"
)

type Product struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}

func getProducts(db *sql.DB) ([]Product, error) {
	query := "Select * from products"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := []Product{}

	for rows.Next() {
		var product Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Quantity, &product.Price); err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil

}

func (p *Product) getProductById(db *sql.DB) error {
	query := fmt.Sprintf("Select * from products where id = %v", p.ID)
	row := db.QueryRow(query)
	err := row.Scan(&p.ID, &p.Name, &p.Quantity, &p.Price)
	if err != nil {
		return err
	}
	return nil
}

func (p *Product) createProduct(db *sql.DB) error {
	query := fmt.Sprintf("Insert into products(name,quantity,price) values('%v', %v, %v)",p.Name, p.Quantity, p.Price)
	result, err := db.Exec(query)
	if err != nil {
		return err
	}
	
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	p.ID = int(id)
	return nil

}

func (p *Product) updateProduct(db *sql.DB) error {
	query := fmt.Sprintf("update products set name='%v', quantity=%v, price=%v where id=%v", p.Name, p.Quantity, p.Price, p.ID)
	_, err := db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func (p *Product) deleteProduct(db *sql.DB) error {
	query := fmt.Sprintf("delete from products where id=%v", p.ID)
	_, err := db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}