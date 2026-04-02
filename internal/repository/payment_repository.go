package repository

import (
	"context"

	"Payments/internal/domain"
)

type PaymentRepository interface {
	Save(ctx context.Context, p *domain.Payment) error
	FindByOrderID(ctx context.Context, orderID string) (*domain.Payment, error)
}
