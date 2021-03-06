package main

import (
	"fmt"
	"github.com/VoneChain-CS/fabric-sdk-go-gm/enroll"
	"github.com/VoneChain-CS/fabric-sdk-go-gm/pkg/client/channel"
	"github.com/VoneChain-CS/fabric-sdk-go-gm/pkg/client/msp"
	"github.com/VoneChain-CS/fabric-sdk-go-gm/pkg/core/config"
	"github.com/VoneChain-CS/fabric-sdk-go-gm/pkg/fabsdk"
	"github.com/VoneChain-CS/fabric-sdk-go-gm/pkg/gateway"
	"io/ioutil"
	"log"
	//"os"
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
	credPath := "D:\\go-sdk\\fabric-sdk-go-gm-master\\fabric-sdk-go-gm-master"

	certPath := filepath.Join(credPath, "wallet", "appuser@Org1MSP-cert.pem")
	// read the certificate pem
	cert, err := ioutil.ReadFile(filepath.Clean(certPath))
	if err != nil {
		return err
	}

	keyDir := filepath.Join(credPath, "keystore", "keystore")
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

	return wallet.Put("appuser", identity)
}
func main() {
	user = "pandau"
	secret = "pandaupw"
	channelName = "mychannel"

	walletPath := "wallet"
	wallet, err := gateway.NewFileSystemWallet(walletPath)
	if err != nil {
		log.Fatalf("Failed to create wallet: %v", err)
	}
	//err = populateWallet(wallet)
	configPath :="D:\\go-sdk\\fabric-sdk-go-gm-master\\fabric-sdk-go-gm-master\\main\\config_test.yaml"
	if !wallet.Exists(user) {
		enroll.Register(user,secret,configPath)
		enroll.EnrollUser(user,secret,configPath,walletPath)
	}
	gw, err := gateway.Connect(
		gateway.WithConfig(config.FromFile(configPath)),
		gateway.WithIdentity(wallet, user),
	)
	if err != nil {
		log.Fatalf("Failed to connect to gateway: %v", err)
	}
	defer gw.Close()

	network, err := gw.GetNetwork(channelName)
	if err != nil {
		log.Fatalf("Failed to get network: %v", err)
	}
	contract := network.GetContract("basic")
	contract.SubmitTransaction("create","user","6","?????????")
	resp ,_ :=contract.EvaluateTransaction("get","user","6")
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

	_, err = mspClient.GetSigningIdentity(user)
	if err == msp.ErrUserNotFound {
		fmt.Println("Going to enroll user")
		err = mspClient.Enroll(user, msp.WithSecret(secret))

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