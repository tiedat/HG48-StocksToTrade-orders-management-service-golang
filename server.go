package main

import (
	"database/sql"
	"net/http"
)

type server struct {
	hs *http.Server

	addr string
	db   *sql.DB
}

type option func(*server)

func withDB(db *sql.DB) option {
	return func(s *server) {
		s.db = db
	}
}

func withAddr(addr string) option {
	return func(s *server) {
		s.addr = addr
	}
}

func newServer(opts ...option) (*server, error) {
	s := &server{addr: "0.0.0.0:3333"}

	for _, opt := range opts {
		opt(s)
	}

	s.hs = &http.Server{Addr: s.addr}
	s.setupRouter()

	return s, nil
}
