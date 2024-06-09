package main

import (
	"database/sql"
	"fmt"
	"net"

	"github.com/celsopires1999/grpc_v2/internal/database"
	"github.com/celsopires1999/grpc_v2/internal/pb"
	"github.com/celsopires1999/grpc_v2/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "./db.sqlite")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	categoryDb := database.NewCategory(db)
	categoryService := service.NewCategoryService(*categoryDb)

	grpcServer := grpc.NewServer()
	pb.RegisterCategoryServiceServer(grpcServer, categoryService)

	// Register reflection service on gRPC server to work with evans
	reflection.Register(grpcServer)

	port := "50051"
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		panic(err)
	}

	fmt.Printf("Server running on port: %s\n", port)
	if err := grpcServer.Serve(lis); err != nil {
		panic(err)
	}
}
