package main

import (
	"log"
	"net"
	"errors"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	pb "streamserver/protocol"
	. "streamserver/common"
	. "streamserver/fabric"
	. "streamserver/dbase"
	common "streamserver/common"
	. "streamserver/models"
	"google.golang.org/grpc/reflection"
	"encoding/hex"
	"fmt"
	"syscall"
	"os"
	"net/http"
    _ "net/http/pprof"
	"runtime"
	"time"
)

type server struct{
	AuthUser string
	AuthPwd  string
	IdServer IDServer //Generate a snowflake ID
	db       ManageDB
	myca     map[string]*CA
	fabric   FabricServer
}

const (
	ORG1_TYPE int = 1
	ORG2_TYPE int = 2
)

var OrgType map[string]int

func init(){
	OrgType = make(map[string]int)
	OrgType["peerorg1"] = ORG1_TYPE
	OrgType["peerorg2"] = ORG2_TYPE
}

var myserver server

func (s *server) InitServer() error {
	s.AuthUser, s.AuthPwd = GetRegisterAuthConfig()
	log.Printf("AuthUser : %s, AuthPwd : %s", s.AuthUser, s.AuthPwd)

	err := s.IdServer.InitSnowflake(GetServerID())
	if err != nil {
		return  err
	}

	configFile := common.GetConfigPath()
	channels, err := common.GetChannels()
	if err != nil {
		return  err
	}

	if GetDatabaseType() == DATATYPE {
		log.Printf("init postgres database...")
		postgresOpt := GetPostgres()
		s.db.InitDB(postgresOpt.User,postgresOpt.Passwd,postgresOpt.DbName, uint32(postgresOpt.DbPort), postgresOpt.DbHost)
		err := s.db.RegisterDB()
		if err != nil {
			log.Fatalf("init postgres database error: %v", err)
		}
		if s.db.Database == nil {
			log.Printf("db.Database is nil")
			return errors.New("db.Database is nil")
		}
		s.db.StatusDB = true
	} else {
		s.db.StatusDB = false
	}


	s.myca = make(map[string]*CA)
	InitCA(configFile)

	for i:=0; i < len(channels); i++ {
		s.myca[channels[i].CA_Org] = new(CA)
		err := s.myca[channels[i].CA_Org].InitCaServer(channels[i].CA_Org, channels[i].EnrollDir)
		if err != nil {
			log.Fatalf("Init CA FAILT Org: ", channels[i].CA_Org, err.Error())
		} else {
			fmt.Println("Init CA SUCCESS Org:", channels[i].CA_Org)
		}
	}

	err = s.fabric.Init()
	if err != nil {
		log.Fatalf("Init Fabric Service failed %v", err)
	}
	return nil
}

func (s *server) UnInitServer() {
	log.Printf("destory database...")
	if s.db.StatusDB {
		s.db.UnRegisterDB()
	}
}


func (s *server) RegisterClient(ctx context.Context, in *pb.RegisterRequest) (*pb.RegisterReply, error) {
	log.Println("_______________________RegisterClient_________________________________")
	if in.User != s.AuthUser && in.Pwd != s.AuthPwd {
		return &pb.RegisterReply{Message:"login user or password is incorrect", Appid:"", Appkey:""}, nil
	}

	if in.Org == "" || in.Affiliation == "" {
		return &pb.RegisterReply{Message:"Org or Affiliation can not be empty", Appid:"", Appkey:""}, nil
	}

	appid := s.IdServer.GetTokenIdFromSnowflake()
	appkey := s.IdServer.GetTokenIdFromSnowflake()

	if appid == "" || appkey == "" {
		return &pb.RegisterReply{Message:"generate appid or appkey error",  Appid:"", Appkey:""}, nil
	}

	for s.db.StatusDB {
		isExist, err := s.db.IsExist(appid)
		if err != nil {
			return &pb.RegisterReply{Message:"query appid exception"}, err
		}

		if isExist == false {
			break
		}
		log.Printf("appid is exist")
		reqNum := 0
		for {
			appid = s.IdServer.GetTokenIdFromSnowflake()
			if appid == "" {
				reqNum++
				if reqNum > 4 {
					break
				}
				continue
			} else {
				break
			}
		}
	}

	dayTimes := GetUTCTimeStr()
	if caService, ok := s.myca[in.Org]; ok {
		if s.db.StatusDB {
			_, err := s.db.InsertAppInfo(appid, appkey, OrgType[caService.OrgId], dayTimes)
			if err != nil {
				return &pb.RegisterReply{Message: "insert appid and appkey failed", Appid: "", Appkey: ""}, err
			}
		}

		fmt.Println("user req prikey:", appid," cert: ", appkey)
		prikey, ecert, err := caService.RegisterAndEnrollUser(appid, appkey, in.Affiliation)
		if err != nil {
			return &pb.RegisterReply{Message:"registering user to ca-server failed", Appid:"", Appkey:""}, err
		}
		_= prikey
		fmt.Println("user prikey:",hex.EncodeToString(prikey[:]))

		return &pb.RegisterReply{Message : "OK",  Appid:appid, Appkey:appkey,
			Prikey:hex.EncodeToString(prikey[:]), Cert:string(ecert[:])}, nil

	} else {
		errMsg := "org is not exist " + in.Org
		return &pb.RegisterReply{Message:"org is not exist", Appid:"", Appkey:""}, errors.New(errMsg)
	}
	return nil,nil
	//return &pb.RegisterReply{Message:"org is not exist", Appid:"", Appkey:""}, nil
}

func (s *server) EnrollAsset(ctx context.Context, in *pb.AssetEnroll) (*pb.ResultsReply, error) {
	log.Println("_______________________EnrollAsset_________________________________")
	if in.Channel == "" || in.Chaincode == "" {
		return &pb.ResultsReply{Message:"chainid or chaincodeid can not be empty", Payload:""}, nil
	}

	if in.Key == "" {
		return &pb.ResultsReply{Message:"Key can not be empty", Payload:""}, nil
	}

	if in.Payload == "" {
		return &pb.ResultsReply{Message:"payload can not be empty", Payload:""}, nil
	}

	err := s.fabric.InitAsset(in.Channel, in.Chaincodepath, in.Chaincode,  in.Chaincodeversion, in.Key, in.Payload)
	if err != nil {
		return &pb.ResultsReply{Message : "execution func enrollasset failed", Payload:""}, err
	}

	return &pb.ResultsReply{Message : "OK", Payload:""}, nil
}

func (s *server) RegisterAsset(ctx context.Context, in *pb.AssetRegister) (*pb.ResultsReply, error) {
	log.Println("_______________________RegisterAsset_________________________________")
	if in.Channel == "" || in.Chaincode == "" {
		return &pb.ResultsReply{Message:"chainid or chaincodeid can not be empty", Payload:""}, nil
	}
	if in.Appid == "" {
		return &pb.ResultsReply{Message:"appid can not be empty", Payload:""}, nil
	}
	if in.Payload == "" {
		return &pb.ResultsReply{Message:"payload can not be empty", Payload:""}, nil
	}

	if s.db.StatusDB {
		isExist, err := s.db.IsExist(in.Appid)
		if err != nil {
			return &pb.ResultsReply{Message: "query database exception", Payload: ""}, err
		}
		if isExist == false {
			return &pb.ResultsReply{Message: "appid does not exist", Payload: ""}, err
		}
	}

	txID, err := s.fabric.InvokeInit(in.Channel, in.Chaincode, in.Appid, in.Key, in.Payload)
	if err != nil {
		return &pb.ResultsReply{Message : "execution func enrollasset failed", Payload:""}, err
	}
	return &pb.ResultsReply{Message : "OK", Payload : txID}, nil
}


func (s *server) TransactionAsset(ctx context.Context, in *pb.Transaction) (*pb.ResultsReply, error) {
	log.Println("_______________________TransactionAsset_________________________________")
	if in.Channel == "" || in.Chaincode == "" {
		return &pb.ResultsReply{Message:"chainid or chaincodeid can not be empty", Payload:""}, nil
	}
	if in.Payload == "" {
		return &pb.ResultsReply{Message:"payload can not be empty", Payload:""}, nil
	}

	if s.db.StatusDB {
		isExist1, err1 := s.db.IsExist(in.Ownerid)
		if err1 != nil {
			return &pb.ResultsReply{Message: "query dbase abort", Payload: ""}, err1
		}
		if isExist1 == false {
			return &pb.ResultsReply{Message: "appid is not exist", Payload: ""}, err1
		}

		isExist2, err2 := s.db.IsExist(in.Receiverid)
		if err2 != nil {
			return &pb.ResultsReply{Message: "query dbase abort", Payload: ""}, err2
		}
		if isExist2 == false {
			return &pb.ResultsReply{Message: "appid is not exist", Payload: ""}, err2
		}
	}

	txID, err := s.fabric.InvokeTransaction(in.Channel, in.Chaincode, in.Ownerid, in.Receiverid, in.Payload)
	if err != nil {
		return &pb.ResultsReply{Message : "execution func enrollasset failed", Payload:""}, err
	}
	return &pb.ResultsReply{Message : "OK", Payload : txID}, nil
}

func (s *server) QueryAsset(ctx context.Context, in *pb.AssetQuery) (*pb.ResultsReply, error) {
	log.Println("_______________________QueryAsset_________________________________")
	if in.Channel == "" || in.Chaincode == "" {
		return &pb.ResultsReply{Message:"chainid or chaincodeid can not be empty", Payload:""}, nil
	}

	if in.Appid == "" {
		return &pb.ResultsReply{Message:"appid can not be empty", Payload:""}, nil
	}

	if in.Key == "" {
		return &pb.ResultsReply{Message:"payload can not be empty", Payload:""}, nil
	}

	if s.db.StatusDB {
		isExist, err := s.db.IsExist(in.Appid)
		if err != nil {
			return &pb.ResultsReply{Message: "query database exception", Payload: ""}, err
		}
		if isExist == false {
			return &pb.ResultsReply{Message: "appid does not exist", Payload: ""}, err
		}
	}

	payload, err := s.fabric.InvokeQuery(in.Channel, in.Chaincode, in.Appid, in.Key)
	if err != nil {
		return &pb.ResultsReply{Message : "execution func InvokeQuery failed", Payload:""}, err
	}
	fmt.Println("payload:",string(payload))
	return &pb.ResultsReply{Message : "OK", Payload : payload}, nil
}

//define main
func setlimit() {
	var rlim syscall.Rlimit
	err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rlim)
	if err != nil {
		fmt.Println("get rlimit error: " + err.Error())
		os.Exit(1)
	}
	rlim.Cur = 65535
	rlim.Max = 65535
	err = syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rlim)
	if err != nil {
		fmt.Println("set rlimit error: " + err.Error())
		os.Exit(1)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello world!")
}



func main() {
	/*
	go func() {
		http.ListenAndServe("0.0.0.0:6060", nil) // 启动默认的 http 服务，可以使用自带的路由
	}()
	*/
	ticker := time.NewTicker(time.Second * 10)
	go func() {
		var m runtime.MemStats
		for {
			select {
				case <-ticker.C:
					runtime.ReadMemStats(&m)
					Tracefile(fmt.Sprintf("%d KB,%d KB,%d KB,%d KB", m.HeapSys/1024, m.HeapAlloc/1024, m.HeapIdle/1024, m.HeapReleased/1024))
					runtime.GC()
			}
		}
	}()

	setlimit()
	myserver.InitServer()
	serveraddr := GetServerAddr()
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
