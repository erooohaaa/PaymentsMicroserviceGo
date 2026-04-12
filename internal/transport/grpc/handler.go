package grpc

import (
	"Payments/internal/usecase"
	"context"
	"github.com/erooohaaa/orders-generated/api"
)

type PaymentGRPCHandler struct {
	api.UnimplementedPaymentServiceServer
	uc *usecase.PaymentUseCase
}

func NewPaymentGRPCHandler(uc *usecase.PaymentUseCase) *PaymentGRPCHandler {
	return &PaymentGRPCHandler{uc: uc}
}

func (h *PaymentGRPCHandler) ProcessPayment(ctx context.Context, req *api.PaymentRequest) (*api.PaymentResponse, error) {
	// Вызываем существующий UseCase из первого асика
	payment, err := h.uc.Authorize(ctx, req.OrderId, req.Amount)
	if err != nil {
		return nil, err
	}

	return &api.PaymentResponse{
		TransactionId: payment.TransactionID,
		Status:        payment.Status,
	}, nil
}
