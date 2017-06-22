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
	"google.golang.org/grpc/reflection"
	. "streamserver/fabric"
	"strconv"
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
	log.Printf("destory database...");
	s.db.UnRegisterDB()
}

func (s *server) RegisterClient(ctx context.Context, in *pb.RegisterReq) (*pb.RegisterReply, error) {
	isExist := true
	var err error
	if in.User != s.AuthUser && in.Pwd != s.AuthPwd {
		return &pb.RegisterReply{Message:"login user or password is abort", Info: nil}, nil
	}

	appid := GetAppId()
	appkey := GetAppKey()
	if appid == nil || appkey == nil {
		return &pb.RegisterReply{Message:"appid or appkey break abort", Info: nil}, nil
	}

	for {
		if appid != nil {
			isExist, err = s.db.IsExist(hex.EncodeToString(appid[:]))
		}
		if err != nil {
			return &pb.RegisterReply{Message:"query appid abort", Info: nil}, err
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
		return &pb.RegisterReply{Message:"insert appid and appkey abort", Info: nil}, err
	}
	err = s.myca.RegisterAndEnrollUser(hex.EncodeToString(appid[:]), hex.EncodeToString(appkey[:]))
	if err != nil {
		return &pb.RegisterReply{Message:"request ca ecert abort", Info: nil}, err
	}
	appInfo := &pb.RegisterInfo{Appid : hex.EncodeToString(appid[:]), Appkey: hex.EncodeToString(appkey[:])}
	return &pb.RegisterReply{Message : "OK", Info : appInfo},nil
}

func (s *server) InitAsset(ctx context.Context, in *pb.Asset) (*pb.MsgReply, error) {
	isExist, err := s.db.IsExist(in.Userid)
	if err != nil {
		return &pb.MsgReply{Message : "query dbase abort!"}, err
	}
	if isExist == false {
		return &pb.MsgReply{Message : "appid is not exist!"}, err
	}

	user, err := s.myca.LoadUser(in.Userid)
	if err != nil || user != nil{

	}

	err = s.fabric.Invoke(in.Userid)
	if err != nil {
		return &pb.MsgReply{Message : "Failt"}, err
	}

	return &pb.MsgReply{Message : "OK"}, nil
}

func (s *server) DealTransaction(ctx context.Context, in *pb.Transaction) (*pb.MsgReply, error){
	isExist1, err1 := s.db.IsExist(in.Ownerid)
	if err1 != nil {
		return &pb.MsgReply{Message : "query dbase abort!"}, err1
	}
	if isExist1 == false {
		return &pb.MsgReply{Message : "appid is not exist!"}, err1
	}

	isExist2, err2 := s.db.IsExist(in.Receiverid)
	if err2 != nil {
		return &pb.MsgReply{Message : "query dbase abort!"}, err2
	}
	if isExist2 == false {
		return &pb.MsgReply{Message : "appid is not exist!"}, err2
	}

	user1, err1 := s.myca.LoadUser(in.Ownerid)
	if err1 != nil || user1 != nil{

	}
	user2, err2 := s.myca.LoadUser(in.Receiverid)
	if err2 != nil || user2 != nil{

	}

	return nil,nil
}

func (s *server)  QueryAsset(ctx context.Context, in *pb.Asset) (*pb.Asset, error){
	isExist, err := s.db.IsExist(in.Userid)
	if err != nil {
		return &pb.Asset{Userid : "query abort", Value: 0}, err
	}
	if isExist == false {
		return &pb.Asset{Userid : "appid not exist", Value: 0}, err
	}
	user, err := s.myca.LoadUser(in.Userid)
	if err != nil || user != nil{

	}
	value, err := s.fabric.Query(in.Userid)
	if err != nil {
		return &pb.Asset{Userid : "query failt", Value: 0}, err
	}
	//log.Fatalf("user %s, value %s", in.Userid, value)
	ivalue,_ := strconv.ParseInt(value, 10, 32)
	return &pb.Asset{Userid : in.Userid, Value: int32(ivalue)}, nil
}

func main() {
	myserver.InitServer()
	listen, err := net.Listen("tcp", serveraddr)
	if err != nil {
		log.Fatalf("failed to listen %v\n",err)
	}

	s := grpc.NewServer()
	pb.RegisterStreamServerServer(s, &myserver)
	reflection.Register(s)
	if err := s.Serve(listen); err != nil {
		log.Fatalf("failed to server %v",err)
	}
	myserver.UnInitServer()
}






