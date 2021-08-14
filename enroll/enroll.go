package enroll

import (
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/core"
	pmsp "github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/keyvaluestore"
	"github.com/pkg/errors"
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
	user = "appuser"
	secret = "appuserpw"
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
func Register()  {
	user = "appuser"
	secret = "appuserpw"
	c := config.FromFile("D:\\go-sdk\\fabric-sdk-go-gm-master\\fabric-sdk-go-gm-master\\main\\config_test.yaml")
	sdk, err := fabsdk.New(c)
	if err != nil {
		fmt.Printf("Failed to create new SDK: %s\n", err)
		os.Exit(1)
	}
	defer sdk.Close()

	registerUser(user,secret,sdk)
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

type GatewayStore struct {
	store core.KVStore
}

func storeKeyFromUserIdentifier(key pmsp.IdentityIdentifier) string {
	return key.ID + "@" + key.MSPID + "-cert.pem"
}

// NewCertFileUserStore1 creates a new instance of CertFileUserStore
func NewCertFileUserStore1(store core.KVStore) (*GatewayStore, error) {
	return &GatewayStore{
		store: store,
	}, nil
}

// NewCertFileUserStore creates a new instance of CertFileUserStore
func NewCertFileUserStore(path string) (*GatewayStore, error) {
	if path == "" {
		return nil, errors.New("path is empty")
	}
	store, err := keyvaluestore.New(&keyvaluestore.FileKeyValueStoreOptions{
		Path: path,
	})
	if err != nil {
		return nil, errors.WithMessage(err, "user store creation failed")
	}
	return NewCertFileUserStore1(store)
}

// Load returns the User stored in the store for a key.
func (s *GatewayStore) Load(key pmsp.IdentityIdentifier) (*pmsp.UserData, error) {
	cert, err := s.store.Load(storeKeyFromUserIdentifier(key))
	if err != nil {
		if err == core.ErrKeyValueNotFound {
			return nil, msp.ErrUserNotFound
		}
		return nil, err
	}
	certBytes, ok := cert.([]byte)
	if !ok {
		return nil, errors.New("user is not of proper type")
	}
	userData := &pmsp.UserData{
		MSPID:                 key.MSPID,
		ID:                    key.ID,
		EnrollmentCertificate: certBytes,
	}
	return userData, nil
}

// Store stores a User into store
func (s *GatewayStore) Store(user *pmsp.UserData) error {
	key := storeKeyFromUserIdentifier(pmsp.IdentityIdentifier{MSPID: user.MSPID, ID: user.ID})
	return s.store.Store(key, user.EnrollmentCertificate)
}

// Delete deletes a User from store
func (s *GatewayStore) Delete(key pmsp.IdentityIdentifier) error {
	return s.store.Delete(storeKeyFromUserIdentifier(key))
}