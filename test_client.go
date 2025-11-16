package main

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	authpb "go-app-layer/gen/auth"
	menupb "go-app-layer/gen/menu"
	orderpb "go-app-layer/gen/order"
	paymentpb "go-app-layer/gen/payment"
	userpb "go-app-layer/gen/user"
)

// ===== Async Wrapper =====
func runTestClientAsync() <-chan string {
	resultChan := make(chan string)

	go func() {
		defer close(resultChan)
		resultChan <- runTestClient()
	}()

	return resultChan
}

// ===== gRPC calls =====
func runTestClient() string {
	// gRPC dailer ใหม่ เลิกใช้ grpc.Dial()
	conn, err := grpc.NewClient(
		"java-core:9001",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		return fmt.Sprintf("Connection failed: %v", err)
	}
	defer conn.Close()

	result := ""

	// === AUTH ===
	authClient := authpb.NewAuthServiceClient(conn)
	authRes, _ := authClient.Login(context.Background(), &authpb.LoginRequest{
		Email:    "a@b.com",
		Password: "1234",
	})
	result += fmt.Sprintf("AUTH: %v\n", authRes)

	// === USER ===
	userClient := userpb.NewUserServiceClient(conn)
	profile, _ := userClient.GetProfile(context.Background(),
		&userpb.GetProfileRequest{UserId: "user-001"})
	result += fmt.Sprintf("USER PROFILE: %v\n", profile)

	// === MENU ===
	menuClient := menupb.NewMenuServiceClient(conn)
	rests, _ := menuClient.ListRestaurants(context.Background(),
		&menupb.ListRestaurantsRequest{})
	result += fmt.Sprintf("RESTAURANTS: %v\n", rests)

	items, _ := menuClient.ListMenuItems(context.Background(),
		&menupb.ListMenuItemsRequest{RestaurantId: "R001"})
	result += fmt.Sprintf("MENU ITEMS: %v\n", items)

	// === PAYMENT ===
	payClient := paymentpb.NewPaymentServiceClient(conn)
	payRes, _ := payClient.CreatePayment(context.Background(),
		&paymentpb.CreatePaymentRequest{
			UserId: "user-001",
			Amount: 150.50,
		})
	result += fmt.Sprintf("PAYMENT: %v\n", payRes)

	// === ORDER ===
	orderClient := orderpb.NewOrderServiceClient(conn)
	orderRes, _ := orderClient.CreateOrder(context.Background(),
		&orderpb.CreateOrderRequest{
			UserId:       "user-001",
			RestaurantId: "R001",
		})
	result += fmt.Sprintf("ORDER: %v\n", orderRes)

	return result
}
