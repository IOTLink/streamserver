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
		log.Printf("appInfo: %s", appInfo)
	}

	//Asset
	/*
	asset := &pb.Asset{Userid: "2bf934667120d18ea30c4c55e59e56ed", Value:100}
	msgReply, err := client.InitAsset(context.Background(), asset)
	if err != nil {
		log.Fatalf("msg:%s,%v",msgReply.Message, err)
	}
	if msgReply != nil {
		msgAsset, _ := json.Marshal(msgReply)
		log.Printf("Asset: %s", msgAsset)
	}*/

	//DealTransaction
	tx := &pb.Transaction{Ownerid:"2bf934667120d18ea30c4c55e59e56ed", Receiverid:"cc84e1f40fd9d36c90061ea2ef371f40",Value : 10}
	msgRep, err := client.DealTransaction(context.Background(), tx)
	if err != nil {
		log.Fatalf("msg:%v", err)
	}
	if msg != nil {
		info, _ := json.Marshal(msgRep)
		log.Printf("Asset: %s", info)
	}


	//QueryAsset
	asset := &pb.Asset{Userid: "2bf934667120d18ea30c4c55e59e56ed", Value: 0 }
	msgasset, err := client.QueryAsset(context.Background(), asset)
	if err != nil {
		log.Fatalf("msg:%v", err)
	}
	if msgasset != nil {
		msgasset, _ := json.Marshal(msgasset)
		log.Printf("Asset: %s", msgasset)
	}
}
