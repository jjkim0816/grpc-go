package service

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"

	"bitbucket.org/xinapsedev/malmoy-file-server/pbs"
	pb "bitbucket.org/xinapsedev/malmoy-file-server/pbs"
	"bitbucket.org/xinapsedev/malmoy-file-server/util"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// FileServer is the server that provide file services.
// If want to service for FileServer, make *.go and add interface.
type FileServer struct {
	fileStore                         FileStore
	pb.UnimplementedFileServiceServer // not use, remove error.
}

// NewFileServer returns a new FileServer
func NewFileServer(fileStore FileStore) *FileServer {
	return &FileServer{fileStore, pb.UnimplementedFileServiceServer{}}
}

// HeartBeat check server lives
func (server *FileServer) HeartBeat(context.Context, *emptypb.Empty) (*pbs.HeartBeatResponse, error) {
	return &pbs.HeartBeatResponse{Beat: ""}, nil
}

// UploadFile is a client-streaming RPC to upload a file
func (server *FileServer) UploadFile(stream pb.FileService_UploadFileServer) error {
	fmt.Println("UploadFile start")

	// receive file metadata
	req, err := stream.Recv()
	if err != nil {
		fmt.Println("Recv : ", err)
		return status.Errorf(codes.Unknown, "cannot receive file info : %v", err)
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

		// fmt.Println("waiting to receive more data")

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

		// fmt.Println("received a chunk with size : ", size)

		fileSize += size

		_, err = fileData.Write(chunk)
		if err != nil {
			fmt.Println("fileData.Write : ", err)
			return err
		}
	}

	// save chunk data
	fileName := fmt.Sprintf("tmp_%s_%s.tar", speakerID, modelID)
	if err := server.fileStore.Save(speakerID, modelID, fileName, fileData); err != nil {
		return err
	}

	// 압축파일 풀기 및 이전 압축 파일 삭제
	target := fmt.Sprintf("/Users/jjkim/workspace/src/sampleProject/grpcProject/files/%s/%s/%s", speakerID, modelID, fileName)
	output := fmt.Sprintf("/Users/jjkim/workspace/src/sampleProject/grpcProject/files/%s/%s", speakerID, modelID)
	if err := util.DecompressTarFile(target, output); err == nil {
		os.Remove(target)
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
