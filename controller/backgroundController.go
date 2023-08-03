package controller

import (
	"Ecommerce/config"
	"Ecommerce/model"
	"database/sql"
	"fmt"
	"github.com/go-gomail/gomail"
	"github.com/robfig/cron/v3"
	"log"
)

func StartBackgroundTask() {
	c := cron.New()
	_, err := c.AddFunc("0 0 * * *", sendReminderEmails)
	if err != nil {
		log.Fatal("Failed to schedule the cron job:", err)
	}
	c.Start()
	select {}
}

func sendReminderEmails() {
	db, err := config.OpenDatabase()
	if err != nil {
		log.Println("Failed to connect to the database:", err)
		return
	}
	defer db.Close()
	customers, err := getCustomersWithPendingOrders(db)
	if err != nil {
		log.Println("Failed to retrieve customers with pending orders:", err)
		return
	}
	for _, customer := range customers {
		err = sendReminderEmailToCustomer(customer)
		if err != nil {
			log.Println("Failed to send reminder email to customer:", err)
		}
	}
}

func getCustomersWithPendingOrders(db *sql.DB) ([]model.Customer, error) {
	query := `
		SELECT c.customer_id, c.name, c.email
		FROM customers c
		INNER JOIN orders o ON c.customer_id = o.customer_id
		WHERE o.status = 'pending'
		GROUP BY c.customer_id, c.name, c.email
	`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var customers []model.Customer
	for rows.Next() {
		var customer model.Customer
		err := rows.Scan(&customer.ID, &customer.Name, &customer.Email)
		if err != nil {
			return nil, err
		}
		customers = append(customers, customer)
	}
	return customers, nil
}

func sendReminderEmailToCustomer(customer model.Customer) error {
	subject := "Reminder: Pending Order"
	body := fmt.Sprintf("Dear %s,\n\nYou have pending order(s) in our store. Please complete the checkout process by following the link below:\n\n%s\n\nThank you for shopping with us!\n\nSincerely,\nThe Store Team", customer.Name, "https://yourstore.com/checkout")

	m := gomail.NewMessage()
	m.SetHeader("From", "yourstore@example.com")
	m.SetHeader("To", customer.Email)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	d := gomail.NewDialer("smtp.example.com", 587, "your-smtp-username", "your-smtp-password")

	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
