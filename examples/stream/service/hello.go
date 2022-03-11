package service

import (
	"fmt"
	"io"
	"time"

	pb "github.com/SeeMusic/kratos/examples/stream/hello"
)

type HelloService struct {
	pb.UnimplementedHelloServer
}

func NewHelloService() *HelloService {
	return &HelloService{}
}

func (s *HelloService) GetNumber(req *pb.GetNumberRequest, conn pb.Hello_GetNumberServer) error {
	var number int64
	for {
		fmt.Println(req.Data)
		err := conn.Send(&pb.GetNumberReply{Number: number})
		if err != nil {
			return err
		}
		number++
		time.Sleep(time.Second)
	}
}

func (s *HelloService) UploadLog(conn pb.Hello_UploadLogServer) error {
	for {
		req, err := conn.Recv()
		if err == io.EOF {
			return conn.SendAndClose(&pb.UploadLogReply{Res: "ok"})
		}
		if err != nil {
			return err
		}
		fmt.Println(req.Log)
	}
}

func (s *HelloService) Chat(conn pb.Hello_ChatServer) error {
	for {
		req, err := conn.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		err = conn.Send(&pb.ChatReply{DownMsg: "hello " + req.UpMsg})
		if err != nil {
			return err
		}
	}
}
