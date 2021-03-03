package service

import (
	"bytes"
	"fmt"
	"io"
	"sampleProject/grpcProject/pb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// FileServer is the server that provide file services.
// If want to service for FileServer, make *.go and add interface.
type FileServer struct {
	fileStore                         FileStore
	pb.UnimplementedFileServiceServer // not use, erase error.
}

// NewFileServer returns a new FileServer
func NewFileServer(fileStore FileStore) *FileServer {
	return &FileServer{fileStore, pb.UnimplementedFileServiceServer{}}
}

// UploadFile is a client-streaming RPC to upload a file
func (server *FileServer) UploadFile(stream pb.FileService_UploadFileServer) error {
	fmt.Println("UploadFile start")

	// receive file metadata
	req, err := stream.Recv()
	if err != nil {
		fmt.Println("Recv : ", err)
		return status.Errorf(codes.Unknown, "cannot receive image info : %v", err)
	}

	speakerID := req.GetMetadata().GetSpeakerId()
	modelID := req.GetMetadata().GetModelId()
	fmt.Printf("speakerID : %s, modelID : %s\n", speakerID, modelID)

	// receive chunk data
	fileData := bytes.Buffer{}
	fileSize := 0

	for {
		if err := stream.Context(); err.Err() != nil {
			fmt.Println("stream.Context : ", err)
			return err.Err()
		}

		fmt.Println("waiting to receive more data")

		req, err := stream.Recv()
		if err == io.EOF {
			fmt.Println("no more data")
			break
		}
		if err != nil {
			fmt.Println("stream.Recv : ", err)
			return err
		}

		chunk := req.GetChunkData()
		size := len(chunk)

		fmt.Println("received a chunk with size : ", size)

		fileSize += size

		_, err = fileData.Write(chunk)
		if err != nil {
			fmt.Println("fileData.Write : ", err)
			return err
		}
	}

	// save chunk data
	if err := server.fileStore.Save(speakerID, modelID, "test.txt", fileData); err != nil {
		return err
	}

	// response
	res := &pb.UploadFileResponse{
		Status:    0,
		SpeakerId: speakerID,
		ModelId:   modelID,
	}

	if err := stream.SendAndClose(res); err != nil {
		fmt.Println("SendAndClose : ", err)
		return status.Errorf(codes.Unknown, "cannot send response : %v", err)
	}

	fmt.Println("UploadFile end")
	return nil
}
