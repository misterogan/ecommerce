package controller

import (
	"Ecommerce/config"
	"Ecommerce/model"
	"encoding/json"
	"net/http"
)

func PlaceOrder(w http.ResponseWriter, r *http.Request) {
	var order model.Order
	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if order.CustomerID == 0 {
		http.Error(w, "Customer ID is required", http.StatusBadRequest)
		return
	}

	if len(order.Items) == 0 {
		http.Error(w, "At least one item is required", http.StatusBadRequest)
		return
	}

	db, err := config.OpenDatabase()
	if err != nil {
		http.Error(w, "Failed to connect to the database", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		http.Error(w, "Failed to start the transaction", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	insertOrderSQL := "INSERT INTO orders (customer_id, status) VALUES (?, ?)"
	result, err := tx.Exec(insertOrderSQL, order.CustomerID, "pending")
	if err != nil {
		http.Error(w, "Failed to create the order", http.StatusInternalServerError)
		return
	}

	orderID, err := result.LastInsertId()
	if err != nil {
		http.Error(w, "Failed to get the order ID", http.StatusInternalServerError)
		return
	}

	insertItemSQL := "INSERT INTO order_items (order_id, product_id, quantity) VALUES (?, ?, ?)"
	for _, item := range order.Items {
		_, err = tx.Exec(insertItemSQL, orderID, item.ProductID, item.Quantity)
		if err != nil {
			http.Error(w, "Failed to add order items", http.StatusInternalServerError)
			return
		}
	}

	err = tx.Commit()
	if err != nil {
		http.Error(w, "Failed to commit the transaction", http.StatusInternalServerError)
		return
	}

	response := map[string]int64{"order_id": orderID}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
