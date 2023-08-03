package controller

import (
	"database/sql"
	"encoding/csv"
	"net/http"
	"os"
	"strconv"

	"Ecommerce/config"
	"Ecommerce/model"
)

func GenerateCSVReport(w http.ResponseWriter, r *http.Request) {
	db, err := config.OpenDatabase()
	if err != nil {
		http.Error(w, "Failed to connect to the database", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	orders, err := getAllOrders(db)
	if err != nil {
		http.Error(w, "Failed to retrieve orders", http.StatusInternalServerError)
		return
	}

	file, err := os.Create("orders_report.csv")
	if err != nil {
		http.Error(w, "Failed to create CSV file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{"Order ID", "Customer Name", "Order Date", "Total Price", "Status"})

	for _, order := range orders {
		customerName, err := getCustomerNameByID(db, order.CustomerID)
		if err != nil {
			http.Error(w, "Failed to retrieve customer name", http.StatusInternalServerError)
			return
		}

		totalPrice, err := getTotalPriceForOrder(db, order.ID)
		if err != nil {
			http.Error(w, "Failed to retrieve total price", http.StatusInternalServerError)
			return
		}

		writer.Write([]string{
			strconv.FormatInt(order.ID, 10),
			customerName,
			order.OrderDate,
			strconv.FormatFloat(totalPrice, 'f', 2, 64),
			order.Status,
		})
	}

	w.Write([]byte("CSV report generated successfully"))
}

func getAllOrders(db *sql.DB) ([]model.Order, error) {
	query := "SELECT * FROM orders"
	rows, err := db.Query(query)
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

func getCustomerNameByID(db *sql.DB, customerID int64) (string, error) {
	var customerName string
	query := "SELECT name FROM customers WHERE customer_id = ?"
	err := db.QueryRow(query, customerID).Scan(&customerName)
	if err != nil {
		return "", err
	}
	return customerName, nil
}

func getTotalPriceForOrder(db *sql.DB, orderID int64) (float64, error) {
	var totalPrice float64
	query := "SELECT SUM(price) FROM products JOIN order_items ON products.product_id = order_items.product_id WHERE order_items.order_id = ?"
	err := db.QueryRow(query, orderID).Scan(&totalPrice)
	if err != nil {
		return 0, err
	}
	return totalPrice, nil
}
