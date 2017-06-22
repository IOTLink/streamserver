package fabric


import (

	config "github.com/hyperledger/fabric-sdk-go/config"
	fabricClient "github.com/hyperledger/fabric-sdk-go/fabric-client"
	kvs "github.com/hyperledger/fabric-sdk-go/fabric-client/keyvaluestore"
	bccspFactory "github.com/hyperledger/fabric/bccsp/factory"

	fabricCAClient "github.com/hyperledger/fabric-sdk-go/fabric-ca-client"
	"encoding/pem"
	"crypto/x509"
	"log"
	"fmt"
)


type CA struct{
	Cli fabricClient.Client
	CaService fabricCAClient.Services
	Admin fabricClient.User
}

func (c *CA)InitCaServer() {
	testSetup := BaseSetupImpl{
		ConfigFile: "./config_test.yaml",
	}

	testSetup.InitConfig()
	client := fabricClient.NewClient()

	err := bccspFactory.InitFactories(config.GetCSPConfig())
	if err != nil {
		log.Fatalf("Failed getting ephemeral software-based BCCSP [%s]", err)
	}

	cryptoSuite := bccspFactory.GetDefault()
	client.SetCryptoSuite(cryptoSuite)
	stateStore, err := kvs.CreateNewFileKeyValueStore("./enroll_user")
	if err != nil {
		log.Fatalf("CreateNewFileKeyValueStore return error[%s]", err)
	}
	client.SetStateStore(stateStore)

	caClient, err := fabricCAClient.NewFabricCAClient()
	if err != nil {
		log.Fatalf("NewFabricCAClient return error: %v", err)
	}

	// Admin user is used to register, enrol and revoke a test user
	adminUser, err := client.LoadUserFromStateStore("admin")

	if err != nil {
		log.Fatalf("client.LoadUserFromStateStore return error: %v", err)
	}
	if adminUser == nil {
		key, cert, err := caClient.Enroll("admin", "adminpw")
		if err != nil {
			log.Fatalf("Enroll return error: %v", err)
		}
		if key == nil {
			log.Fatalf("private key return from Enroll is nil")
		}
		if cert == nil {
			log.Fatalf("cert return from Enroll is nil")
		}

		certPem, _ := pem.Decode(cert)
		if err != nil {
			log.Fatalf("pem Decode return error: %v", err)
		}

		cert509, err := x509.ParseCertificate(certPem.Bytes)
		if err != nil {
			log.Fatalf("x509 ParseCertificate return error: %v", err)
		}
		if cert509.Subject.CommonName != "admin" {
			log.Fatalf("CommonName in x509 cert is not the enrollmentID")
		}
		adminUser = fabricClient.NewUser("admin")
		adminUser.SetPrivateKey(key)
		adminUser.SetEnrollmentCertificate(cert)
		err = client.SaveUserToStateStore(adminUser, false)
		if err != nil {
			log.Fatalf("client.SaveUserToStateStore return error: %v", err)
		}
		adminUser, err = client.LoadUserFromStateStore("admin")
		if err != nil {
			log.Fatalf("client.LoadUserFromStateStore return error: %v", err)
		}
		if adminUser == nil {
			log.Fatalf("client.LoadUserFromStateStore return nil")
		}
	}
	c.Cli = client
	c.CaService = caClient
	c.Admin = adminUser
}

func (c *CA)RegisterAndEnrollUser(appid string, appkey string)  error {
	registerRequest := fabricCAClient.RegistrationRequest{
		Name:        appid,
		Type:        "user",
		Affiliation: "org1.department1",
		CAName:      config.GetFabricCAName(),
		Secret:      appkey,
	}
	enrolmentSecret, err := c.CaService.Register(c.Admin, &registerRequest)
	if err != nil {
		log.Fatalf("Error from Register: %s", err)
		return err
	}
	fmt.Printf("Registered User: %s, Secret: %s\n", appid, enrolmentSecret)

	// Enrol the previously registered user
	ekey, ecert, err := c.CaService.Enroll(appid, enrolmentSecret)
	if err != nil {
		log.Fatalf("Error enroling user: %s", err.Error())
		return err
	}

	enrolleduser := fabricClient.NewUser(appid)
	enrolleduser.SetEnrollmentCertificate(ecert)
	enrolleduser.SetPrivateKey(ekey)

	err = c.Cli.SaveUserToStateStore(enrolleduser, false)
	if err != nil {
		log.Fatalf("client.SaveUserToStateStore return error: %v", err)
		return err
	}
	return nil
}

func (c *CA)LoadUser(appid string) (fabricClient.User, error){

	user, err := c.Cli.LoadUserFromStateStore(appid)
	if err != nil {
		log.Fatalf("client.LoadUserFromStateStore return error: %v", err)
		return nil,err
	}
	return user,nil
}

