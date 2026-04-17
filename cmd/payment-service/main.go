package main

import (
	"database/sql"
	"log"
	"net"
	"os"
	"time"

	"Payments/internal/repository"
	transportGRPC "Payments/internal/transport/grpc"
	"Payments/internal/usecase"

	api "github.com/erooohaaa/orders-generated"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment")
	}

	dsn := os.Getenv("DB_DSN")
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("failed to open db: %v", err)
	}
	defer db.Close()

	db.SetMaxOpenConns(10)
	db.SetConnMaxLifetime(time.Minute * 3)

	paymentRepo := repository.NewPostgresPaymentRepository(db)
	paymentUC := usecase.NewPaymentUseCase(paymentRepo)
	grpcHandler := transportGRPC.NewPaymentGRPCHandler(paymentUC)

	grpcPort := os.Getenv("GRPC_PORT")
	if grpcPort == "" {
		grpcPort = "50051"
	}

	lis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer(
		grpc.UnaryInterceptor(transportGRPC.LoggingInterceptor),
	)
	api.RegisterPaymentServiceServer(server, grpcHandler)

	log.Printf("Payment gRPC service listening on :%s", grpcPort)
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
