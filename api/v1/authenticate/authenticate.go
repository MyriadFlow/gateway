package authenticate

import (
	"net/http"

	"github.com/MyriadFlow/gateway/models/claims"

	"github.com/MyriadFlow/gateway/config/dbconfig"
	"github.com/MyriadFlow/gateway/config/envconfig"
	"github.com/MyriadFlow/gateway/models"
	"github.com/MyriadFlow/gateway/util/pkg/auth"
	"github.com/MyriadFlow/gateway/util/pkg/cryptosign"
	"github.com/MyriadFlow/gateway/util/pkg/flowid"
	"github.com/MyriadFlow/gateway/util/pkg/httphelper"
	"github.com/MyriadFlow/gateway/util/pkg/logwrapper"
	"github.com/ethereum/go-ethereum/common/hexutil"
	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

// ApplyRoutes applies router to gin Router
func ApplyRoutes(r *gin.RouterGroup) {
	g := r.Group("/auth")
	{
		g.GET("/web3", GetFlowId)
		g.POST("/web3", authenticate)
	}
}

func authenticate(c *gin.Context) {

	db := dbconfig.GetDb()
	var req AuthenticateRequest
	err := c.BindJSON(&req)
	if err != nil {
		httphelper.ErrResponse(c, http.StatusForbidden, "payload is invalid")
		return
	}
	//Get flowid type
	var flowIdData models.FlowId
	err = db.Model(&models.FlowId{}).Where("flow_id = ?", req.FlowId).First(&flowIdData).Error
	if err != nil {
		logwrapper.Errorf("failed to get flowId, error %v", err)
		httphelper.ErrResponse(c, http.StatusNotFound, "flow id not found")
		return
	}

	if flowIdData.FlowIdType != models.AUTH {
		httphelper.ErrResponse(c, http.StatusBadRequest, "Flow id not created for auth")
		return
	}

	if err != nil {
		logwrapper.Error(err)
		httphelper.ErrResponse(c, 500, "Unexpected error occured")
		return
	}
	userAuthEULA := envconfig.EnvVars.AUTH_EULA
	message := userAuthEULA + req.FlowId
	walletAddress, isCorrect, err := cryptosign.CheckSign(req.Signature, req.FlowId, message)

	if err == cryptosign.ErrFlowIdNotFound {
		httphelper.ErrResponse(c, http.StatusNotFound, "Flow Id not found")
		return
	}

	if err != nil {
		logwrapper.Errorf("failed to CheckSignature, error %v", err.Error())
		httphelper.ErrResponse(c, http.StatusInternalServerError, "Unexpected error occured")
		return
	}
	if isCorrect {
		customClaims := claims.New(walletAddress)
		pasetoToken, err := auth.GenerateTokenPaseto(customClaims)
		if err != nil {
			httphelper.NewInternalServerError(c, "failed to generate token, error %v", err.Error())
			return
		}

		payload := AuthenticatePayload{
			Token: pasetoToken,
		}
		httphelper.SuccessResponse(c, "Token generated successfully", payload)
	} else {
		httphelper.ErrResponse(c, http.StatusForbidden, "Wallet Address is not correct")
		return
	}
}

func GetFlowId(c *gin.Context) {
	walletAddress := c.Query("walletAddress")
	chain := c.Query("chain")

	// Validate required parameters
	if walletAddress == "" {
		httphelper.ErrResponse(c, http.StatusBadRequest, "Wallet address (walletAddress) is required")
		return
	}
	if chain == "" {
		httphelper.ErrResponse(c, http.StatusBadRequest, "Chain (chain) is required")
		return
	}

	// Handle chain-specific validation
	switch chain {
	case "ethereum":
		// Validate Ethereum wallet address
		_, err := hexutil.Decode(walletAddress)
		if err != nil {
			httphelper.ErrResponse(c, http.StatusBadRequest, "Wallet address (walletAddress) is not valid for Ethereum")
			return
		}
	case "solana":
		// No validation needed for Solana; directly proceed
	default:
		httphelper.ErrResponse(c, http.StatusBadRequest, "Unsupported chain. Valid values are 'ethereum' or 'solana'")
		return
	}

	// Generate Flow ID
	flowId, err := flowid.GenerateFlowId(walletAddress, models.AUTH, "")
	if err != nil {
		log.Error(err)
		httphelper.ErrResponse(c, http.StatusInternalServerError, "Unexpected error occurred")
		return
	}

	// Build and send response
	userAuthEULA := envconfig.EnvVars.AUTH_EULA
	payload := GetFlowIdPayload{
		FlowId: flowId,
		Eula:   userAuthEULA,
	}
	httphelper.SuccessResponse(c, "Flow ID successfully generated", payload)
}
