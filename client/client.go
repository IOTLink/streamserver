package main

import (
	"log"
	"google.golang.org/grpc"
	pb "streamserver/protocol"
	"golang.org/x/net/context"
	"encoding/json"
)

const (
	serveraddr = "127.0.0.1:50055"
)

func main() {
	conn, err := grpc.Dial(serveraddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewStreamServerClient(conn)

	//RegisterReq
	regInfo := &pb.RegisterReq{User:"client", Pwd:"123456"}
	msg, err := client.RegisterClient(context.Background(), regInfo)
	if err != nil {
		log.Fatalf("msg:%s,%v",msg.Message, err)
	}
	if msg != nil {
		appInfo, _ := json.Marshal(msg)
		log.Printf("register: %s", appInfo)
	}

	//Asset
	/*
	asset := &pb.Asset{Userid: "b5a2f040374c79a11aa27d5020cc9a6f", Value:100}
	msgReply, err := client.InitAsset(context.Background(), asset)
	if err != nil {
		log.Fatalf("msg:%s,%v",msgReply.Message, err)
	}
	if msgReply != nil {
		msgAsset, _ := json.Marshal(msgReply)
		log.Printf("Asset: %s", msgAsset)
	}
	*/

	//DealTransaction

	tx := &pb.Transaction{Ownerid:"b5a2f040374c79a11aa27d5020cc9a6f", Receiverid:"6ea981c23737f0624c16ed7de2d2b8d0",Value : 10}
	msgRep, err := client.DealTransaction(context.Background(), tx)
	if err != nil {
		log.Fatalf("msg:%v", err)
	}
	if msg != nil {
		info, _ := json.Marshal(msgRep)
		log.Printf("move txid: %s", info)
	}

	asset := &pb.Asset{Userid: "b5a2f040374c79a11aa27d5020cc9a6f", Value: 0 }
	msgasset, err := client.QueryAsset(context.Background(), asset)
	if err != nil {
		log.Fatalf("msg:%v", err)
	}
	if msgasset != nil {
		msgasset, _ := json.Marshal(msgasset)
		log.Printf("query1: %s", msgasset)
	}
	//QueryAsset

	asset = &pb.Asset{Userid: "6ea981c23737f0624c16ed7de2d2b8d0", Value: 0 }
	msgasset, err = client.QueryAsset(context.Background(), asset)
	if err != nil {
		log.Fatalf("msg:%v", err)
	}
	if msgasset != nil {
		msgasset, _ := json.Marshal(msgasset)
		log.Printf("query2: %s", msgasset)
	}

}
