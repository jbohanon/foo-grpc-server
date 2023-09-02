package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"

	pb "github.com/jbohanon/foo-grpc-server/api"

	"google.golang.org/grpc"
)

var (
	port       = flag.Int("port", 5000, "gRPC listen port")
	healthport = flag.Int("healthport", 5001, "health check listen port")
)

type fooServer struct {
	pb.UnimplementedFooServer
}

func (s *fooServer) GetFoo(ctx context.Context, req *pb.FooRequest) (*pb.FooResponse, error) {
	log.Println("received request to GetFoo")
	return &pb.FooResponse{}, nil
}

func serveHealth(cancel context.CancelFunc) {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	http.HandleFunc("/healthbad", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	http.HandleFunc("/quitquitquit", func(w http.ResponseWriter, r *http.Request) {
		log.Println("quitquitquit request received. shutting down.")
		cancel()
	})
	addr := fmt.Sprintf("localhost:%d", *healthport)
	log.Printf("health server listening on %s", addr)
	http.ListenAndServe(addr, nil)
}

func serveGrpc(cancel context.CancelFunc) {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterFooServer(grpcServer, &fooServer{})
	log.Println("gRPC FooServer listening on localhost:5000")

	if err = grpcServer.Serve(lis); err != nil {
		cancel()
	}
}

func main() {
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())

	go serveHealth(cancel)

	go serveGrpc(cancel)

	<-ctx.Done()

}
