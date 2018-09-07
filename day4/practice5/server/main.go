package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"

	"github.com/tohutohu/Donuts/day4/practice5/proto"

	"google.golang.org/grpc"
)

type server struct{}

func (s *server) CalcFizzBuzz(c context.Context, p *proto.CalcFizzBuzzRequest) (*proto.CalcFizzBuzzReply, error) {
	var response string
	num := p.Num
	if num%15 == 0 {
		response = "fizzbuzz"
	} else if num%5 == 0 {
		response = "buzz"
	} else if num%3 == 0 {
		response = "fizz"
	} else {
		response = strconv.FormatInt(num, 10)
	}

	return &proto.CalcFizzBuzzReply{Res: response}, nil
}

func (s *server) SquareList(c context.Context, p *proto.SquareListRequest) (*proto.SquareListReply, error) {
	numList := p.NumList
	for i, num := range numList {
		numList[i] = num * num
	}

	return &proto.SquareListReply{NumList: numList}, nil
}

func (s *server) Stream(stream proto.FizzBuzz_StreamServer) error {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			fmt.Println("Failed recieve: ", err)
			return err
		}
		fmt.Println("Recieve: ", in.Num)
		stream.Send(&proto.StreamStruct{Num: in.Num * in.Num})
	}
}

func main() {
	port := "50051"
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failt to listen: %v", err)
	}
	s := grpc.NewServer()
	proto.RegisterFizzBuzzServer(s, &server{})
	go func() {
		log.Printf("start grpc server port :%v", port)
		s.Serve(lis)
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Printf("stopping grpc server...")
	s.GracefulStop()
}
