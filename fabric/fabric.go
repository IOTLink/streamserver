package fabric

import (
	"log"
	common "streamserver/common"
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/api/apitxn"
	"errors"
)

type FabricServer struct {
	setupMap map[string]*BaseSetupImpl
}

func (fabric *FabricServer) Init() error {
	fabric.setupMap = make(map[string]*BaseSetupImpl)
	configPath := common.GetConfigPath()
	channels,_ := common.GetChannels()
	if channels != nil {
		for i:= 0; i < len(channels); i++ {
			if channels[i].ChannelID == "" {
				errors.New("channelid is null")
			}
			fabric.setupMap[channels[i].ChannelID] = &BaseSetupImpl{
				ConfigFile:      configPath,
				ChannelID:       channels[i].ChannelID,
				ChannelConfig:   channels[i].ChannelConfig,
				ConnectEventHub: channels[i].ConnectEventHub,
				OrgID:           channels[i].OrgID[0],
				EnrollUserDir:   channels[i].EnrollDir,
			}
		}
	}

	for i:= 0; i < len(channels); i++ {
		if baseSetup, ok := fabric.setupMap[channels[i].ChannelID]; ok {
			if err := baseSetup.Initialize(); err != nil {
				log.Println("fabric server init abort %s",err.Error())
				return err
			} else {
				fmt.Println("Success Init Channel:",channels[i].ChannelID)
			}
		} else {
			fmt.Println("not find channelid",channels[i].ChannelID)
		}
	}
	return nil
}

//install chaincode
func (fabric *FabricServer)InitAsset(channelid string, chaincodepath, chaincode, chaincodeversion string ,key string, payload string) error{
	if channelid == "" || chaincode == "" || chaincodeversion == "" {
		return errors.New("parameter can not be empty")
	}

	if baseSetup, ok := fabric.setupMap[channelid]; ok {
		if err := baseSetup.InitCC(chaincodepath, chaincode, chaincodeversion, key, payload); err != nil {
			log.Println("install chaincode return error: %v", err)
			return err
		}
	} else {
		errMsg := "not find channelid " + channelid
		return errors.New(errMsg)
	}
	return nil
}

func (fabric *FabricServer) InvokeInit(channelid string, chaincode string, appid string, key string, payload string) (string,error) {
	if channelid == "" || chaincode == "" || appid == "" || key == "" || payload == "" {
		return fmt.Sprintf("%s","parameter can not be empty"), errors.New("parameter can not be empty")
	}
	fmt.Println("Channel Nmae:",channelid)
	var txID string
	var err error
	fcn := "invoke"

	var args []string
	args = append(args, "funcinit")
	args = append(args, key)
	args = append(args, payload)

	transientDataMap := make(map[string][]byte)
	transientDataMap["result"] = []byte("Transient data in move funds...")
	//defer runtime.GC()

	if baseSetup, ok := fabric.setupMap[channelid]; ok {
		sdk  := baseSetup.FabricSDK
		client, nil := sdk.NewSystemClient(nil)
		user, err := client.GetPreEnrolledUser(appid)
		//user, err := client.LoadUserFromStateStore(appid) //mem leaks
		if err != nil {
			log.Println(appid, "LoadUserFromStateStore ERROR:",err.Error())
		} else {
			log.Println("success load user appid:", appid,user)
		}

		channel, err := baseSetup.GetChannel(client, channelid, []string{baseSetup.OrgID})
		if err != nil {
			return fmt.Sprintf("Create channel %s failed: %v", channelid, err), err
		}
		txID, err = baseSetup.InvokeFunc(client, channel, []apitxn.ProposalProcessor{channel.PrimaryPeer()}, baseSetup.EventHub, chaincode, fcn, args, transientDataMap)
		if err != nil {
			return "", err
		}
		//txID, err = baseSetup.InvokeFunc(baseSetup.Client, baseSetup.Channel, []apitxn.ProposalProcessor{baseSetup.Channel.PrimaryPeer()}, baseSetup.EventHub, chaincode, fcn, args, transientDataMap)
	} else {
		errMsg := "not find channelid " + channelid
		return "", errors.New(errMsg)
	}
	return 	txID, err
}

func (fabric *FabricServer) InvokeTransaction(channelid string, chaincode string , ownerid string, receiverid string, payload string) (string, error) {
	var txID string
	var err error
	fcn := "invoke"

	var args []string
	args = append(args, "functransaction")
	args = append(args, ownerid)
	args = append(args, receiverid)
	args = append(args, payload)

	transientDataMap := make(map[string][]byte)
	transientDataMap["result"] = []byte("Transient data in move funds...")
	//defer runtime.GC()
	
	if baseSetup, ok := fabric.setupMap[channelid]; ok {
		sdk  := baseSetup.FabricSDK
		client, nil := sdk.NewSystemClient(nil)
		//user, err := client.LoadUserFromStateStore(ownerid)
		user, err := client.GetPreEnrolledUser(ownerid)
		if err != nil {
			log.Println(ownerid, "LoadUserFromStateStore ERROR:",err.Error())
		} else {
			log.Println("success load user appid:", ownerid,user)
		}

		channel, err := baseSetup.GetChannel(client, channelid, []string{baseSetup.OrgID})
		if err != nil {
			return fmt.Sprintf("Create channel %s failed: %v", channelid, err), err
		}
		txID, err = baseSetup.InvokeFunc(client, channel, []apitxn.ProposalProcessor{channel.PrimaryPeer()}, baseSetup.EventHub, chaincode, fcn, args, transientDataMap)
		if err != nil {
			return "", err
		}
	} else {
		errMsg := "not find channelid " + channelid
		return fmt.Sprintf("%s", errMsg), errors.New(errMsg)
	}
	return 	txID, err
}

func (fabric *FabricServer) InvokeQuery(channelid string, chaincode string, appid string, key string) (string,error) {
	var payload string
	var err error
	fcn := "invoke"

	var args []string
	args = append(args, "funcquery")
	args = append(args, key)

	if baseSetup, ok := fabric.setupMap[channelid]; ok {
		Client  := baseSetup.Client
		Channel := baseSetup.Channel
		payload, err = baseSetup.InvokeFuncQuery(Client, Channel, chaincode, fcn, args)
		if err != nil {
			return "", err
		}
	} else {
		errMsg := "not find channelid " + channelid
		return fmt.Sprintf("%s", errMsg), errors.New(errMsg)
	}
	return 	payload, err
}
