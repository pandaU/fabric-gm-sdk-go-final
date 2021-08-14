package enroll

import (
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"os"
)
var (
	cc            = ""
	user          = ""
	secret        = ""
	channelName   = ""
	chaincodrPath = "github.com/VoneChain-CS/fabric-gm/scripts/fabric-samples/chaincode/abstore/go"
)
func EnrollUser()  {
	user = "admin"
	secret = "adminpw"
	channelName = "mychannel"
	cc = "mycc_3"
	fmt.Println("Reading connection profile..")
	c := config.FromFile("D:\\go-sdk\\fabric-sdk-go-gm-master\\fabric-sdk-go-gm-master\\main\\config_test.yaml")
	sdk, err := fabsdk.New(c)
	if err != nil {
		fmt.Printf("Failed to create new SDK: %s\n", err)
		os.Exit(1)
	}
	defer sdk.Close()

	//registerUser(user,secret,sdk)
	enrollUser(sdk)
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