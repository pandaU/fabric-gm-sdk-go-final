package enroll

import (
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/cryptosuite"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
	"io/ioutil"
	"os"
	"path/filepath"
)
var (
	cc            = ""
	user          = ""
	secret        = ""
)
func EnrollUser(user string,secret string,configPath string,walletPath string)  {
	c := config.FromFile(configPath)
	sdk, err := fabsdk.New(c)
	if err != nil {
		fmt.Printf("Failed to create new SDK: %s\n", err)
		os.Exit(1)
	}
	defer sdk.Close()

	//registerUser(user,secret,sdk)
	enrollUser(user,secret,sdk,walletPath)
}
func Register(user string,secret string,configPath string)  {
	c := config.FromFile(configPath)
	sdk, err := fabsdk.New(c)
	if err != nil {
		fmt.Printf("Failed to create new SDK: %s\n", err)
		os.Exit(1)
	}
	defer sdk.Close()

	registerUser(user,secret,sdk)
}
func enrollUser(user string,secret string,sdk *fabsdk.FabricSDK,walletPath string) {
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
		if err != nil {
			fmt.Printf("Failed to enroll user: %s\n", err)
		} else {
			fmt.Printf("Success enroll user: %s\n", user)
		}
		dir, err := os.Getwd()
		keystr :=filepath.Join(dir,keyStore,userDta.KeyPath + "_sk")
		key , err:= ioutil.ReadFile(keystr)
		identity := gateway.NewX509Identity(userDta.MSPID, string(userDta.EnrollmentCertificate), string(key))
		wallet, err := gateway.NewFileSystemWallet(walletPath)
		wallet.Put(user, identity)
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
