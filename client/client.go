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
	asset := &pb.Asset{Userid: "9fb05ce2c57905b41be708425bdda6c8", Value:100}
	msgReply, err := client.InitAsset(context.Background(), asset)
	if err != nil {
		log.Fatalf("msg:%s,%v",msgReply.Message, err)
	}
	if msgReply != nil {
		msgAsset, _ := json.Marshal(msgReply)
		log.Printf("Asset: %s", msgAsset)
	}
	*/


	//QueryAsset
	asset := &pb.Asset{Userid: "9fb05ce2c57905b41be708425bdda6c8", Value:0}
	msgasset, err := client.QueryAsset(context.Background(), asset)
	if err != nil {
		log.Fatalf("msg:%v", err)
	}
	if msgasset != nil {
		msgasset, _ := json.Marshal(msgasset)
		log.Printf("Asset: %s", msgasset)
	}

}
