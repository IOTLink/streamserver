package fabric

import (
	"log"
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
		ChainCodeID:     "mychaincodev6",

	}
	if err := testSetup.Initialize(); err != nil {
		log.Fatalf("fabric server init abort %s",err.Error())
	}
	fabric.setup = testSetup

}

func (fabric *FabricServer)Invoke(appid string) error{

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

	if err := fabric.setup.InstallAndInstantiateExampleCC(appid); err != nil {
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