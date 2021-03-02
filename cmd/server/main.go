package main

import (
	"fmt"
	"log"
	"net"
	"time"

	pb "sampleProject/grpcProject/pb/proto"

	"google.golang.org/grpc"
)

// sample grpc stream server-client example
// https://www.freecodecamp.org/news/grpc-server-side-streaming-with-go/
// https://github.com/pahanini/go-grpc-bidirectional-streaming-example

type server struct {
	pb.UnsafeStreamServiceServer
}

func (s *server) FetchResponse(in *pb.Request, svr pb.StreamService_FetchResponseServer) error {

	for {
		resp := pb.Response{Result: fmt.Sprintf("success request")}
		if err := svr.Send(&resp); err != nil {
			fmt.Println("Send : ", err.Error())
			break
		}

		time.Sleep(time.Second * 5)
	}

	return nil
}

func main() {
	listen, err := net.Listen("tcp", ":8888")
	if err != nil {
		fmt.Printf("listen : %s\n", err.Error())
		return
	}

	fServer := grpc.NewServer()
	pb.RegisterStreamServiceServer(fServer, &server{})

	log.Println("start server")
	if err := fServer.Serve(listen); err != nil {
		fmt.Printf("Serve : %s\n", err.Error())
		return
	}
}
