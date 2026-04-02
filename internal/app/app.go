package app

import (
	"Payments/internal/repository"
	transportHTTP "Payments/internal/transport/http"
	"Payments/internal/usecase"
	"database/sql"
	"net/http"
)

type App struct {
	Server *http.Server
}

func New(db *sql.DB) *App {
	paymentRepo := repository.NewPostgresPaymentRepository(db)
	paymentUC := usecase.NewPaymentUseCase(paymentRepo)

	handler := transportHTTP.NewPaymentHandler(paymentUC)
	mux := http.NewServeMux()
	handler.Register(mux)

	return &App{
		Server: &http.Server{
			Addr:    ":8081",
			Handler: mux,
		},
	}
}
