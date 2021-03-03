package main

import (
	"fmt"
	"net"
	"sampleProject/grpcProject/pb"
	"sampleProject/grpcProject/service"

	"google.golang.org/grpc"
)

func main() {
	listen, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Listen : ", err)
		return
	}

	svr := grpc.NewServer()
	fileStore := service.NewDiskFileStore("", "")
	fileServer := service.NewFileServer(fileStore)
	pb.RegisterFileServiceServer(svr, fileServer)

	fmt.Println("server listen success : ", listen.Addr().String())

	if err := svr.Serve(listen); err != nil {
		fmt.Println("Serve : ", err)
		return
	}
}
