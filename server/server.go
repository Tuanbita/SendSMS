package main

import (
	"fmt"
	"log"
	"net"
	pb "example/SendSMS"
	"google.golang.org/grpc"

	"git.apache.org/thrift.git/lib/go/thrift"
	"example/TestThrift/gen-go/mythrift/demo"

	"strings"
	"os"
	"github.com/claudiu/gocron"
)
const (
	address = "0.0.0.0:10023"
)
const (
	HOST = "127.0.0.1"
	PORT = "9090"
)
type SendSMSService struct {
}

var clients = make(map[int]*Client)

type Client struct {
	key int
	stream pb.SendSMS_SendServer
}
func task() {
	fmt.Println("I am runnning task.")
}
func abc(){
	s := gocron.NewScheduler()
	s.Every(1).Day().Do(task)
	//s.Every(4).Seconds().Do(vijay)

	sc := s.Start() // keep the channel
	//go test(s, sc)  // wait
	<-sc            // it will happens if the channel is closed
}

var dem = 0
func (svr *SendSMSService) Send(stream pb.SendSMS_SendServer) error {
//luu cau truc thong tin cua 1 client
	client := &Client{
		key:  dem,
		stream: stream,
	}
	log.Println("[Send]: ", dem)
	clients[dem] = client
	dem ++
	abc()
	log.Println("[Send]: ket thuc ham Send")
	return nil
}
//nhan lenh tu thrift
func listenFromThrift() {
	for {
		transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
		protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()

		transport, err := thrift.NewTSocket(net.JoinHostPort(HOST, PORT))
		if err != nil {
			fmt.Fprintln(os.Stderr, "error resolving address:", err)
			os.Exit(1)
		}

		useTransport := transportFactory.GetTransport(transport)
		client := demo.NewMyThriftClientFactory(useTransport, protocolFactory)
		if err := transport.Open(); err != nil {
			fmt.Fprintln(os.Stderr, "Error opening socket to "+HOST+":"+PORT, " ", err)
			os.Exit(1)
		}

		defer transport.Close()
//gui tin hieu cho serverthrift va doi lenh
		sendmess, errr := client.SendSMS()

		if len(clients) > 0 && sendmess != "" && errr == nil {
			s := strings.Split(string(sendmess), " ")

			var mess pb.SendMessage
			mess.Content = s[1]
			mess.ToNumber = s[0]
			broascast(mess)
		}
	}
}
func broascast(mess pb.SendMessage){
	for i := 0; i<dem; i++ {
		clients[i].stream.Send(&mess)
	}
}
func main() {
	listen, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal(err)
	}
	s := grpc.NewServer()
	pb.RegisterSendSMSServer(s, &SendSMSService{})

	fmt.Println("Listening on the 0.0.0.0:10023")

	go listenFromThrift()

	if err := s.Serve(listen); err != nil {
		log.Fatal(err)
	}
}