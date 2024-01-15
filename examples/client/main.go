package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"myMiniBlog/internal/pkg/log"
	proto "myMiniBlog/pkg/proto/miniblog/v1"
	"time"
)

var (
	addr  = flag.String("addr", "localhost:9090", "The address to connect to.")
	limit = flag.Int64("limit", 10, "Limit to list users.")
)

func main() {
	flag.Parse()
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalw("Did not connect", "err", err)
	}
	defer conn.Close()
	c := proto.NewMiniBlogClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.ListUser(ctx, &proto.ListUserRequest{Offset: 0, Limit: *limit})
	if err != nil {
		log.Fatalw("could not greet: %v", err)
	}
	fmt.Println("TotalCount:", r.TotalCount)
	for _, u := range r.Users {
		d, _ := json.Marshal(u)
		fmt.Println(string(d))
	}
}
