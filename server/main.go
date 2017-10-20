package main

import (
	. "streamserver/common"
//	. "streamserver/fabric"
	//"github.com/stretchr/testify/assert"
	//"fmt"
	//"reflect"
	"fmt"
)


type User struct  {
	Id int
	Name string
	//addr string
}

func main() {
	channels,_ := GetChannels()
	ShowChannels(channels)

	path := GetConfigPath()
	fmt.Println("path:",path)
	/*
	GetServerAddr()

	var fabric FabricServer
	fabric.Init()


	var ca CA

	ca.InitCaServer("peerorg1")

	ca.RegisterAndEnrollUser("lhy","mobile")
	*/

}