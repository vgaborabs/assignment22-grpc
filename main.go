package main

import (
	"github.com/vgaborabs/assignment22-grpc/internal/db"
	"github.com/vgaborabs/assignment22-grpc/internal/user"
	pb "github.com/vgaborabs/assignment22-grpc/proto"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen on port 50051: %v", err)
	}
	ur := db.NewInMemoryUserRepo()
	s := grpc.NewServer()
	srv := user.NewUserService(ur)
	pb.RegisterUserServiceServer(s, srv)
	log.Printf("gRPC server listening at %v", lis.Addr())
	log.Fatal(s.Serve(lis))
}
