package paseto

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	gopaseto "aidanwoods.dev/go-paseto"

	"github.com/MyriadFlow/gateway/models/claims"
	"github.com/MyriadFlow/gateway/util/pkg/auth"
	"github.com/MyriadFlow/gateway/util/pkg/httphelper"
	"github.com/MyriadFlow/gateway/util/pkg/logwrapper"

	"github.com/gin-gonic/gin"
)

var (
	ErrAuthHeaderMissing = errors.New("authorization header is required")
)

func PASETO(c *gin.Context) {
	var headers GenericAuthHeaders
	err := c.BindHeader(&headers)
	if err != nil {
		err = fmt.Errorf("failed to bind header, %s", err)
		logValidationFailed(headers.Authorization, err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if headers.Authorization == "" {
		logValidationFailed(headers.Authorization, ErrAuthHeaderMissing)
		httphelper.ErrResponse(c, http.StatusBadRequest, ErrAuthHeaderMissing.Error())
		c.Abort()
		return
	}
	token := headers.Authorization
	splitToken := strings.Split(token, "Bearer ")
	pasetoToken := splitToken[1]
	parser := gopaseto.NewParser()
	parser.AddRule(gopaseto.NotExpired())
	publickey := auth.Getpublickey()
	parsedToken, err := parser.ParseV4Public(publickey, pasetoToken, nil)
	if err != nil {
		err = fmt.Errorf("failed to parsed paseto token, %s", err)
		logValidationFailed(headers.Authorization, err)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	jsonvalue := parsedToken.ClaimsJSON()
	ClaimsValue := claims.CustomClaims{}
	err = json.Unmarshal(jsonvalue, &ClaimsValue)
	if err != nil {
		err = fmt.Errorf("failed to scan claims for paseto token, %s", err)
		logValidationFailed(headers.Authorization, err)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	} else {
		c.Set("walletAddress", ClaimsValue.WalletAddress)
		c.Next()
	}

}

func logValidationFailed(token string, err error) {
	logwrapper.Warnf("validation failed with token %v and error: %v", token, err)
}
