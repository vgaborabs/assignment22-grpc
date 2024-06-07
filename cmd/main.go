package main

import (
	"fmt"
	"github.com/vgaborabs/assignment22-grpc/internal/db"
	"github.com/vgaborabs/assignment22-grpc/internal/user"
	pb "github.com/vgaborabs/assignment22-grpc/proto"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

func main() {
	port := os.Getenv("SERVER_PORT")
	if len(port) == 0 {
		port = "5000"
	}
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("failed to listen on port %s: %v", err, port)
	}
	ur := db.NewInMemoryUserRepo()
	s := grpc.NewServer()
	srv := user.NewUserService(ur)
	pb.RegisterUserServiceServer(s, srv)
	log.Printf("gRPC server listening at %v", lis.Addr())
	log.Fatal(s.Serve(lis))
}
