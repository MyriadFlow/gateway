package testingcommon

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/MyriadFlow/gateway/api/types"

	"github.com/MyriadFlow/gateway/models/claims"
	"github.com/MyriadFlow/gateway/util/pkg/auth"

	"crypto/ecdsa"
	"log"

	"github.com/gin-gonic/gin"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/MyriadFlow/gateway/config/envconfig"
)

func PrepareAndGetAuthHeader(t *testing.T, testWalletAddress string) string {
	gin.SetMode(gin.TestMode)
	customClaims := claims.New(testWalletAddress)
	token, err := auth.GenerateTokenPaseto(customClaims)
	if err != nil {
		t.Fatal(err)
	}
	header := fmt.Sprintf("Bearer %v", token)

	return header
}

func GenerateWallet() *TestWallet {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal(err)
	}

	privateKeyBytes := crypto.FromECDSA(privateKey)
	pvtkey := hexutil.Encode(privateKeyBytes)[2:]
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	testWallet := TestWallet{
		PrivateKey:    pvtkey,
		WalletAddress: address,
	}
	return &testWallet
}

// Converts map created by json decoder to struct
// out should be pointer (&payload)
func ExtractPayload(response *types.ApiResponse, out interface{}) {
	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(response.Payload)
	json.NewDecoder(buf).Decode(out)
}

func InitializeEnvVars() {

	envconfig.EnvVars.APP_NAME = "storefront-gateway"
	envconfig.EnvVars.APP_PORT = 8000
	envconfig.EnvVars.APP_MODE = "debug"
	envconfig.EnvVars.APP_ALLOWED_ORIGIN = []string{"*"}

	envconfig.EnvVars.DB_HOST = "localhost"
	envconfig.EnvVars.DB_USERNAME = "postgres"
	envconfig.EnvVars.DB_PASSWORD = "postgres"
	envconfig.EnvVars.DB_NAME = "gateway"
	envconfig.EnvVars.DB_PORT = 5432

	envconfig.EnvVars.AUTH_EULA = "I Accept the MyriadFlow Terms of Service https://myriadflow.com/terms.html for accessing the application. Challenge ID: "

}
