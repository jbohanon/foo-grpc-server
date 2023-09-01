package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	pb "github.com/jbohanon/foo-grpc-server/api"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	port = flag.Int("port", 5000, "port to listen. Use 5000 to directly dial a server running from this repo's pkg/server/server.go")
)

func main() {
	flag.Parse()
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	addr := fmt.Sprintf("localhost:%d", *port)
	log.Printf("dialing %s\n", addr)
	conn, err := grpc.Dial(addr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := pb.NewFooClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err = client.GetFoo(ctx, &pb.FooRequest{})
	if err != nil {
		log.Fatalf("client.GetFoo failed: %v", err)
	}
	log.Println("client.GetFoo successful")
}
