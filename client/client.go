package main

import (
	"fmt"
	"log"
	pb "example/SendSMS"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)
const (
	address = "0.0.0.0:10023"
)
func chat(client pb.SendSMSClient) {
	stream, err := client.Send(context.Background())
	if err != nil {
		log.Fatal("client Chat get stream error: ", err)
	}

	for i:=1; ; i++{
		recv,_ := stream.Recv()
		fmt.Println("toNumber: ",recv.ToNumber,"  content: ",recv.Content)
	}
}
func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatal("dail error:", err)
	}
	defer conn.Close()
	c := pb.NewSendSMSClient(conn)

	chat(c)
}