package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	paymentpb "go-app-layer/gen/payment"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// ========================================
// gRPC Payment Client
// ========================================

func getPaymentClient() (paymentpb.PaymentServiceClient, *grpc.ClientConn, error) {
	conn, err := grpc.NewClient(
		"localhost:9001",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, nil, err
	}

	return paymentpb.NewPaymentServiceClient(conn), conn, nil
}

// ========================================
// POST /payments
// ========================================

func CreatePayment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		UserID string  `json:"user_id"`
		Amount float64 `json:"amount"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	client, conn, err := getPaymentClient()
	if err != nil {
		http.Error(
			w,
			"failed to connect to gRPC server: "+err.Error(),
			http.StatusInternalServerError,
		)
		return
	}
	defer conn.Close()

	res, err := client.CreatePayment(
		ctx,
		&paymentpb.CreatePaymentRequest{
			UserId: req.UserID,
			Amount: req.Amount,
		},
	)
	if err != nil {
		http.Error(
			w,
			"gRPC create payment failed: "+err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(res)
}

// ========================================
// GET /payments/{payment_id}
// ========================================

func GetPayment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Expected:
	// /payments/PAY-123
	//
	// Split:
	// ["payments", "PAY-123"]

	parts := strings.Split(
		strings.Trim(r.URL.Path, "/"),
		"/",
	)

	if len(parts) != 2 || parts[0] != "payments" {
		http.Error(
			w,
			"invalid payment path",
			http.StatusBadRequest,
		)
		return
	}

	paymentID := parts[1]

	if paymentID == "" {
		http.Error(
			w,
			"payment_id is required",
			http.StatusBadRequest,
		)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	client, conn, err := getPaymentClient()
	if err != nil {
		http.Error(
			w,
			"failed to connect to gRPC server: "+err.Error(),
			http.StatusInternalServerError,
		)
		return
	}
	defer conn.Close()

	res, err := client.GetPaymentStatus(
		ctx,
		&paymentpb.GetPaymentStatusRequest{
			PaymentId: paymentID,
		},
	)
	if err != nil {
		http.Error(
			w,
			"gRPC get payment status failed: "+err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(res)
}
