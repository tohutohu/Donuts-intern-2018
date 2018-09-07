package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"strconv"
	"time"

	"github.com/tohutohu/Donuts/day4/practice5/proto"
	"google.golang.org/grpc"
)

func main() {
	var method string
	flag.StringVar(&method, "m", "fizzbuzz", "method")
	flag.Parse()

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := proto.NewFizzBuzzClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	switch method {
	case "fizzbuzz":
		fizzBuzzInput, err := strconv.ParseInt(flag.Arg(0), 10, 64)
		if err != nil {
			panic(err)
		}
		r, err := c.CalcFizzBuzz(ctx, &proto.CalcFizzBuzzRequest{Num: fizzBuzzInput})
		if err != nil {
			log.Fatalf("could not fizzbuzz: %v", err)
		}
		log.Printf("%v", r.Res)

	case "square":
		numList := []int64{}
		for _, v := range flag.Args() {
			num, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				panic(err)
			}
			numList = append(numList, num)
		}
		r, err := c.SquareList(ctx, &proto.SquareListRequest{NumList: numList})
		if err != nil {
			log.Fatalf("could not square: %v", err)
		}
		log.Printf("%v", r.NumList)

	case "stream":
		waitc := make(chan struct{})
		stream, err := c.Stream(context.Background())
		if err != nil {
			panic(err)
		}
		go func() {
			for {
				in, err := stream.Recv()
				if err == io.EOF {
					waitc <- struct{}{}
					return
				}
				if err != nil {
					fmt.Println("Failed recieve: ", err)
				}
				fmt.Println(in)
			}
		}()
		i := 0
		for {
			i++
			time.Sleep(time.Second)
			if err := stream.Send(&proto.StreamStruct{Num: int64(i)}); err != nil {
				fmt.Println("Failed send: ", err)
			}
			if i > 100 {
				break
			}
		}
		stream.CloseSend()
		<-waitc

	default:
		fizzBuzzInput, err := strconv.ParseInt(flag.Arg(0), 10, 64)
		if err != nil {
			panic(err)
		}
		r, err := c.CalcFizzBuzz(ctx, &proto.CalcFizzBuzzRequest{Num: fizzBuzzInput})
		if err != nil {
			log.Fatalf("could not fizzbuzz: %v", err)
		}
		log.Printf("%v", r.Res)

	}
}
