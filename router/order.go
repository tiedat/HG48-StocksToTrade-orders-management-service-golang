package router

import (
	"encoding/json"
	db "gitlab.com/StocksToTrade2/stt-order-mgmt/database"
	"log"
	"net/http"
)

func GetAllOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	users, err := db.GetAllOrder()

	if err != nil {
		log.Fatalf("Unable to get all user. %v", err)
	}

	json.NewEncoder(w).Encode(users)

	if err:= json.NewEncoder(w).Encode(users); err != nil {
		http.Error(w, "Server Error", http.StatusInternalServerError)
		return
	}
}
