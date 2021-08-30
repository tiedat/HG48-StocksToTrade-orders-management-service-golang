package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/gchaincl/dotsql"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/ory/dockertest"
)

var db *sql.DB

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	resource, err := pool.Run("postgres", "13", []string{"POSTGRES_PASSWORD=123456"})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	dbHost := "localhost"
	if postgresDBHost := os.Getenv("POSTGRES_DB_HOST"); postgresDBHost != "" {
		dbHost = postgresDBHost
	}
	if err := pool.Retry(func() error {
		var err error
		db, err = sql.Open("pgx", fmt.Sprintf("postgres://postgres:123456@%s:%s/postgres", dbHost, resource.GetPort("5432/tcp")))
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	dot, err := dotsql.LoadFromFile("db_test.sql")
	if err != nil {
		log.Fatal(err)
	}
	for _, sqlQueryName := range []string{
		"create_citext_extension",
		"create_orders_table",
		"create_email_index",
		"create_orders_data",
		"create_users_table",
		"create_users_data",
		"create_entitlement_logs_table",
		"create_nyse_entries_table",
		"create_nyse_entries_data",
		"create_referrals_table",
		"create_referrals_data",
		"create_recurly_subs_table",
		"create_recurly_subs_data",
	} {
		if _, err := dot.Exec(db, sqlQueryName); err != nil {
			log.Fatalf("failed to exectue sql command: %s, error: %v", sqlQueryName, err)
		}
	}

	code := m.Run()

	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}
