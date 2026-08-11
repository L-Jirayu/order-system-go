package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	orderpb "go-app-layer/gen/order"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// ========================================
// gRPC Order Client
// ========================================

func getOrderClient() (orderpb.OrderServiceClient, *grpc.ClientConn, error) {
	conn, err := grpc.NewClient(
		"localhost:9001",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, nil, err
	}

	return orderpb.NewOrderServiceClient(conn), conn, nil
}

// ========================================
// POST /orders
// ========================================

func CreateOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// HTTP JSON request
	var req struct {
		UserID       string `json:"user_id"`
		RestaurantID string `json:"restaurant_id"`
		Items        []struct {
			ItemID   string `json:"item_id"`
			Quantity int32  `json:"quantity"`
		} `json:"items"`
		AddressID string `json:"address_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	client, conn, err := getOrderClient()
	if err != nil {
		http.Error(
			w,
			"failed to connect to gRPC server: "+err.Error(),
			http.StatusInternalServerError,
		)
		return
	}
	defer conn.Close()

	// Convert HTTP JSON items -> protobuf OrderItemInput
	items := make([]*orderpb.OrderItemInput, 0, len(req.Items))

	for _, item := range req.Items {
		items = append(items, &orderpb.OrderItemInput{
			ItemId:   item.ItemID,
			Quantity: item.Quantity,
		})
	}

	// Call Java gRPC service
	res, err := client.CreateOrder(
		ctx,
		&orderpb.CreateOrderRequest{
			UserId:       req.UserID,
			RestaurantId: req.RestaurantID,
			Items:        items,
			AddressId:    req.AddressID,
		},
	)
	if err != nil {
		http.Error(
			w,
			"gRPC create order failed: "+err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(res)
}

// ========================================
// GET /orders/{order_id}
// ========================================

func GetOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Expected:
	// /orders/ORDER-123
	//
	// Split:
	// ["orders", "ORDER-123"]

	parts := strings.Split(
		strings.Trim(r.URL.Path, "/"),
		"/",
	)

	if len(parts) != 2 || parts[0] != "orders" {
		http.Error(
			w,
			"invalid order path",
			http.StatusBadRequest,
		)
		return
	}

	orderID := parts[1]

	if orderID == "" {
		http.Error(
			w,
			"order_id is required",
			http.StatusBadRequest,
		)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	client, conn, err := getOrderClient()
	if err != nil {
		http.Error(
			w,
			"failed to connect to gRPC server: "+err.Error(),
			http.StatusInternalServerError,
		)
		return
	}
	defer conn.Close()

	res, err := client.GetOrderStatus(
		ctx,
		&orderpb.GetOrderStatusRequest{
			OrderId: orderID,
		},
	)
	if err != nil {
		http.Error(
			w,
			"gRPC get order status failed: "+err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(res)
}
