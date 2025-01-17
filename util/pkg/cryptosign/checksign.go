package cryptosign

import (
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/MyriadFlow/gateway/config/dbconfig"
	"github.com/MyriadFlow/gateway/models"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/mr-tron/base58"
	"golang.org/x/crypto/ed25519"
)

var (
	ErrFlowIdNotFound = errors.New("flow id not found")
)

func CheckSign(signature string, flowId string, message string) (string, bool, error) {
	// Ethereum verification first
	walletAddress, isCorrect, err := verifyEthereumSignature(signature, flowId, message)
	if err == nil {
		return walletAddress, isCorrect, nil
	}
	log.Printf("Ethereum verification failed: %v", err)

	// Fall back to Solana verification
	return verifySolanaSignature(signature, flowId, message)
}

func verifyEthereumSignature(signature string, flowId string, message string) (string, bool, error) {
	// Construct Ethereum Signed Message format
	db := dbconfig.GetDb()

	newMsg := fmt.Sprintf("\x19Ethereum Signed Message:\n%v%v", len(message), message)
	newMsgHash := crypto.Keccak256Hash([]byte(newMsg))

	// Decode signature
	signatureInBytes, err := hexutil.Decode(signature)
	if err != nil {
		return "", false, fmt.Errorf("failed to decode Ethereum signature: %v", err)
	}

	// Normalize Ethereum signature
	if signatureInBytes[64] == 27 || signatureInBytes[64] == 28 {
		signatureInBytes[64] -= 27
	}

	// Recover public key from Ethereum signature
	pubKey, err := crypto.SigToPub(newMsgHash.Bytes(), signatureInBytes)
	if err != nil {
		return "", false, fmt.Errorf("failed to recover public key from Ethereum signature: %v", err)
	}

	// Get wallet address from the public key
	walletAddress := crypto.PubkeyToAddress(*pubKey)

	// Fetch flow ID data from the database
	var flowIdData models.FlowId
	res := db.Model(&models.FlowId{}).Where("flow_id = ?", flowId).First(&flowIdData)
	if res.RowsAffected == 0 {
		return "", false, ErrFlowIdNotFound
	}
	if err := res.Error; err != nil {
		return "", false, err
	}

	// Compare wallet address
	if strings.EqualFold(flowIdData.WalletAddress, walletAddress.String()) {
		return flowIdData.WalletAddress, true, nil
	}

	return "", false, nil
}

func verifySolanaSignature(signature string, flowId string, message string) (string, bool, error) {

	db := dbconfig.GetDb()

	// Decode the hex signature
	signatureBytes, err := hex.DecodeString(signature)
	if err != nil {
		return "", false, fmt.Errorf("failed to decode Solana signature: %v", err)
	}

	// Fetch flow ID data from the database
	var flowIdData models.FlowId
	res := db.Model(&models.FlowId{}).Where("flow_id = ?", flowId).First(&flowIdData)
	if res.RowsAffected == 0 {
		return "", false, ErrFlowIdNotFound
	}
	if err := res.Error; err != nil {
		return "", false, err
	}

	// Decode the public key from the wallet address
	pubKeyBytes, err := base58.Decode(flowIdData.WalletAddress)
	if err != nil {
		return "", false, fmt.Errorf("failed to decode wallet address: %v", err)
	}

	// Verify the signature
	if ed25519.Verify(pubKeyBytes, []byte(message), signatureBytes) {
		return flowIdData.WalletAddress, true, nil
	}

	return "", false, nil
}
