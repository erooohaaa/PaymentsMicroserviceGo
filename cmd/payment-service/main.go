package main

import (
	"Payments/internal/app"
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	dsn := "host=localhost user=postgres password=0000 dbname=payment_db sslmode=disable"

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("failed to open db: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	log.Println("Payment Service: connected to database")

	a := app.New(db)
	log.Println("Payment Service: listening on :8081")
	if err := a.Server.ListenAndServe(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
