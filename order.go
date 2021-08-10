package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type Order struct {
	ID         *uint64    `json:"id"`
	Email      *string    `json:"email"`
	ProductID  *int       `json:"product_id"`
	RecurlyUID *string    `json:"-"`
	CreatedAt  *time.Time `json:"create_at"`
	UpdatedAt  *time.Time `json:"updated_at"`
}

type RecurlySubscription struct {
	ID                    *uint64 `json:"id"`
	UserId                *uint64 `json:"user_id"`
	RecurlySubscriptionId *string `json:"recurly_subscription_id"`
	UtmSource             *string `json:"utm_source"`
	UtmContent            *string `json:"utm_content"`
	UtmMedium             *string `json:"utm_medium"`
	UtmTerm               *string `json:"utm_term"`
	UtmCampaign           *string `json:"utm_campaign"`
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
	email := chi.URLParam(r, "email")
	if email == "" {
		logrus.Error("path param email is empty")
		w.WriteHeader(http.StatusNotFound)
		return
	}
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

// hard code handler function for test
func (s *server) cancelSubscription(w http.ResponseWriter, r *http.Request) {
	email := chi.URLParam(r, "email")
	//todo select user have email = email and is not dev
	row := s.db.QueryRow(`SELECT id FROM users WHERE email LIKE $1 AND NOT is_dev`, email)
	var userId uint64
	if err := row.Scan(&userId); err != nil {
		json.NewEncoder(w).Encode(fmt.Sprintf("user not found for %v", email))
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// todo entitlementLog create new event with user, SUBSCRIPTION_EXPIRED, false
	// self.create user: user, event: name, ent_nasdaq: flag,
	// ent_nyse: flag, ent_otc: flag, ent_date_nasdaq: date,
	// ent_date_nyse: date, ent_date_otc: date
	s.db.Exec(`INSERT INTO entitlement_logs (ent_date_nasdaq,ent_date_nyse,ent_date_otc,created_at,updated_at,user_id,"event") VALUES (NOW(),NOW(),NOW(),NOW(),NOW(),$1,'subscription expired')`, userId)

	//todo user.nyse_entry.present? -> user.nyse_entry.update!
	row = s.db.QueryRow(`select id from nyse_entries ne where user_id = $1 order by created_at asc offset 1`, userId)

	var nyseEntryId uint64
	if err := row.Scan(&nyseEntryId); err == nil {
		s.db.Exec(`update nyse_entries set void_agreement = true, void_agreement_date = NOW() where id = $1`, nyseEntryId)
	}
	//todo user.update
	s.db.Exec(`update users set subscription_expired = true, signed_agreements = false, nyse_token = null where id = $1`, userId)

	//todo referral exist -> referral.update
	row = s.db.QueryRow(`select id from referrals where not is_cancelled and not is_expired where user_id = $1`, userId)

	var referralId uint64
	if err := row.Scan(&referralId); err == nil {
		s.db.Exec(`update referrals set date_expired = NOW(), is_expired = true where id = $1`, referralId)
	}
}

func (s *server) reactiveSubscription(w http.ResponseWriter, r *http.Request) {
	email := chi.URLParam(r, "email")
	//todo select user have email = email and is not dev
	row := s.db.QueryRow(`SELECT id FROM users WHERE email LIKE $1 AND NOT is_dev`, email)
	var userId uint64
	if err := row.Scan(&userId); err != nil {
		json.NewEncoder(w).Encode(fmt.Sprintf("user not found for %v", email))
		w.WriteHeader(http.StatusNotFound)
		return
	}

	s.db.Exec(`update users set subscription_expired = false where id = $1`, userId)
}

func (s *server) getRecurlySubscription(subscriptionId string) ([]*RecurlySubscription, error) {
	sqlStatement := `SELECT id, user_id, recurly_subscription_id, utm_source, utm_content, utm_medium, utm_term, utm_campaign FROM recurly_subs where recurly_subscription_id = $1`

	rows, err := s.db.Query(sqlStatement, subscriptionId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subs []*RecurlySubscription
	for rows.Next() {
		sub := new(RecurlySubscription)
		if err := rows.Scan(&sub.ID, &sub.UserId, &sub.RecurlySubscriptionId, &sub.UtmSource, &sub.UtmContent, &sub.UtmMedium, &sub.UtmTerm, &sub.UtmCampaign); err != nil {
			return nil, err
		}
		subs = append(subs, sub)
	}

	return subs, err
}
