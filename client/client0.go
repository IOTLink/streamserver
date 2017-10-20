package main

import (
	"log"
	"google.golang.org/grpc"
	pb "streamserver/protocol"
	"golang.org/x/net/context"
	"encoding/json"
	//"github.com/hyperledger/fabric/bccsp"
	//"github.com/hyperledger/fabric-sdk-go/api/apitxn"
	//"time"
	//"time"
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
	regInfo := &pb.RegisterRequest{User: "admin", Pwd: "adminw", Org: "peerorg1", Affiliation: "org1.department1"}

	for {
		userx, err := client.RegisterClient(context.Background(), regInfo)
		if err != nil {
			log.Fatalf("msg:%v", err)
		}
		if userx != nil {
			appInfo, _ := json.Marshal(userx)
			log.Printf("register: %s", appInfo)
		}
	}
	return

	user1, err := client.RegisterClient(context.Background(), regInfo)
	if err != nil {
		log.Fatalf("msg:%v", err)
	}
	if user1 != nil {
		appInfo, _ := json.Marshal(user1)
		log.Printf("register: %s", appInfo)
	}

	user2, err := client.RegisterClient(context.Background(), regInfo)
	if err != nil {
		log.Fatalf("msg:%s,%v", user2.Message, err)
	}
	if user2 != nil {
		appInfo, _ := json.Marshal(user2)
		log.Printf("register: %s", appInfo)
	}
	//time.Sleep(time.Second *1 )


	mychannel := "mychannel" //businesschannel
	chaincodeName := "mychannel_cc_v2"

	/*
	//install chaincode mychannel businesschannel
	assetInfo := &pb.AssetEnroll{Channel: mychannel, Chaincodepath:"github.com/chaincode", Chaincode:chaincodeName,
		Chaincodeversion:"v2", Key:"initkey", Payload:"100"}
	msgReply, err := client.EnrollAsset(context.Background(), assetInfo)
	if err != nil {
		log.Println("msg:%s","", err.Error())
	}
	if msgReply != nil {
		info, _ := json.Marshal(msgReply)
		log.Printf("EnrollAsset: %s", info)
	}
	*/


	//fun invoke init user1

	assetRegInfo := &pb.AssetRegister{Channel: mychannel, Chaincode: chaincodeName,
		Appid:                                 user1.Appid, Key: user1.Appid, Payload: "100"}
	msgReply, err := client.RegisterAsset(context.Background(), assetRegInfo)
	if err != nil {
		log.Println("msg:", "", err.Error())
	}
	if msgReply != nil {
		info, _ := json.Marshal(msgReply)
		log.Printf("AssetRegister: %s", info)
	}

	//fun invoke init user2
	assetRegInfo = &pb.AssetRegister{Channel: mychannel, Chaincode: chaincodeName,
		Appid:                                 user2.Appid, Key: user2.Appid, Payload: "100"}
	msgReply, err = client.RegisterAsset(context.Background(), assetRegInfo)
	if err != nil {
		log.Println("msg:", "", err.Error())
	}
	if msgReply != nil {
		info, _ := json.Marshal(msgReply)
		log.Printf("AssetRegister: %s", info)
	}

	transInfo := &pb.Transaction{Channel: mychannel, Chaincode: chaincodeName,
		Ownerid:                          user1.Appid, Receiverid: user2.Appid, Payload: "50"}
	msgReply, err = client.TransactionAsset(context.Background(), transInfo)
	if err != nil {
		log.Println("msg:", "", err.Error())
	}
	if msgReply != nil {
		info, _ := json.Marshal(msgReply)
		log.Printf("TransactionAsset: %s", info)
	}

	//invoke query user1
	assetQueryInfo := &pb.AssetQuery{Channel: mychannel, Chaincode: chaincodeName,
		Appid: user1.Appid, Key: user1.Appid}
	msgReply, err = client.QueryAsset(context.Background(), assetQueryInfo)
	if err != nil {
		log.Println("msg:", "", err.Error())
	}
	if msgReply != nil {
		info, _ := json.Marshal(msgReply)
		log.Printf("AssetQuery: %s", info)
	}

	//invoke query user1
	assetQueryInfo = &pb.AssetQuery{Channel: mychannel, Chaincode: chaincodeName,
		Appid:                               user2.Appid, Key: user1.Appid}
	msgReply, err = client.QueryAsset(context.Background(), assetQueryInfo)
	if err != nil {
		log.Println("msg:", "", err.Error())
	}
	if msgReply != nil {
		info, _ := json.Marshal(msgReply)
		log.Printf("AssetQuery: %s", info)
	}

}
