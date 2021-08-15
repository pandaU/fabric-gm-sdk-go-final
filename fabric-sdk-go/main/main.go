package main

import (
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/cryptosuite"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)
var (
	cc            = ""
	user          = ""
	secret        = ""
	channelName   = ""
	chaincodrPath = "github.com/VoneChain-CS/fabric-gm/scripts/fabric-samples/chaincode/abstore/go"
)
func populateWallet(wallet *gateway.Wallet) error {
	log.Println("============ Populating wallet ============")
	credPath := "C:\\Users\\xxx\\Desktop\\gm-sdk\\fabric-sdk-go-v1.0.0-gm\\main\\organizations\\peerOrganizations\\org1.xxzx.com\\users\\Admin@org1.xxzx.com\\msp"

	certPath := filepath.Join(credPath, "signcerts", "cert.pem")
	// read the certificate pem
	cert, err := ioutil.ReadFile(filepath.Clean(certPath))
	if err != nil {
		return err
	}

	keyDir := filepath.Join(credPath, "keystore")
	// there's a single file in this dir containing the private key
	files, err := ioutil.ReadDir(keyDir)
	if err != nil {
		return err
	}
	if len(files) != 1 {
		return fmt.Errorf("keystore folder should have contain one file")
	}
	keyPath := filepath.Join(keyDir, files[0].Name())
	key, err := ioutil.ReadFile(filepath.Clean(keyPath))
	if err != nil {
		return err
	}

	identity := gateway.NewX509Identity("Org1MSP", string(cert), string(key))

	return wallet.Put("appUser", identity)
}
func main() {
	user = "admin"
	secret = "adminpw"
	channelName = "mychannel"
	cc = "mycc_3"
	//fmt.Println("Reading connection profile..")
	//c := config.FromFile("C:\\Users\\xxx\\Desktop\\gm-sdk\\fabric-sdk-go-v1.0.0-gm\\main\\config_test.yaml")
	//sdk, err := fabsdk.New(c)
	//if err != nil {
	//	fmt.Printf("Failed to create new SDK: %s\n", err)
	//	os.Exit(1)
	//}
	//defer sdk.Close()
	//
	////registerUser(user,secret,sdk)
	//sdk.Config()
	//enrollUser(sdk)
	//clientChannelContext := sdk.ChannelContext(channelName, fabsdk.WithUser(user))
	//client, err := channel.New(clientChannelContext)
	//resp := queryCC(client, []byte("user"), []byte("1"))
	//print(resp)
	wallet, err := gateway.NewFileSystemWallet("wallet")
	if err != nil {
		log.Fatalf("Failed to create wallet: %v", err)
	}
	err = populateWallet(wallet)
	if !wallet.Exists("admin") {
		err = populateWallet(wallet)
		if err != nil {
			log.Fatalf("Failed to populate wallet contents: %v", err)
		}
	}
	gw, err := gateway.Connect(
		gateway.WithConfig(config.FromFile("C:\\Users\\xxx\\Desktop\\gm-sdk\\fabric-sdk-go-v1.0.0-gm\\main\\organizations\\peerOrganizations\\org1.xxzx.com\\connection-org1.yaml")),
		gateway.WithIdentity(wallet, "admin"),
	)
	if err != nil {
		log.Fatalf("Failed to connect to gateway: %v", err)
	}
	defer gw.Close()

	network, err := gw.GetNetwork("mychannel")
	if err != nil {
		log.Fatalf("Failed to get network: %v", err)
	}

	contract := network.GetContract("basic")
	resp ,_ :=contract.EvaluateTransaction("get","user","1")
	print(string(resp))
}
func queryCC(client *channel.Client, k1 []byte ,k2 []byte) string {
	var queryArgs = [][]byte{k1,k2}
	response, err := client.Query(channel.Request{
		ChaincodeID: cc,
		Fcn:         "get",
		Args:        queryArgs,
	})

	if err != nil {
		fmt.Println("Failed to query: ", err)
	}

	ret := string(response.Payload)
	fmt.Println("Chaincode status: ", response.ChaincodeStatus)
	fmt.Println("Payload: ", ret)
	return ret
}
func enrollUser(sdk *fabsdk.FabricSDK) {
	ctx := sdk.Context()
	mspClient, err := msp.New(ctx)
	if err != nil {
		fmt.Printf("Failed to create msp client: %s\n", err)
	}
	configBackend, err :=sdk.Config()
	cryptoSuiteConfig := cryptosuite.ConfigFromBackend(configBackend)
	keyStore :=cryptoSuiteConfig.KeyStorePath()
	_, err = mspClient.GetSigningIdentity(user)
	if err == msp.ErrUserNotFound {
		fmt.Println("Going to enroll user")
		userDta,err := mspClient.Enroll(user, msp.WithSecret(secret))
		dir, err := os.Getwd()
		keystr :=filepath.Join(dir,keyStore,userDta.KeyPath + "_sk")
		key , err:= ioutil.ReadFile(keystr)
		identity := gateway.NewX509Identity(userDta.MSPID, string(userDta.EnrollmentCertificate), string(key))
		wallet, err := gateway.NewFileSystemWallet("wallet")
		wallet.Put(user, identity)
		if err != nil {
			fmt.Printf("Failed to enroll user: %s\n", err)
		} else {
			fmt.Printf("Success enroll user: %s\n", user)
		}

	} else if err != nil {
		fmt.Printf("Failed to get user: %s\n", err)
	} else {
		fmt.Printf("User %s already enrolled, skip enrollment.\n", user)
	}
}

func registerUser(user string, secret string, sdk *fabsdk.FabricSDK) {

	ctxProvider := sdk.Context()

	// Get the Client.
	// Without WithOrg option, it uses default client organization.
	msp1, err := msp.New(ctxProvider)
	if err != nil {
		fmt.Printf("failed to create CA client: %s", err)
	}

	request := &msp.RegistrationRequest{Name: user, Secret: secret, Type: "client", Affiliation: "org1.department1"}
	_, err = msp1.Register(request)
	if err != nil {
		fmt.Printf("Register return error %s", err)
	}

}