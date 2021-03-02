package main

import (
	"context"
	"fmt"
	"io"

	pb "sampleProject/grpcProject/pb/proto"

	"google.golang.org/grpc"
)

// sample grpc stream server-client example
// https://www.freecodecamp.org/news/grpc-server-side-streaming-with-go/
// https://github.com/pahanini/go-grpc-bidirectional-streaming-example

func main() {
	fmt.Println("start client")

	conn, err := grpc.Dial(":8888", grpc.WithInsecure())
	if err != nil {
		fmt.Println("Dial : ", err)
		return
	}

	client := pb.NewStreamServiceClient(conn)
	in := &pb.Request{Id: 1}
	stream, err := client.FetchResponse(context.Background(), in)
	if err != nil {
		fmt.Println("open stream error : ", err.Error())
		return
	}

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			fmt.Println("Recv : ", err.Error())
			break
		}

		if err != nil {
			fmt.Println("Resp received : ", err)
			break
		}

		fmt.Println("resp : ", resp)
	}
}
