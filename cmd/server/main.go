package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	pb "github.com/iamlilze/gRPC/api/auth"
	"github.com/iamlilze/gRPC/internal/service"
	"github.com/iamlilze/gRPC/internal/storage"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	defaultPort = "50051"
)

func main() {
	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	// Create storage
	store := storage.NewInMemoryStorage()

	// Create auth service
	authService := service.NewAuthService(store)

	// Create gRPC server
	grpcServer := grpc.NewServer()

	// Register service
	pb.RegisterAuthServiceServer(grpcServer, authService)

	// Register reflection service for debugging with grpcurl
	reflection.Register(grpcServer)

	// Start listening
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Handle graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		<-sigChan
		log.Println("Shutting down gRPC server...")
		grpcServer.GracefulStop()
	}()

	log.Printf("gRPC server listening on port %s", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
