package authenticate

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/MyriadFlow/gateway/api/types"
	"github.com/MyriadFlow/gateway/util/pkg/logwrapper"
	testingcommon "github.com/MyriadFlow/gateway/util/testingcommon"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

func Test_GetFlowId(t *testing.T) {
	testingcommon.InitializeEnvVars()
	logwrapper.Init()
	t.Cleanup(testingcommon.DeleteCreatedEntities())
	gin.SetMode(gin.TestMode)
	testWalletAddress := testingcommon.GenerateWallet().WalletAddress
	u, err := url.Parse("/api/v1.0/auth/web3")
	if err != nil {
		t.Fatal(err)
	}
	t.Run("Should fail if wallet address is not hexadecimal", func(t *testing.T) {
		q := url.Values{}
		q.Set("walletAddress", "invalidwalletaddr")
		u.RawQuery = q.Encode()
		rr := httptest.NewRecorder()

		req, err := http.NewRequest("GET", u.String(), nil)
		if err != nil {
			t.Error(err)
		}
		c, _ := gin.CreateTestContext(rr)
		c.Request = req
		GetFlowId(c)
		assert.Equal(t, http.StatusBadRequest, rr.Result().StatusCode)
	})
	t.Run("Should be able to get flow id", func(t *testing.T) {

		q := url.Values{}
		q.Set("walletAddress", testWalletAddress)
		u.RawQuery = q.Encode()
		rr := httptest.NewRecorder()

		req, err := http.NewRequest("GET", u.String(), nil)
		if err != nil {
			t.Error(err)
		}
		c, _ := gin.CreateTestContext(rr)
		c.Request = req
		GetFlowId(c)
		assert.Equal(t, http.StatusOK, rr.Result().StatusCode)
	})

}

// TODO: Write test to verify expiry
func Test_PostAuthenticate(t *testing.T) {

	testingcommon.InitializeEnvVars()
	logwrapper.Init()

	t.Cleanup(testingcommon.DeleteCreatedEntities())
	gin.SetMode(gin.TestMode)

	url := "/api/v1.0/auth/web3"

	t.Run("Should return 403 with different wallet address", func(t *testing.T) {
		testWallet := testingcommon.GenerateWallet()
		eula, flowId := callFlowIdApi(testWallet.WalletAddress, t)
		// Different private key will result in different wallet address
		differentPrivatekey := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
		signature := getSignature(eula, flowId, differentPrivatekey)
		body := AuthenticateRequest{Signature: signature, FlowId: flowId}
		jsonBody, err := json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}
		newWalletAddress := testWallet.WalletAddress + "ba"
		callFlowIdApi(newWalletAddress, t)

		rr := httptest.NewRecorder()

		//Request with signature stil created from different walletAddress
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
		if err != nil {
			t.Fatal(err)
		}
		c, _ := gin.CreateTestContext(rr)
		c.Request = req
		authenticate(c)
		assert.Equal(t, http.StatusForbidden, rr.Code, rr.Body.String())
	})

	t.Run("Should return 200 with correct wallet address", func(t *testing.T) {
		testWallet := testingcommon.GenerateWallet()
		eula, flowId := callFlowIdApi(testWallet.WalletAddress, t)

		signature := getSignature(eula, flowId, testWallet.PrivateKey)
		body := AuthenticateRequest{Signature: signature, FlowId: flowId}
		jsonBody, err := json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()

		//Request with signature created from correct wallet address
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
		if err != nil {
			t.Fatal(err)
		}

		c, _ := gin.CreateTestContext(rr)
		c.Request = req
		authenticate(c)
		assert.Equal(t, http.StatusOK, rr.Code, rr.Body.String())
	})

}

func callFlowIdApi(walletAddress string, t *testing.T) (eula string, flowidString string) {
	// Call flowid api
	u, err := url.Parse("/api/v1.0/auth/web3")
	q := url.Values{}
	q.Set("walletAddress", walletAddress)
	u.RawQuery = q.Encode()
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	req, err := http.NewRequest("GET", u.String(), nil)
	req.URL.RawQuery = q.Encode()
	if err != nil {
		t.Error(err)
	}
	c, _ := gin.CreateTestContext(rr)
	c.Request = req
	GetFlowId(c)
	assert.Equal(t, http.StatusOK, rr.Code, "Failed to call flowApi")
	var flowIdPayload GetFlowIdPayload
	var res types.ApiResponse
	decoder := json.NewDecoder(rr.Result().Body)
	err = decoder.Decode(&res)
	testingcommon.ExtractPayload(&res, &flowIdPayload)
	if err != nil {
		t.Fatal(err)
	}
	return flowIdPayload.Eula, flowIdPayload.FlowId
}

func getSignature(eula string, flowId string, hexPrivateKey string) string {
	message := eula + flowId
	newMsg := fmt.Sprintf("\x19Ethereum Signed Message:\n%v%v", len(message), message)
	privateKey, err := crypto.HexToECDSA(hexPrivateKey)
	if err != nil {
		log.Fatal("HexToECDSA failed ", err)
	}

	// keccak256 hash of the data
	dataBytes := []byte(newMsg)
	hashData := crypto.Keccak256Hash(dataBytes)

	signatureBytes, err := crypto.Sign(hashData.Bytes(), privateKey)
	if err != nil {
		log.Fatal("len", err)
	}

	signature := hexutil.Encode(signatureBytes)

	return signature
}

// func Test_Web2Auth(t *testing.T) {
// 	gin.SetMode(gin.TestMode)

// 	t.Run("Should return 403 with invalid payload", func(t *testing.T) {
// 		// Prepare the request payload
// 		payload := web2AuthRequest{
// 			// Provide an invalid payload here
// 		}
// 		jsonBody, err := json.Marshal(payload)
// 		if err != nil {
// 			t.Fatal(err)
// 		}

// 		rr := httptest.NewRecorder()

// 		// Create a test request
// 		req, err := http.NewRequest("POST", "/api/v1.0/auth/web2", bytes.NewBuffer(jsonBody))
// 		if err != nil {
// 			t.Fatal(err)
// 		}

// 		// Create a Gin context with the test request
// 		c, _ := gin.CreateTestContext(rr)
// 		c.Request = req

// 		// Call the Web2Auth handler function
// 		Web2Auth(c)

// 		// Assert the response status code and body
// 		assert.Equal(t, http.StatusForbidden, rr.Code, "Unexpected response status code")
// 		expectedResponse := ErrorResponse{
// 			Message: "payload is invalid",
// 		}
// 		expectedResponseBody, err := json.Marshal(expectedResponse)
// 		if err != nil {
// 			t.Fatal(err)
// 		}
// 		assert.Equal(t, string(expectedResponseBody), rr.Body.String(), "Unexpected response body")
// 	})

// 	t.Run("Should return 403 with Supabase error", func(t *testing.T) {
// 		// Prepare the request payload
// 		payload := web2AuthRequest{
// 			Provider: "supabase",
// 			Token:    "invalid_token",
// 		}
// 		jsonBody, err := json.Marshal(payload)
// 		if err != nil {
// 			t.Fatal(err)
// 		}

// 		rr := httptest.NewRecorder()

// 		// Create a test request
// 		req, err := http.NewRequest("POST", "/api/v1.0/auth/web2", bytes.NewBuffer(jsonBody))
// 		if err != nil {
// 			t.Fatal(err)
// 		}

// 		// Create a Gin context with the test request
// 		c, _ := gin.CreateTestContext(rr)
// 		c.Request = req

// 		// Call the Web2Auth handler function
// 		Web2Auth(c)

// 		// Assert the response status code and body
// 		assert.Equal(t, http.StatusForbidden, rr.Code, "Unexpected response status code")
// 		expectedResponse := ErrorResponse{
// 			Message: "Supabase error message",
// 		}
// 		expectedResponseBody, err := json.Marshal(expectedResponse)
// 		if err != nil {
// 			t.Fatal(err)
// 		}
// 		assert.Equal(t, string(expectedResponseBody), rr.Body.String(), "Unexpected response body")
// 	})

// 	t.Run("Should return 404 when user is not found", func(t *testing.T) {
// 		// Prepare the request payload
// 		payload := web2AuthRequest{
// 			Provider: "supabase",
// 			Token:    "valid_token",
// 		}
// 		jsonBody, err := json.Marshal(payload)
// 		if err != nil {
// 			t.Fatal(err)
// 		}

// 		rr := httptest.NewRecorder()

// 		// Create a test request
// 		req, err := http.NewRequest("POST", "/api/v1.0/auth/web2", bytes.NewBuffer(jsonBody))
// 		if err != nil {
// 			t.Fatal(err)
// 		}

// 		// Create a Gin context with the test request
// 		c, _ := gin.CreateTestContext(rr)
// 		c.Request = req

// 		// Call the Web2Auth handler function
// 		Web2Auth(c)

// 		// Assert the response status code and body
// 		assert.Equal(t, http.StatusNotFound, rr.Code, "Unexpected response status code")
// 		expectedResponse := ErrorResponse{
// 			Message: "Web2Auth failed to count user",
// 		}
// 		expectedResponseBody, err := json.Marshal(expectedResponse)
// 		if err != nil {
// 			t.Fatal(err)
// 		}
// 		assert.Equal(t, string(expectedResponseBody), rr.Body.String(), "Unexpected response body")
// 	})

// 	// Add more test cases as needed for different scenarios

// }

// func Test_Web2Auth(t *testing.T) {
// 	testingcommon.InitializeEnvVars()
// 	logwrapper.Init()
// 	t.Cleanup(testingcommon.DeleteCreatedEntities())
// 	gin.SetMode(gin.TestMode)

// 	t.Run("Should return 403 with invalid payload", func(t *testing.T) {
// 		// Prepare the request payload
// 		payload := web2AuthRequest{
// 			// Provide an invalid payload here
// 		}
// 		jsonBody, err := json.Marshal(payload)
// 		if err != nil {
// 			t.Fatal(err)
// 		}

// 		rr := httptest.NewRecorder()

// 		// Create a test request
// 		req, err := http.NewRequest("POST", "/api/v1.0/auth/web2", bytes.NewBuffer(jsonBody))
// 		if err != nil {
// 			t.Fatal(err)
// 		}

// 		// Create a Gin context with the test request
// 		c, _ := gin.CreateTestContext(rr)
// 		c.Request = req

// 		// Call the Web2Auth handler function
// 		Web2Auth(c)

// 		// Assert the response status code and body
// 		assert.Equal(t, http.StatusForbidden, rr.Code, "Unexpected response status code")
// 		expectedResponse := ErrorResponse{
// 			Message: "payload is invalid",
// 		}
// 		expectedResponseBody, err := json.Marshal(expectedResponse)
// 		if err != nil {
// 			t.Fatal(err)
// 		}
// 		assert.Equal(t, string(expectedResponseBody), rr.Body.String(), "Unexpected response body")
// 	})

// 	t.Run("Should return 403 with Supabase error", func(t *testing.T) {
// 		// Prepare the request payload
// 		payload := web2AuthRequest{
// 			Token:     "invalid_token",
// 			Provider:  "supabase",
// 			User_Type: "web2",
// 		}
// 		jsonBody, err := json.Marshal(payload)
// 		if err != nil {
// 			t.Fatal(err)
// 		}

// 		rr := httptest.NewRecorder()

// 		// Create a test request
// 		req, err := http.NewRequest("POST", "/api/v1.0/auth/web2", bytes.NewBuffer(jsonBody))
// 		if err != nil {
// 			t.Fatal(err)
// 		}

// 		// Create a Gin context with the test request
// 		c, _ := gin.CreateTestContext(rr)
// 		c.Request = req

// 		// Call the Web2Auth handler function
// 		Web2Auth(c)

// 		// Assert the response status code and body
// 		assert.Equal(t, http.StatusForbidden, rr.Code, "Unexpected response status code")
// 		expectedResponse := ErrorResponse{
// 			Message: "Supabase error message",
// 		}
// 		expectedResponseBody, err := json.Marshal(expectedResponse)
// 		if err != nil {
// 			t.Fatal(err)
// 		}
// 		assert.Equal(t, string(expectedResponseBody), rr.Body.String(), "Unexpected response body")
// 	})

// 	t.Run("Should return 404 when user is not found", func(t *testing.T) {
// 		// Prepare the request payload
// 		payload := web2AuthRequest{
// 			Token:     "valid_token",
// 			Provider:  "supabase",
// 			User_Type: "web2",
// 		}
// 		jsonBody, err := json.Marshal(payload)
// 		if err != nil {
// 			t.Fatal(err)
// 		}

// 		rr := httptest.NewRecorder()

// 		// Create a test request
// 		req, err := http.NewRequest("POST", "/api/v1.0/auth/web2", bytes.NewBuffer(jsonBody))
// 		if err != nil {
// 			t.Fatal(err)
// 		}

// 		// Create a Gin context with the test request
// 		c, _ := gin.CreateTestContext(rr)
// 		c.Request = req

// 		// Call the Web2Auth handler function
// 		Web2Auth(c)

// 		// Assert the response status code and body
// 		assert.Equal(t, http.StatusNotFound, rr.Code, "Unexpected response status code")
// 		expectedResponse := ErrorResponse{
// 			Message: "Web2Auth failed to count user",
// 		}
// 		expectedResponseBody, err := json.Marshal(expectedResponse)
// 		if err != nil {
// 			t.Fatal(err)
// 		}
// 		assert.Equal(t, string(expectedResponseBody), rr.Body.String(), "Unexpected response body")
// 	})

// 	// Add more test cases as needed for different scenarios

// }
