package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

type Order struct {
	ID         uint64    `json:"id"`
	Email      string    `json:"email"`
	ProductID  int       `json:"-"`
	RecurlyUID string    `json:"-"`
	CreatedAt  time.Time `json:"create_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (s *server) ordersHandler(w http.ResponseWriter, r *http.Request) {
	orders, err := s.orders()
	if err != nil {
		logrus.WithError(err).Error("failed to get orders from db")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(orders); err != nil {
		logrus.WithError(err).Error("failed to encode orders response")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (s *server) orders() ([]*Order, error) {
	sqlStatement := `SELECT id, email, created_at, updated_at FROM orders`

	rows, err := s.db.Query(sqlStatement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*Order
	for rows.Next() {
		order := new(Order)
		if err := rows.Scan(&order.ID, &order.Email, &order.CreatedAt, &order.UpdatedAt); err != nil {
			return nil, err
		}

		// append the user in the users slice
		orders = append(orders, order)
	}

	return orders, err
}
