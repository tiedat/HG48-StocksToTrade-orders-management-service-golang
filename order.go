package main

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type Order struct {
	ID         uint64     `json:"id,omitempty"`
	Email      string     `json:"email,omitempty"`
	ProductID  int        `json:"product_id,omitempty"`
	RecurlyUID string     `json:"recurly_uid,omitempty"`
	CreatedAt  *time.Time `json:"create_at,omitempty"`
	UpdatedAt  *time.Time `json:"updated_at,omitempty"`
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

		orders = append(orders, order)
	}

	return orders, err
}

func (s *server) orderDetailHandler(w http.ResponseWriter, r *http.Request) {
	order, err := s.getOrderByEmail(chi.URLParam(r, "email"))
	if err != nil {
		logrus.WithError(err).Error("failed to get orders from db")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(order); err != nil {
		logrus.WithError(err).Error("failed to encode orders response")
		w.WriteHeader(http.StatusInternalServerError)
		return

	}
}

func (s *server) getOrderByEmail(email string) (*Order, error) {
	sqlStatement := `SELECT email, product_id FROM orders WHERE email = $1`
	row := s.db.QueryRow(sqlStatement, email)

	order := new(Order)
	if err := row.Scan(&order.Email, &order.ProductID); err != nil {
		return nil, err
	}
	return order, nil

}
