package fabric

import (
	"log"
	//"fmt"
	//"time"
	//"github.com/hyperledger/fabric-sdk-go/fabric-client/util"
	//"time"
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
		ChainCodeID:     "mychaincodev8",

	}
	if err := testSetup.Initialize(); err != nil {
		log.Fatalf("fabric server init abort %s",err.Error())
	}
	fabric.setup = testSetup

}

func (fabric *FabricServer)Invoke(appid string, value int32) error{
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

	if err := fabric.setup.InstallAndInstantiateExampleCC(appid, value); err != nil {
		log.Fatalf("InstallAndInstantiateExampleCC return error: %v", err)
		return err
	}
	return nil

}

func (fabric *FabricServer)Query(appid string) (string, error){
	/*
	testSetup := BaseSetupImpl{
		ConfigFile:      "./config_test.yaml",
		ChainID:         "mychannel",
		ChannelConfig:   "./channel.tx",
		ConnectEventHub: true,
		ChainCodeID:     "mychaincodev6",
	}
	if err := testSetup.Initialize(); err != nil {
		log.Fatalf("fabric server init abort %s",err.Error())
	}
	fabric.setup = testSetup

	if err := fabric.setup.InstallAndInstantiateExampleCC(appid); err != nil {
		log.Fatalf("InstallAndInstantiateExampleCC return error: %v", err)
		return "",err
	}
	*/
	var args []string
	args = append(args, "invoke")
	args = append(args, "query")
	args = append(args, appid)
	return fabric.setup.Query(fabric.setup.ChainID, fabric.setup.ChainCodeID, args)
}

func (fabric *FabricServer)Move(appidA string, appidB string, value string) (string, error) {

	//eventID := "test([a-zA-Z]+)"
	//done, rce := util.RegisterCCEvent(fabric.setup.ChainCodeID, eventID, fabric.setup.EventHub)

	txId, err := fabric.setup.MoveFunds(appidA, appidB, value)
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