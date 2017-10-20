package common

import (
	"fmt"
	"github.com/spf13/viper"
	"strings"
	"log"
	//"strconv"
	"strconv"
	//"encoding/json"
	"path/filepath"
	"os"
)

var myViper = viper.New()

type ChannelConfig struct {
	ChannelID string
	ChannelConfig string
	OrgID []string
	ConnectEventHub bool
	CA_Org string
	EnrollDir string
}

const (
	SERVERADDR = "0.0.0.0:50055"
	SERVERID = 1
	CONFIGPATH = "/conf/config.yaml"

	AUTHCLIENT = "admin"
	AUTHPASSWD = "adminw"

	DATATYPE = "postgres"
	POSTGRES_USER = "root"
	POSTGRES_PASSWD = "123456"
	POSTGRES_DBNAME = "fabric"
	POSTGRES_DBPORT = 5432
	POSTGRES_DBHOST = "127.0.0.1"

)

func init() {
	execPath := getCurrentDirectory()
	myViper.SetEnvPrefix("streamserver")
	myViper.AutomaticEnv()
	replacer := strings.NewReplacer(".", "_")
	myViper.SetEnvKeyReplacer(replacer)

	execPath = execPath + "/conf/server.yaml"
	myViper.SetConfigFile(execPath)
	myViper.SetConfigType("yaml")
	err := myViper.ReadInConfig()
	if err != nil {
		fmt.Println(err.Error())
		log.Fatalf("Fatal error config file: %v", err)
		return
	}

	if err == nil {
		log.Println("Using config file: ", myViper.ConfigFileUsed())
	} else {
		log.Fatalf("Fatal error config file: %v", err)
		return
	}
}

func  GetChannelConfig() (map[string]interface{}, error) {
	channelflag := myViper.InConfig("channels")
	if channelflag == true {
		groups := myViper.GetStringMap("channels")
		return groups, nil
	}
	return nil, fmt.Errorf("not found channels in config file")
}

func GetChannels() ([]ChannelConfig, error) {
	mapChannels, err := GetChannelConfig()
	if err != nil {
		return nil, err
	}

	i := 0
	channels := make([]ChannelConfig, len(mapChannels))
	for key, value := range mapChannels {
		log.Println("key:", key, "value:", value)

		m := value.(map[string]interface{})
		for k, v := range m {
			switch vv := v.(type) {
			case string:
				if k == "orgid" {
					channels[i].OrgID = strings.Split(vv, ",")
				} else if k == "channelid" {
					channels[i].ChannelID = vv
				} else if k == "channelconfig" {
					channels[i].ChannelConfig = vv
				} else if k == "ca_org" {
					channels[i].CA_Org = vv
				} else if k == "enrolldir" {
					channels[i].EnrollDir = vv
				}
			case bool:
				if k == "connecteventhub" {
					channels[i].ConnectEventHub = vv
				}
			default:
				log.Println(k, "is of a type I don't know how to handle")
			}
		}
		log.Println("GetChannels:", channels[i].ChannelID,channels[i].OrgID,channels[i].EnrollDir,channels[i].ConnectEventHub,channels[i].ChannelConfig)
		i++
	}
	return channels,nil
}

func ShowChannels(channels []ChannelConfig) {
	if channels != nil {
		for _, channel := range channels {
			log.Println("ChannelID:",channel.ChannelID)
			log.Println("ChannelConfig:",channel.ChannelConfig)
			log.Println("ConnectEventHub:",channel.ConnectEventHub)
			log.Println("OrgID:",channel.OrgID)
			log.Println("EnrollDir:",channel.EnrollDir)
			log.Println()
		}
	}
}


func GetServerAddr() string {
	serverAddr := myViper.GetString("server.serverAddr")
	log.Println("server listen addr:",serverAddr)
	if serverAddr == "" {
		return SERVERADDR
	}
	return serverAddr
}

func GetServerID() int64 {
	serverId := myViper.GetString("server.serverId")
	log.Println("server serverId :",serverId)
	if serverId == "" {
		return SERVERID
	}
	id, err :=strconv.ParseInt(serverId, 10, 64)
	if err != nil {
		return SERVERID
	}
	return id
}

func GetConfigPath() string {
	execPath := getCurrentDirectory()
	configPath := myViper.GetString("configFile.path")
	log.Println("server listen addr:",configPath)
	if configPath == "" {
		return execPath + CONFIGPATH
	}
	return execPath + configPath
}

func GetRegisterAuthConfig() (string, string) {
	client := myViper.GetString("registerAuth.client")
	if client == "" {
		client = AUTHCLIENT
	}

	passwd := myViper.GetString("registerAuth.passwd")
	if passwd == "" {
		passwd = AUTHPASSWD
	}
	return client,passwd
}


//db
func GetDatabaseType() string {
	dbtype := myViper.GetString("dataBase.type")
	if dbtype == "" {
		return DATATYPE
	}
	return dbtype
}

type PostgresConfig struct {
	User string
	Passwd string
	DbName string
	DbPort int64
	DbHost string
}

func GetPostgres() PostgresConfig {
	postgres := PostgresConfig{
		User:   myViper.GetString("postgres.user"),
		Passwd: myViper.GetString("postgres.passwd"),
		DbName: myViper.GetString("postgres.dbname"),
		DbPort: myViper.GetInt64("postgres.dbport"),
		DbHost: myViper.GetString("postgres.dbhost"),
	}
	if postgres.User == "" {
		postgres.User = POSTGRES_USER
	}
	if postgres.Passwd == "" {
		postgres.Passwd = POSTGRES_PASSWD
	}
	if postgres.DbName == "" {
		postgres.DbName = POSTGRES_DBNAME
	}
	if postgres.DbPort == 0 {
		postgres.DbPort = POSTGRES_DBPORT
	}
	if postgres.DbHost == "" {
		postgres.DbHost = POSTGRES_DBHOST
	}
	return postgres
}

func substr(s string, pos, length int) string {
	runes := []rune(s)
	l := pos + length
	if l > len(runes) {
		l = len(runes)
	}
	return string(runes[pos:l])
}

func getParentDirectory(dirctory string) string {
	return substr(dirctory, 0, strings.LastIndex(dirctory, "/"))
}

func getCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}