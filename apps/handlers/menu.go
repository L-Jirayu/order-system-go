package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	menupb "go-app-layer/gen/menu"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// ========================================
// gRPC Menu Client
// ========================================

func getMenuClient() (menupb.MenuServiceClient, *grpc.ClientConn, error) {
	conn, err := grpc.NewClient(
		"localhost:9001",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, nil, err
	}

	return menupb.NewMenuServiceClient(conn), conn, nil
}

// ========================================
// GET /restaurants
// ========================================

func Restaurants(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	client, conn, err := getMenuClient()
	if err != nil {
		http.Error(
			w,
			"failed to connect to gRPC server: "+err.Error(),
			http.StatusInternalServerError,
		)
		return
	}
	defer conn.Close()

	res, err := client.ListRestaurants(
		ctx,
		&menupb.ListRestaurantsRequest{},
	)
	if err != nil {
		http.Error(
			w,
			"gRPC list restaurants failed: "+err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(res)
}

// ========================================
// GET /restaurants/{restaurant_id}/menu
// ========================================

func RestaurantMenu(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Expected:
	// /restaurants/R001/menu
	//
	// Split:
	// ["", "restaurants", "R001", "menu"]

	parts := strings.Split(
		strings.Trim(r.URL.Path, "/"),
		"/",
	)

	if len(parts) != 3 ||
		parts[0] != "restaurants" ||
		parts[2] != "menu" {

		http.Error(
			w,
			"invalid restaurant menu path",
			http.StatusBadRequest,
		)
		return
	}

	restaurantID := parts[1]

	if restaurantID == "" {
		http.Error(
			w,
			"restaurant_id is required",
			http.StatusBadRequest,
		)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	client, conn, err := getMenuClient()
	if err != nil {
		http.Error(
			w,
			"failed to connect to gRPC server: "+err.Error(),
			http.StatusInternalServerError,
		)
		return
	}
	defer conn.Close()

	res, err := client.ListMenuItems(
		ctx,
		&menupb.ListMenuItemsRequest{
			RestaurantId: restaurantID,
		},
	)
	if err != nil {
		http.Error(
			w,
			"gRPC list menu items failed: "+err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(res)
}
