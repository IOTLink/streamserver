package main

import (
	"log"
	"google.golang.org/grpc"
	pb "streamserver/protocol"
	"golang.org/x/net/context"
	"encoding/json"
	"time"
	"fmt"
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

	regInfo := &pb.RegisterRequest{User:"client", Pwd:"123456", Chainid:"mychannel"}
	msgReply1, err := client.RegisterClient(context.Background(), regInfo)
	if err != nil {
		log.Fatalf("msg:%s,%v",msgReply1.Message, err)
	}
	if msgReply1 != nil {
		info, _ := json.Marshal(msgReply1)
		log.Printf("RegisterClient: %s", info)
	}

	//Asset
	/*
	assetInfo := &pb.AssetEnroll{Chainid: "mychannel", Chaincodeid:"mychaincodev5",
		Appid:"2d89b1ba0f869c49f4dbcee15c766c5d", Payload:"100"}
	msgReply2, err := client.EnrollAsset(context.Background(), assetInfo)
	if err != nil {
		log.Fatalf("msg:%s,%v",msgReply2.Message, err)
	}
	if msgReply2 != nil {
		info, _ := json.Marshal(msgReply2)
		log.Printf("EnrollAsset: %s", info)
	}*/


	//return
	//DealTransaction
	start := time.Now()
	transactionInfo := &pb.TransactionRequest{Chainid: "mychannel", Chaincodeid:"mychaincodev5",
		Appidower:"2d89b1ba0f869c49f4dbcee15c766c5d",
		Appidreceive:"5e9761ccc2f59df2dc4f42173486c627",
		Payload : "10"}
	msgReply3,err := client.TransactionAsset(context.Background(), transactionInfo)
	if err != nil {
		log.Fatalf("msg:%v", err)
	}
	if msgReply3 != nil {
		info, _ := json.Marshal(msgReply3)
		log.Printf("TransactionAsset txid: %s", info)
	}
	fmt.Printf("SendTransaction %.2fs elapsed\n", time.Since(start).Seconds())

	start = time.Now()
	queryInfo1 := &pb.QueryRequest{Chainid: "mychannel", Chaincodeid:"mychaincodev5",
		Appid:"2d89b1ba0f869c49f4dbcee15c766c5d"}
	msgReply4, err := client.QueryAsset(context.Background(), queryInfo1)
	if err != nil {
		log.Fatalf("msg:%v", err)
	}
	if msgReply4 != nil {
		info, _ := json.Marshal(msgReply4)
		log.Printf("QueryAsset: %s", info)
	}
	fmt.Printf("QueryAsset %.2fs elapsed\n", time.Since(start).Seconds())
	//QueryAsset

	start = time.Now()
	queryInfo2 := &pb.QueryRequest{Chainid: "mychannel", Chaincodeid:"mychaincodev5",
		Appid:"5e9761ccc2f59df2dc4f42173486c627"}
	msgReply5, err := client.QueryAsset(context.Background(), queryInfo2)
	if err != nil {
		log.Fatalf("msg:%v", err)
	}
	if msgReply5 != nil {
		info, _ := json.Marshal(msgReply5)
		log.Printf("QueryAsset: %s", info)
	}
	fmt.Printf("QueryAsset %.2fs elapsed\n", time.Since(start).Seconds())

	/*
	start = time.Now()
	assetInfo1 := &pb.AssetEnroll{Chainid: "mychannel", Chaincodeid:"mychaincodev5", Appid:"dadsd",Payload:"100"}
	msgReply6, err := client.EnrollAsset(context.Background(), assetInfo1)
	if err != nil {
		log.Fatalf("msg:%v", err)
	}
	if msgReply6 != nil {
		info, _ := json.Marshal(msgReply6)
		log.Printf("query2: %s", info)
	}
	fmt.Printf("QueryAsset %.2fs elapsed\n", time.Since(start).Seconds())
	*/
}
