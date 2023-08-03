// main.go
package main

import (
	"Ecommerce/controller"
	"Ecommerce/ratelimiter"
	"net/http"
	"time"
)

func main() {
	limiter := ratelimiter.NewLimiter(100, time.Minute)
	http.HandleFunc("/api/place-order", limiter.Middleware(controller.PlaceOrder))
	http.HandleFunc("/api/view-orders", limiter.Middleware(controller.GetCustomerOrders))
	http.HandleFunc("/api/generate-report", limiter.Middleware(controller.GenerateCSVReport))
	controller.StartBackgroundTask()
	http.ListenAndServe(":8080", nil)
}
