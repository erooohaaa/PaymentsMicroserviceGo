package main

import (
	"database/sql"
	"log"
	"net"
	"os"

	"Payments/internal/repository"
	transportGRPC "Payments/internal/transport/grpc"
	"Payments/internal/usecase"
	"github.com/erooohaaa/orders-generated/api"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using defaults")
	}

	dsn := os.Getenv("DB_DSN")
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("failed to open db: %v", err)
	}
	defer db.Close()

	// Инициализация слоев Clean Architecture
	paymentRepo := repository.NewPostgresPaymentRepository(db)
	paymentUC := usecase.NewPaymentUseCase(paymentRepo)
	grpcHandler := transportGRPC.NewPaymentGRPCHandler(paymentUC)

	// Запуск grpc сервера
	grpcPort := os.Getenv("GRPC_PORT") // Добавь в .env: GRPC_PORT=50051
	if grpcPort == "" {
		grpcPort = "50051"
	}

	lis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer()
	api.RegisterPaymentServiceServer(server, grpcHandler)

	log.Printf("Payment grpc Service listening on %v", grpcPort)
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve grpc: %v", err)
	}
}
