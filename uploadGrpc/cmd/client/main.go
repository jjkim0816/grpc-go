package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"sampleProject/grpcProject/pb"

	"google.golang.org/grpc"
)

const (
	// ChunkSize define buffer size when send to server
	ChunkSize = 2048
)

func main() {
	transportOption := grpc.WithInsecure()
	client, err := grpc.Dial(":8080", transportOption)
	if err != nil {
		fmt.Println("Dial : ", err.Error())
		return
	}

	// NewFileServiceClient create a new FileServiceClient interface
	service := pb.NewFileServiceClient(client)

	// UploadFile create a new FileService_UploadFileClient interface
	stream, err := service.UploadFile(context.Background())
	if err != nil {
		fmt.Println("UploadFile")
		return
	}

	// create Metadata
	req := &pb.UploadFileRequest{
		Data: &pb.UploadFileRequest_Metadata{
			Metadata: &pb.FileInfo{
				SpeakerId: "A2C23DDC66",
				ModelId:   "C2AE5FCCF0",
			},
		},
	}

	// Send which is FileService_UploadFileClient interface send first message of oneof to server
	if err := stream.Send(req); err != nil {
		fmt.Println("Send : ", err)
		return
	}

	// 파일 압축

	// 파일 열기
	filePath := "/Users/jjkim/workspace/src/sampleProject/grpcProject/files/201119_0000.txt"
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Open : ", err)
		return
	}

	// 파일 읽기
	reader := bufio.NewReader(file)
	buffer := make([]byte, ChunkSize)
	size := 0

	// chunk_data 보내기
	for {
		read, err := reader.Read(buffer)
		if err == io.EOF {
			fmt.Println("file is end")
			break
		}

		size += read

		req := &pb.UploadFileRequest{
			Data: &pb.UploadFileRequest_ChunkData{
				ChunkData: buffer[:read],
			},
		}

		if err = stream.Send(req); err != nil {
			fmt.Println("Send : ", err)
			break
		}
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		fmt.Println("CloseAndRecv : ", err)
		return
	}

	fmt.Printf("res.Status : %d, Speaker_id : %s, Model_id : %s\n",
		res.GetStatus(), res.GetSpeakerId(), res.GetModelId())
}
