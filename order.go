package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
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
	ID                    uint64 `json:"id"`
	UserId                uint64 `json:"user_id"`
	RecurlySubscriptionId string `json:"recurly_subscription_id"`
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

	s.writeBody(w, orders)
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

	if errors.Is(err, sql.ErrNoRows) {
		w.WriteHeader(http.StatusNotFound)
		s.writeBody(w, "email not found")
		return
	}

	if err != nil {
		logrus.WithError(err).Error("failed to get orders from db")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	s.writeBody(w, order)
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
func (s *server) cancelSubscription(email string) error {

	//select user have email = email and is not dev
	var userId uint64
	getUserErr := s.db.QueryRow(`SELECT id FROM users WHERE email LIKE $1 AND NOT is_dev`, email).Scan(&userId)

	if errors.Is(getUserErr, sql.ErrNoRows) {
		return errors.New(fmt.Sprintf("user not found for %v", email))
	} else if getUserErr != nil {
		return getUserErr
	}

	// entitlementLog create new event with user, SUBSCRIPTION_EXPIRED, false
	// self.create user: user, event: name, ent_nasdaq: flag,
	// ent_nyse: flag, ent_otc: flag, ent_date_nasdaq: date,
	// ent_date_nyse: date, ent_date_otc: date
	_, insertEntitlementLogErr := s.db.Exec(`INSERT INTO entitlement_logs (ent_date_nasdaq,ent_date_nyse,ent_date_otc,created_at,updated_at,user_id,"event") VALUES (NOW(),NOW(),NOW(),NOW(),NOW(),$1,'subscription expired')`, userId)
	if insertEntitlementLogErr != nil {
		return insertEntitlementLogErr
	}

	// user.nyse_entry.present? -> user.nyse_entry.update!
	var nyseEntryId uint64
	getNyseEntryErr := s.db.QueryRow(`SELECT id FROM nyse_entries ne WHERE user_id = $1 ORDER BY created_at DESC LIMIT 1`, userId).Scan(&nyseEntryId)

	if getNyseEntryErr == nil {
		_, updateNyseErr := s.db.Exec(`UPDATE nyse_entries SET void_agreement = TRUE, void_agreement_date = NOW() WHERE id = $1`, nyseEntryId)
		if updateNyseErr != nil {
			return updateNyseErr
		}
	} else if !errors.Is(getNyseEntryErr, sql.ErrNoRows) {
		return getNyseEntryErr
	}
	// user.update
	_, updateUserErr := s.db.Exec(`UPDATE users SET subscription_expired = TRUE, signed_agreements = FALSE, nyse_token = NULL WHERE id = $1`, userId)
	if updateUserErr != nil {
		return updateUserErr
	}
	// referral exist -> referral.update
	var referralId uint64
	getReferralErr := s.db.QueryRow(`SELECT id FROM referrals WHERE NOT is_cancelled AND NOT is_expired AND user_id = $1 `, userId).Scan(&referralId)

	if getReferralErr == nil {
		_, updateReferralErr := s.db.Exec(`UPDATE referrals SET date_expired = NOW(), is_expired = TRUE WHERE id = $1`, referralId)
		if updateReferralErr != nil {
			return updateReferralErr
		}
	} else if !errors.Is(getReferralErr, sql.ErrNoRows) {
		return getReferralErr
	}
	return nil
}

func (s *server) reactiveSubscription(email string) error {
	// select user have email = email and is not dev
	var userId uint64
	getUserErr := s.db.QueryRow(`SELECT id FROM users WHERE email LIKE $1 AND NOT is_dev`, email).Scan(&userId)

	if errors.Is(getUserErr, sql.ErrNoRows) {
		return errors.New(fmt.Sprintf("user not found for %v", email))
	} else if getUserErr != nil {
		return getUserErr
	}

	_, updateUserErr := s.db.Exec(`UPDATE users SET subscription_expired = FALSE WHERE id = $1`, userId)
	if updateUserErr != nil {
		return updateUserErr
	}

	return nil
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
		if err := rows.Scan(&sub.ID,
			&sub.UserId,
			&sub.RecurlySubscriptionId,
			&sub.UtmSource,
			&sub.UtmContent,
			&sub.UtmMedium,
			&sub.UtmTerm,
			&sub.UtmCampaign); err != nil {
			return nil, err
		}
		subs = append(subs, sub)
	}

	return subs, err
}

// always return func after writeBody
func (s *server) writeBody(w http.ResponseWriter, b interface{}) {
	if err := json.NewEncoder(w).Encode(b); err != nil {
		logrus.WithError(err).Error("failed to encode response")
		w.WriteHeader(http.StatusInternalServerError)
	}
}
