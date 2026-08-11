package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	userpb "go-app-layer/gen/user"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func getUserClient() (userpb.UserServiceClient, *grpc.ClientConn, error) {
	conn, err := grpc.NewClient(
		"localhost:9001", //java-core
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, nil, err
	}

	return userpb.NewUserServiceClient(conn), conn, nil
}

// GET /profile
func Profile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	client, conn, err := getUserClient()
	if err != nil {
		http.Error(
			w,
			"failed to connect to gRPC server: "+err.Error(),
			http.StatusInternalServerError,
		)
		return
	}
	defer conn.Close()

	res, err := client.GetProfile(ctx, &userpb.GetProfileRequest{
		UserId: "user-001",
	})
	if err != nil {
		http.Error(
			w,
			"gRPC get profile failed: "+err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(res)
}

// GET /addresses
func Addresses(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	client, conn, err := getUserClient()
	if err != nil {
		http.Error(
			w,
			"failed to connect to gRPC server: "+err.Error(),
			http.StatusInternalServerError,
		)
		return
	}
	defer conn.Close()

	res, err := client.ListAddresses(ctx, &userpb.ListAddressesRequest{
		UserId: "user-001",
	})
	if err != nil {
		http.Error(
			w,
			"gRPC list addresses failed: "+err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(res)
}

// POST /addresses
func AddAddress(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		UserID      string `json:"user_id"`
		FullAddress string `json:"full_address"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	client, conn, err := getUserClient()
	if err != nil {
		http.Error(
			w,
			"failed to connect to gRPC server: "+err.Error(),
			http.StatusInternalServerError,
		)
		return
	}
	defer conn.Close()

	res, err := client.AddAddress(ctx, &userpb.AddAddressRequest{
		UserId:      req.UserID,
		FullAddress: req.FullAddress,
	})
	if err != nil {
		http.Error(
			w,
			"gRPC add address failed: "+err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(res)
}
