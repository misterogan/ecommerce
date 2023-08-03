package controller

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"Ecommerce/config"
	"Ecommerce/model"
)

func GetCustomerOrders(w http.ResponseWriter, r *http.Request) {

	customerIDStr := r.URL.Query().Get("customer_id")
	if customerIDStr == "" {
		http.Error(w, "Customer ID is required", http.StatusBadRequest)
		return
	}

	// Convert the customer ID from string to int64
	customerID, err := strconv.ParseInt(customerIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid customer ID", http.StatusBadRequest)
		return
	}

	// Open the database connection
	db, err := config.OpenDatabase()
	if err != nil {
		http.Error(w, "Failed to connect to the database", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Retrieve the list of orders associated with the customer
	orders, err := getOrdersByCustomerID(db, customerID)
	if err != nil {
		http.Error(w, "Failed to retrieve orders", http.StatusInternalServerError)
		return
	}

	// For each order, retrieve the list of products associated with the order
	for i := range orders {
		orders[i].Items, err = getOrderItemsByOrderID(db, orders[i].ID)
		if err != nil {
			http.Error(w, "Failed to retrieve order items", http.StatusInternalServerError)
			return
		}
	}

	// Respond with the list of orders
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}

// Helper functions to retrieve data from the database
func getOrdersByCustomerID(db *sql.DB, customerID int64) ([]model.Order, error) {
	query := "SELECT * FROM orders WHERE customer_id = ?"
	rows, err := db.Query(query, customerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []model.Order
	for rows.Next() {
		var order model.Order
		err := rows.Scan(&order.ID, &order.CustomerID, &order.OrderDate, &order.Status)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	return orders, nil
}

func getOrderItemsByOrderID(db *sql.DB, orderID int64) ([]model.OrderItem, error) {
	query := "SELECT product_id, quantity FROM order_items WHERE order_id = ?"
	rows, err := db.Query(query, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []model.OrderItem
	for rows.Next() {
		var item model.OrderItem
		err := rows.Scan(&item.ProductID, &item.Quantity)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}
