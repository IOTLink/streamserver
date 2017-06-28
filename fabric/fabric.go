package fabric

import (
	"log"
	//fcutil "github.com/hyperledger/fabric-sdk-go/fabric-client/util"
	//fabricClient "github.com/hyperledger/fabric-sdk-go/fabric-client"
)

type FabricServer struct {
	setup BaseSetupImpl
}

func (fabric *FabricServer)Init() {
	testSetup := BaseSetupImpl{
		ConfigFile:      "./config_test.yaml",
		ChainID:         "mychannel",
		ChannelConfig:   "./channel.tx",
		ConnectEventHub: true,
		ChainCodeID:     "mychaincodev10",

	}
	if err := testSetup.Initialize(); err != nil {
		log.Fatalf("fabric server init abort %s",err.Error())
	}
	fabric.setup = testSetup
}

func (fabric *FabricServer)InitAsset(appid string, payload string) error{
	/*
	testSetup := BaseSetupImpl{
		ConfigFile:      "./config_test.yaml",
		ChainID:         "mychannel",
		ChannelConfig:   "./channel.tx",
		ConnectEventHub: true,
		ChainCodeID:     "mychaincodev5",
	}
	if err := testSetup.Initialize(); err != nil {
		log.Fatalf("fabric server init abort %s",err.Error())
	}
	fabric.setup = testSetup
	*/
	user, err := GetUser(fabric.setup.Client, appid, "")
	if err != nil {

	}
	fabric.setup.Client.SetUserContext(user)

	if err := fabric.setup.InitCC(appid, payload); err != nil {
		log.Fatalf("InstallAndInstantiateExampleCC return error: %v", err)
		return err
	}
	return nil
}

func (fabric *FabricServer)Init2Asset(appid string, payload string) (string,error){
	user, err := GetUser(fabric.setup.Client, appid, "")
	if err != nil {

	}
	fabric.setup.Client.SetUserContext(user)

	var args []string
	args = append(args, "invoke")
	args = append(args, "initins")
	args = append(args, appid)
	args = append(args, payload)

	return fabric.setup.InvokeInit(fabric.setup.ChainID, fabric.setup.ChainCodeID, args)
}

func (fabric *FabricServer)Query(appid string) (string, error){
	user, err := GetUser(fabric.setup.Client, appid, "")
	if err != nil {

	}
	fabric.setup.Client.SetUserContext(user)

	var args []string
	args = append(args, "invoke")
	args = append(args, "query")
	args = append(args, appid)

	return fabric.setup.InvokeQuery(fabric.setup.ChainID, fabric.setup.ChainCodeID, args)
}


func (fabric *FabricServer)Transfer(appidA string, appidB string, value string) (string, error) {
	user, err := GetUser(fabric.setup.Client, appidA, "")
	if err != nil {

	}
	fabric.setup.Client.SetUserContext(user)
	//eventID := "test([a-zA-Z]+)"
	//done, rce := util.RegisterCCEvent(fabric.setup.ChainCodeID, eventID, fabric.setup.EventHub)
	var args []string
	args = append(args, "invoke")
	args = append(args, "transfer")
	args = append(args, appidA)
	args = append(args, appidB)
	args = append(args, value)
	txId, err := fabric.setup.InvokeTransfer(fabric.setup.ChainID, fabric.setup.ChainCodeID, args)
	if err != nil {
		log.Fatalf("Move funds return error: %v", err)
	}

	/*
	select {
	case <-done:
	case <-time.After(time.Second * 20):
		log.Fatalf("Did NOT receive CC for eventId(%s)\n", eventID)
	}

	fabric.setup.EventHub.UnregisterChaincodeEvent(rce)
	*/
	return txId,nil
}


func (fabric *FabricServer)Delete(appid string) (string,error){
	user, err := GetUser(fabric.setup.Client, appid, "")
	if err != nil {

	}
	fabric.setup.Client.SetUserContext(user)

	var args []string
	args = append(args, "invoke")
	args = append(args, "delete")
	args = append(args, appid)

	return fabric.setup.InvokeDelete(fabric.setup.ChainID, fabric.setup.ChainCodeID, args)
}

