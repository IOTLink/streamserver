package main

import (
	"log"
	"net"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	pb "streamserver/protocol"
	. "streamserver/models"
	. "streamserver/config"
	"encoding/hex"
	//"google.golang.org/grpc/reflection"
	. "streamserver/fabric"
	"google.golang.org/grpc/reflection"
)

const (
	serveraddr = "0.0.0.0:50055"
)

type server struct{
	AuthUser string
	AuthPwd  string
	db       ManageDB
	myca     CA
	fabric   FabricServer
}
var myserver server

func (s *server) InitServer() {
	s.AuthUser, s.AuthPwd = GetAuthAdmin()
	log.Printf("AuthUser : %s, AuthPwd : %s", s.AuthUser, s.AuthPwd)

	log.Printf("init database...")
	s.db.InitDB("","","",0, "")
	s.db.RegisterDB()
	if s.db.Database == nil {
		log.Printf("db.Database is nil")
	}

	s.myca.InitCaServer()
	s.fabric.Init()
}

func (s *server) UnInitServer() {
	log.Printf("destory database...")
	s.db.UnRegisterDB()
}

func (s *server) RegisterClient(ctx context.Context, in *pb.RegisterRequest) (*pb.RegisterReply, error) {
	isExist := true
	var err error
	if in.User != s.AuthUser && in.Pwd != s.AuthPwd {
		return &pb.RegisterReply{Message:"login user or password is incorrect", Appid:"", Appkey:""}, nil
	}

	if in.Chainid == "" {
		return &pb.RegisterReply{Message:"chainid can not be empty", Appid:"", Appkey:""}, nil
	}

	appid := GetAppId()
	appkey := GetAppKey()
	if appid == nil || appkey == nil {
		return &pb.RegisterReply{Message:"generate appid or appkey error",  Appid:"", Appkey:""}, nil
	}

	for {
		if appid != nil {
			isExist, err = s.db.IsExist(hex.EncodeToString(appid[:]))
		}
		if err != nil {
			return &pb.RegisterReply{Message:"query appid exception"}, err
		}
		if isExist == true {
			appid = GetAppId()
			continue
		} else {
			break
		}
	}

	dayTimes := GetUTCTimeStr()
	_, err = s.db.InsertAppInfo(hex.EncodeToString(appid[:]), hex.EncodeToString(appkey[:]), dayTimes)
	if err != nil {
		return &pb.RegisterReply{Message:"insert appid and appkey failed", Appid:"", Appkey:""}, err
	}
	err = s.myca.RegisterAndEnrollUser(hex.EncodeToString(appid[:]), hex.EncodeToString(appkey[:]))
	if err != nil {
		return &pb.RegisterReply{Message:"registering user to ca-server failed", Appid:"", Appkey:""}, err
	}
	return &pb.RegisterReply{Message : "OK",  Appid:hex.EncodeToString(appid[:]), Appkey:hex.EncodeToString(appkey[:])}, nil
}

func (s *server) EnrollAsset(ctx context.Context, in *pb.AssetEnroll) (*pb.ResultsReply, error) {
	if in.Chainid == "" || in.Chaincodeid == "" {
		return &pb.ResultsReply{Message:"chainid or chaincodeid can not be empty", Payload:""}, nil
	}
	if in.Appid == "" {
		return &pb.ResultsReply{Message:"appid can not be empty", Payload:""}, nil
	}
	if in.Payload == "" {
		return &pb.ResultsReply{Message:"payload can not be empty", Payload:""}, nil
	}

	isExist, err := s.db.IsExist(in.Appid)
	if err != nil {
		return &pb.ResultsReply{Message : "query database exception", Payload:""}, err
	}
	if isExist == false {
		return &pb.ResultsReply{Message : "appid does not exist", Payload:""}, err
	}

	user, err := s.myca.LoadUser(in.Appid)
	if err != nil || user != nil{

	}

	err = s.fabric.InitAsset(in.Appid, in.Payload)
	if err != nil {
		return &pb.ResultsReply{Message : "execution func enrollasset failed", Payload:""}, err
	}

	return &pb.ResultsReply{Message : "OK", Payload:""}, nil
}

func (s *server) RegisterAsset(ctx context.Context, in *pb.AssetRegister) (*pb.ResultsReply, error) {
	if in.Chainid == "" || in.Chaincodeid == "" {
		return &pb.ResultsReply{Message:"chainid or chaincodeid can not be empty", Payload:""}, nil
	}
	if in.Appid == "" {
		return &pb.ResultsReply{Message:"appid can not be empty", Payload:""}, nil
	}
	if in.Payload == "" {
		return &pb.ResultsReply{Message:"payload can not be empty", Payload:""}, nil
	}

	isExist, err := s.db.IsExist(in.Appid)
	if err != nil {
		return &pb.ResultsReply{Message : "query database exception", Payload:""}, err
	}
	if isExist == false {
		return &pb.ResultsReply{Message : "appid does not exist", Payload:""}, err
	}

	user, err := s.myca.LoadUser(in.Appid)
	if err != nil || user != nil{

	}

	txID, err := s.fabric.Init2Asset(in.Appid, in.Payload)
	if err != nil {
		return &pb.ResultsReply{Message : "execution func enrollasset failed", Payload:""}, err
	}
	return &pb.ResultsReply{Message : "OK", Payload : txID}, nil
}


func (s *server) TransactionAsset(ctx context.Context, in *pb.TransactionRequest) (*pb.ResultsReply, error){
	if in.Chainid == "" || in.Chaincodeid == "" {
		return &pb.ResultsReply{Message:"chainid or chaincodeid can not be empty", Payload:""}, nil
	}
	if in.Appidower == "" || in.Appidreceive == ""{
		return &pb.ResultsReply{Message:"appidower or appidreceive can not be empty", Payload:""}, nil
	}
	if in.Payload == "" {
		return &pb.ResultsReply{Message:"payload can not be empty", Payload:""}, nil
	}

	isExist1, err1 := s.db.IsExist(in.Appidower)
	if err1 != nil {
		return &pb.ResultsReply{Message : "query database exception", Payload:""}, err1
	}
	if isExist1 == false {
		return &pb.ResultsReply{Message : "appid does not exist", Payload:""}, err1
	}

	isExist2, err2 := s.db.IsExist(in.Appidreceive)
	if err2 != nil {
		return &pb.ResultsReply{Message : "query database exception", Payload:""}, err2
	}
	if isExist2 == false {
		return &pb.ResultsReply{Message : "appid does not exist", Payload:""}, err2
	}

	user1, err1 := s.myca.LoadUser(in.Appidower)
	if err1 != nil || user1 != nil{

	}
	user2, err2 := s.myca.LoadUser(in.Appidreceive)
	if err2 != nil || user2 != nil{

	}
	txId, err := s.fabric.Transfer(in.Appidower, in.Appidreceive, in.Payload)
	if err != nil {
		return &pb.ResultsReply{Message : "transaction execution failed", Payload:""}, err
	}

	return &pb.ResultsReply{Message:"OK" , Payload:txId}, nil
}

func (s *server)  QueryAsset(ctx context.Context, in *pb.QueryRequest) (*pb.ResultsReply, error){
	if in.Chainid == "" || in.Chaincodeid == "" {
		return &pb.ResultsReply{Message:"chainid or chaincodeid can not be empty", Payload:""}, nil
	}

	isExist, err := s.db.IsExist(in.Appid)
	if err != nil {
		return &pb.ResultsReply{Message : "query database exception", Payload:""}, err
	}
	if isExist == false {
		return &pb.ResultsReply{Message : "appid does not exist", Payload:""}, err
	}

	user, err := s.myca.LoadUser(in.Appid)
	if err != nil || user != nil{

	}
	value, err := s.fabric.Query(in.Appid)
	if err != nil {
		return &pb.ResultsReply{Message : "query appid exception", Payload:""}, err
	}

	return &pb.ResultsReply{Message : value}, nil
}

func main() {
	myserver.InitServer()
	listen, err := net.Listen("tcp", serveraddr)
	if err != nil {
		log.Fatalf("server exec listen failed %v\n",err)
	}

	s := grpc.NewServer()
	pb.RegisterStreamServerServer(s, &myserver)

	reflection.Register(s)
	if err := s.Serve(listen); err != nil {
		log.Fatalf("start service failed %v",err)
	}

	myserver.UnInitServer()
}






