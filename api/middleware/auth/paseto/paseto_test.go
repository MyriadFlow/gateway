package paseto

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MyriadFlow/gateway/config/dbconfig"
	"github.com/MyriadFlow/gateway/config/envconfig"
	"github.com/MyriadFlow/gateway/models"
	"github.com/MyriadFlow/gateway/models/claims"
	"github.com/MyriadFlow/gateway/util/pkg/auth"
	"github.com/MyriadFlow/gateway/util/pkg/logwrapper"
	"github.com/MyriadFlow/gateway/util/testingcommon"

	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func Test_PASETO(t *testing.T) {
	testingcommon.InitializeEnvVars()
	logwrapper.Init()
	db := dbconfig.GetDb()
	t.Cleanup(testingcommon.DeleteCreatedEntities())
	gin.SetMode(gin.TestMode)
	testWalletAddress := testingcommon.GenerateWallet().WalletAddress
	newUser := models.User{
		WalletAddress: testWalletAddress,
	}
	err := db.Model(&models.User{}).Create(&newUser).Error
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		db.Delete(&newUser)
	}()

	t.Run("Should return 200 with correct PASETO", func(t *testing.T) {
		newClaims := claims.New(testWalletAddress)
		token, err := auth.GenerateTokenPaseto(newClaims)
		if err != nil {
			t.Fatal(err)
		}
		token = "Bearer " + token
		rr := callApi(t, token)
		assert.Equal(t, http.StatusOK, rr.Result().StatusCode)
	})

	t.Run("Should return 401 and 4011 with expired PASETO", func(t *testing.T) {
		signedBy := envconfig.EnvVars.PASETO_SIGNED_BY
		newClaims := claims.CustomClaims{
			WalletAddress: testWalletAddress,
			SignedBy:      signedBy,
		}
		fmt.Println("newClaims : ", newClaims)
		token, err := auth.GenerateExpiredTokenPaseto(newClaims)
		if err != nil {
			t.Fatal(err)
		}
		token = "Bearer " + token
		rr := callApi(t, token)
		assert.Equal(t, http.StatusUnauthorized, rr.Result().StatusCode)
	})

}

func callApi(t *testing.T, token string) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	ginTestApp := gin.New()

	rq, err := http.NewRequest("POST", "", nil)
	if err != nil {
		t.Fatal(err)
	}
	rq.Header.Add("Authorization", token)
	ginTestApp.Use(PASETO)
	ginTestApp.Use(successHander)
	ginTestApp.ServeHTTP(rr, rq)
	return rr
}

func successHander(c *gin.Context) {
	c.Status(http.StatusOK)
}
