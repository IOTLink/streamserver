package fabric

import (
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/config"
	"github.com/hyperledger/fabric-sdk-go/fabric-client/events"

	fabricClient "github.com/hyperledger/fabric-sdk-go/fabric-client"
	fcutil "github.com/hyperledger/fabric-sdk-go/fabric-client/util"
	bccspFactory "github.com/hyperledger/fabric/bccsp/factory"
	"os"
	"path"
	"time"
)

// BaseSetupImpl implementation of BaseTestSetup
type BaseSetupImpl struct {
	Client             fabricClient.Client
	OrdererAdminClient fabricClient.Client
	Chain              fabricClient.Chain
	EventHub           events.EventHub
	ConnectEventHub    bool
	ConfigFile         string
	ChainID            string
	ChainCodeID        string
	Initialized        bool
	ChannelConfig      string
}

// GetOrdererAdmin ...
func GetOrdererAdmin(c fabricClient.Client) (fabricClient.User, error) {
	keyDir := "ordererOrganizations/example.com/users/Admin@example.com/msp/keystore"
	certDir := "ordererOrganizations/example.com/users/Admin@example.com/msp/signcerts"
	return fcutil.GetPreEnrolledUser(c, keyDir, certDir, "ordererAdmin")
}

// GetAdmin ...
func GetAdmin(c fabricClient.Client, userOrg string) (fabricClient.User, error) {
	keyDir := fmt.Sprintf("peerOrganizations/%s.example.com/users/Admin@%s.example.com/msp/keystore", userOrg, userOrg)
	certDir := fmt.Sprintf("peerOrganizations/%s.example.com/users/Admin@%s.example.com/msp/signcerts", userOrg, userOrg)
	username := fmt.Sprintf("peer%sAdmin", userOrg)
	return fcutil.GetPreEnrolledUser(c, keyDir, certDir, username)
}

//GetUser
func GetUser(c fabricClient.Client, name string, pwd string) (fabricClient.User, error) {
	return fcutil.GetMember(c, name, pwd)
}

// Initialize reads configuration from file and sets up client, chain and event hub
func (setup *BaseSetupImpl) Initialize() error {

	if err := setup.InitConfig(); err != nil {
		return fmt.Errorf("Init from config failed: %v", err)
	}

	// Initialize bccsp factories before calling get client
	err := bccspFactory.InitFactories(config.GetCSPConfig())
	if err != nil {
		return fmt.Errorf("Failed getting ephemeral software-based BCCSP [%s]", err)
	}

	client, err := fcutil.GetClient("admin", "adminpw", "/tmp/enroll_user")
	if err != nil {
		return fmt.Errorf("Create client failed: %v", err)
	}
	//clientUser := client.GetUserContext()

	setup.Client = client

	org1Admin, err := GetAdmin(client, "org1")
	if err != nil {
		return fmt.Errorf("Error getting org admin user: %v", err)
	}

	chain, err := fcutil.GetChain(setup.Client, setup.ChainID)
	if err != nil {
		return fmt.Errorf("Create chain (%s) failed: %v", setup.ChainID, err)
	}
	setup.Chain = chain

	ordererAdmin, err := GetOrdererAdmin(client)
	if err != nil  {
		return fmt.Errorf("Error getting orderer admin user: %v", err)
	}
	if ordererAdmin != nil {

	}
	// Create and join channel

	if err := fcutil.CreateAndJoinChannel(client, ordererAdmin, org1Admin, chain, setup.ChannelConfig); err != nil {
		return fmt.Errorf("CreateAndJoinChannel return error: %v", err)
	}

	client.SetUserContext(org1Admin)

	if err := setup.setupEventHub(client); err != nil {
		return err
	}

	setup.Initialized = true

	return nil
}

// getEventHub initilizes the event hub
func getEventHub(client fabricClient.Client) (events.EventHub, error) {
	eventHub, err := events.NewEventHub(client)
	if err != nil {
		return nil, fmt.Errorf("Error creating new event hub: %v", err)
	}
	foundEventHub := false
	peerConfig, err := config.GetPeersConfig()
	if err != nil {
		return nil, fmt.Errorf("Error reading peer config: %v", err)
	}
	for _, p := range peerConfig {
		if p.EventHost != "" && p.EventPort != 0 {
			fmt.Printf("******* EventHub connect to peer (%s:%d) *******\n", p.EventHost, p.EventPort)
			eventHub.SetPeerAddr(fmt.Sprintf("%s:%d", p.EventHost, p.EventPort),
				p.TLS.Certificate, p.TLS.ServerHostOverride)
			foundEventHub = true
			break
		}
	}

	if !foundEventHub {
		return nil, fmt.Errorf("No EventHub configuration found")
	}

	return eventHub, nil
}

func (setup *BaseSetupImpl) setupEventHub(client fabricClient.Client) error {
	eventHub, err := getEventHub(client)
	if err != nil {
		return err
	}

	if setup.ConnectEventHub {
		if err := eventHub.Connect(); err != nil {
			return fmt.Errorf("Failed eventHub.Connect() [%s]", err)
		}
	}
	setup.EventHub = eventHub

	return nil
}

// InitConfig ...
func (setup *BaseSetupImpl) InitConfig() error {
	if err := config.InitConfig(setup.ConfigFile); err != nil {
		return err
	}
	return nil
}

func (setup *BaseSetupImpl) GetDeployPath() string {
	pwd, _ := os.Getwd()
	return path.Join(pwd, "../")
}

// start chaincode deal  . deploy chaincode and init
func (setup *BaseSetupImpl) InitCC(appid string, payload string) error {
	chainCodePath := "github.com/mychaincode"
	chainCodeVersion := "v0"


	if err := setup.InstallCC(setup.ChainCodeID, chainCodePath, chainCodeVersion, nil); err != nil {
		return err
	}

	var args []string
	//valueStr := strconv.FormatInt(int64(value), 10)
	args = append(args, "init")
	args = append(args, appid)
	args = append(args, payload)
	return setup.InstantiateCC(setup.ChainCodeID, setup.ChainID, chainCodePath, chainCodeVersion, args)
}

func (setup *BaseSetupImpl) InstallCC(chainCodeID string, chainCodePath string, chainCodeVersion string, chaincodePackage []byte) error {
	if err := fcutil.SendInstallCC(setup.Client, setup.Chain, chainCodeID, chainCodePath, chainCodeVersion, chaincodePackage, setup.Chain.GetPeers(), setup.GetDeployPath()); err != nil {
		return fmt.Errorf("SendInstallProposal return error: %v", err)
	}
	return nil
}

func (setup *BaseSetupImpl) InstantiateCC(chainCodeID string, chainID string, chainCodePath string, chainCodeVersion string, args []string) error {
	if err := fcutil.SendInstantiateCC(setup.Chain, chainCodeID, chainID, args, chainCodePath, chainCodeVersion, []fabricClient.Peer{setup.Chain.GetPrimaryPeer()}, setup.EventHub); err != nil {
		return err
	}
	return nil
}
// end chaincode deal  . deploy chaincode and init

// Init ...
func (setup *BaseSetupImpl) InvokeInit(chainID string, chainCodeID string, args []string) (string, error) {
	transientDataMap := make(map[string][]byte)
	transientDataMap["result"] = []byte("Transient data in move funds...")

	transactionProposalResponses, txID, err := fcutil.CreateAndSendTransactionProposal(setup.Chain, chainCodeID, chainID, args, []fabricClient.Peer{setup.Chain.GetPrimaryPeer()}, nil)
	if err != nil {
		return "", fmt.Errorf("CreateAndSendTransactionProposal return error: %v", err)
	}
	// Register for commit event
	done, fail := fcutil.RegisterTxEvent(txID, setup.EventHub)

	txResponse, err := fcutil.CreateAndSendTransaction(setup.Chain, transactionProposalResponses)
	if err != nil {
		return "", fmt.Errorf("CreateAndSendTransaction return error: %v", err)
	}
	fmt.Println(txResponse)

	select {
	case <-done:
	case <-fail:
		return "", fmt.Errorf("invoke Error received from eventhub for txid(%s) error(%v)", txID, fail)
	case <-time.After(time.Second * 30):
		return "", fmt.Errorf("invoke Didn't receive block event for txid(%s)", txID)
	}

	return txID, nil
	//return string(transactionProposalResponses[0].GetResponsePayload()), nil
}

// Query ...
func (setup *BaseSetupImpl) InvokeQuery(chainID string, chainCodeID string, args []string) (string, error) {
	transactionProposalResponses, _, err := fcutil.CreateAndSendTransactionProposal(setup.Chain, chainCodeID, chainID, args, []fabricClient.Peer{setup.Chain.GetPrimaryPeer()}, nil)
	if err != nil {
		return "", fmt.Errorf("CreateAndSendTransactionProposal return error: %v", err)
	}
	return string(transactionProposalResponses[0].GetResponsePayload()), nil
}

func (setup *BaseSetupImpl)InvokeTransfer(chainID string, chainCodeID string, args []string) (string, error){
	transientDataMap := make(map[string][]byte)
	transientDataMap["result"] = []byte("Transient data in move funds...")

	transactionProposalResponse, txID, err := fcutil.CreateAndSendTransactionProposal(setup.Chain, setup.ChainCodeID, setup.ChainID, args, []fabricClient.Peer{setup.Chain.GetPrimaryPeer()}, transientDataMap)
	if err != nil {
		return "", fmt.Errorf("CreateAndSendTransactionProposal return error: %v", err)
	}
	// Register for commit event
	done, fail := fcutil.RegisterTxEvent(txID, setup.EventHub)

	txResponse, err := fcutil.CreateAndSendTransaction(setup.Chain, transactionProposalResponse)
	if err != nil {
		return "", fmt.Errorf("CreateAndSendTransaction return error: %v", err)
	}
	fmt.Println(txResponse)

	select {
	case <-done:
	case <-fail:
		return "", fmt.Errorf("invoke Error received from eventhub for txid(%s) error(%v)", txID, fail)
	case <-time.After(time.Second * 30):
		return "", fmt.Errorf("invoke Didn't receive block event for txid(%s)", txID)
	}

	return txID, nil
}

func (setup *BaseSetupImpl)InvokeDelete(chainID string, chainCodeID string, args []string) (string, error){
	transientDataMap := make(map[string][]byte)
	transientDataMap["result"] = []byte("Transient data in move funds...")

	transactionProposalResponse, txID, err := fcutil.CreateAndSendTransactionProposal(setup.Chain, setup.ChainCodeID, setup.ChainID, args, []fabricClient.Peer{setup.Chain.GetPrimaryPeer()}, transientDataMap)
	if err != nil {
		return "", fmt.Errorf("CreateAndSendTransactionProposal return error: %v", err)
	}
	// Register for commit event
	done, fail := fcutil.RegisterTxEvent(txID, setup.EventHub)

	txResponse, err := fcutil.CreateAndSendTransaction(setup.Chain, transactionProposalResponse)
	if err != nil {
		return "", fmt.Errorf("CreateAndSendTransaction return error: %v", err)
	}
	fmt.Println(txResponse)

	select {
	case <-done:
	case <-fail:
		return "", fmt.Errorf("invoke Error received from eventhub for txid(%s) error(%v)", txID, fail)
	case <-time.After(time.Second * 30):
		return "", fmt.Errorf("invoke Didn't receive block event for txid(%s)", txID)
	}

	return txID, nil
}



