package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/jackc/pgx/v4"
)

func main() {
	mux := http.NewServeMux()
	srv := &http.Server{Addr: ":" + os.Getenv("ABX_PORT"), Handler: mux}
	mux.HandleFunc("/", index)
	err := srv.ListenAndServe()
	if err != http.ErrServerClosed {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}
}

func index(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	s, err := query()
	if err != nil {
		s = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
	}
	fmt.Fprintf(w, "%s\n", s)
}

func query() (string, error) {
	var s string
	ctx := context.Background()
	dsn := os.Getenv("ABX_STORE_DSN")
	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		return "", err
	}
	defer conn.Close(ctx)
	name := os.Getenv("ABX_NAME")
	err = conn.QueryRow(ctx, "SELECT $1", name).Scan(&s)
	if err != nil {
		return "", err
	}
	return s, nil
}
