package usecase

import (
	"context"

	"Payments/internal/domain"

	"github.com/google/uuid"
)

type PaymentRepository interface {
	Save(ctx context.Context, p *domain.Payment) error
	FindByOrderID(ctx context.Context, orderID string) (*domain.Payment, error)
}

type PaymentUseCase struct {
	repo PaymentRepository
}

func NewPaymentUseCase(repo PaymentRepository) *PaymentUseCase {
	return &PaymentUseCase{repo: repo}
}

func (uc *PaymentUseCase) Authorize(ctx context.Context, orderID string, amount int64) (*domain.Payment, error) {
	payStatus := domain.StatusAuthorized
	if amount > domain.MaxPaymentAmount {
		payStatus = domain.StatusDeclined
	}

	p := &domain.Payment{
		ID:            uuid.New().String(),
		OrderID:       orderID,
		TransactionID: uuid.New().String(),
		Amount:        amount,
		Status:        payStatus,
	}

	if err := uc.repo.Save(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

func (uc *PaymentUseCase) GetByOrderID(ctx context.Context, orderID string) (*domain.Payment, error) {
	return uc.repo.FindByOrderID(ctx, orderID)
}
