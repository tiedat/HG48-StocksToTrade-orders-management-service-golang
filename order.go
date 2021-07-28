package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type Order struct {
	ID        uint64    `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"create_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func getAllOrderHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	orders, err := getAllOrder()

	if err != nil {
		log.Fatalf("Unable to get all user. %v", err)
	}

	json.NewEncoder(w).Encode(orders)

	if err := json.NewEncoder(w).Encode(orders); err != nil {
		http.Error(w, "Server Error", http.StatusInternalServerError)
		return
	}
}

func getAllOrder() ([]Order, error) {
	db, err := createConnection()

	if err != nil {
		return nil, err
	}

	defer db.Close()

	sqlStatement := `SELECT id, email, created_at, updated_at FROM orders`

	rows, _ := db.Query(sqlStatement)

	if _, err := db.Query(sqlStatement); err != nil {
		return nil, err
	}

	defer rows.Close()

	var orders []Order

	for rows.Next() {
		var order Order

		err = rows.Scan(&order.ID, &order.Email, &order.CreatedAt, &order.UpdatedAt)

		if err != nil {
			return nil, err
		}

		// append the user in the users slice
		orders = append(orders, order)
	}

	if err := db.Close(); err != nil {
		return orders, err
	}

	if err := rows.Close(); err != nil {
		return orders, err
	}

	return orders, err
}
