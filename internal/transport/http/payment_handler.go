package http

import (
	"encoding/json"
	"net/http"

	"Payments/internal/usecase"
)

type PaymentHandler struct {
	uc *usecase.PaymentUseCase
}

func NewPaymentHandler(uc *usecase.PaymentUseCase) *PaymentHandler {
	return &PaymentHandler{uc: uc}
}

func (h *PaymentHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc("POST /payments", h.authorize)
	mux.HandleFunc("GET /payments/{order_id}", h.getByOrderID)
}

func (h *PaymentHandler) authorize(w http.ResponseWriter, r *http.Request) {
	var req struct {
		OrderID string `json:"order_id"`
		Amount  int64  `json:"amount"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if req.OrderID == "" || req.Amount <= 0 {
		http.Error(w, "order_id and amount > 0 are required", http.StatusBadRequest)
		return
	}

	payment, err := h.uc.Authorize(r.Context(), req.OrderID, req.Amount)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"transaction_id": payment.TransactionID,
		"status":         payment.Status,
	})
}

func (h *PaymentHandler) getByOrderID(w http.ResponseWriter, r *http.Request) {
	orderID := r.PathValue("order_id")
	payment, err := h.uc.GetByOrderID(r.Context(), orderID)
	if err != nil {
		http.Error(w, "payment not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(payment)
}
