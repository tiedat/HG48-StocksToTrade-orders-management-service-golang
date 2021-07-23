package db

import (
	"gitlab.com/StocksToTrade2/stt-order-mgmt/models"
	"log"
)

func GetAllOrder() ([]models.Order, error) {
	db := createConnection()

	defer db.Close()

	sqlStatement := `SELECT id, email, created_at, updated_at FROM orders`

	rows, err := db.Query(sqlStatement)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	defer rows.Close()

	var orders []models.Order

	for rows.Next() {
		var order models.Order

		// unmarshal the row object to user
		err = rows.Scan(&order.ID, &order.Email, &order.CreatedAt, &order.UpdatedAt)

		if err != nil {
			log.Fatalf("Unable to scan the row. %v", err)
		}

		// append the user in the users slice
		orders = append(orders, order)
	}

	if err := db.Close(); err != nil {
		log.Fatalf("Unable to close db .%v", err)
	}

	if err := rows.Close(); err != nil {
		log.Fatalf("Unable to close rows .%v", err)
	}

	return orders, err
}
