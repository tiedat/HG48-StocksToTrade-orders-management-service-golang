package main

import (
	"log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	_, err := createConnection()
	if err != nil {
		log.Fatal("DB", err)
	}

	os.Exit(m.Run())
}
