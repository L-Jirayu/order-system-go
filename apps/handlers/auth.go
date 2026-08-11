package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	authpb "go-app-layer/gen/auth"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token  string `json:"token"`
	UserID string `json:"user_id"`
}

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	conn, err := grpc.NewClient(
		"localhost:9001", //java-core
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		http.Error(w, "failed to connect to gRPC server", http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	client := authpb.NewAuthServiceClient(conn)

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	res, err := client.Login(ctx, &authpb.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		http.Error(w, "gRPC login failed", http.StatusInternalServerError)
		return
	}

	response := LoginResponse{
		Token:  res.GetToken(),
		UserID: res.GetUserId(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(response)
}
