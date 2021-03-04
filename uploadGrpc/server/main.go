package main

import (
	"fmt"
	"net"

	pb "bitbucket.org/xinapsedev/malmoy-file-server/pbs"
	"bitbucket.org/xinapsedev/malmoy-file-server/service"

	"google.golang.org/grpc"
)

func main() {
	listen, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Listen : ", err)
		return
	}

	svr := grpc.NewServer()
	fileStore := service.NewDiskFileStore()
	fileServer := service.NewFileServer(fileStore)
	pb.RegisterFileServiceServer(svr, fileServer)

	fmt.Println("server listen success : ", listen.Addr().String())

	if err := svr.Serve(listen); err != nil {
		fmt.Println("Serve : ", err)
		return
	}
}
