package models

import(
	"github.com/bwmarrin/snowflake"
	"fmt"
	"encoding/hex"
)

func GetAppId() string {
	randData := RandomData()
	if randData == nil {
		return ""
	}
	appId := GetMd5(randData)
	return hex.EncodeToString(appId[:])
}

func GetAppKey() string {
	randData := RandomData()
	if randData == nil {
		return ""
	}
	appkey := GetMd5(randData)
	return hex.EncodeToString(appkey[:])
}


type IDServer struct{
	IdServer *snowflake.Node
}

func (Ids *IDServer)InitSnowflake(nodeId int64) error {
	var err error
	Ids.IdServer, err = snowflake.NewNode(nodeId)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (Ids *IDServer)GetTokenIdFromSnowflake() string {
	// Generate a snowflake ID.
	id := Ids.IdServer.Generate()
	//return id.Base58()
	byteArry := []byte(id.Base58())
	return hex.EncodeToString(byteArry[:])
}