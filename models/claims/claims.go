package claims

import (
	"fmt"
	//"os"
	"strconv"
	"time"

	"github.com/MyriadFlow/gateway/config/dbconfig"
	"github.com/MyriadFlow/gateway/config/envconfig"
	"github.com/MyriadFlow/gateway/models"
	"github.com/MyriadFlow/gateway/util/pkg/logwrapper"
)

type CustomClaims struct {
	WalletAddress string    `json:"walletAddress"`
	SignedBy      string    `json:"signedBy"`
	Expiration    time.Time `json:"expiryTime"`
}

type CustomClaimsWeb2User struct {
	Email      string    `json:"email"`
	SignedBy   string    `json:"signedBy"`
	Expiration time.Time `json:"expiryTime"`
}

func (c CustomClaims) Valid() error {
	db := dbconfig.GetDb()

	// if time.; err != nil {
	// 	return err
	// }
	err := db.Model(&models.User{}).Where("wallet_address = ?", c.WalletAddress).First(&models.User{}).Error
	if err != nil {
		return err
	}
	return nil
}

func New(walletAddress string) CustomClaims {
	pasetoExpirationInHours := envconfig.EnvVars.PASETO_EXPIRATION_IN_HOURS
	ok := false
	if pasetoExpirationInHours != "" {
		ok = true
	}
	pasetoExpirationInHoursInt := time.Duration(24)
	if ok {
		res, err := strconv.Atoi(pasetoExpirationInHours)
		if err != nil {
			logwrapper.Log.Warnf("Failed to parse PASETO_EXPIRATION_IN_HOURS as int : %v", err.Error())
		} else {
			pasetoExpirationInHoursInt = time.Duration(res)
		}
	}
	pasetoExpirationHours := pasetoExpirationInHoursInt * time.Hour
	expiration := time.Now().Add(pasetoExpirationHours)
	signedBy := envconfig.EnvVars.PASETO_SIGNED_BY
	return CustomClaims{
		walletAddress,
		signedBy,
		expiration,
	}
}

func NewUser(Email string) CustomClaimsWeb2User {
	pasetoExpirationInHours := envconfig.EnvVars.PASETO_EXPIRATION_IN_HOURS
	ok := false
	if pasetoExpirationInHours != "" {
		ok = true
	}
	pasetoExpirationInHoursInt := time.Duration(24)
	fmt.Println("ok value token", ok)
	if ok {
		res, err := strconv.Atoi(pasetoExpirationInHours)
		if err != nil {
			logwrapper.Log.Warnf("Failed to parse PASETO_EXPIRATION_IN_HOURS as int : %v", err.Error())
		} else {
			pasetoExpirationInHoursInt = time.Duration(res)
		}
	}
	pasetoExpirationHours := pasetoExpirationInHoursInt * time.Hour
	expiration := time.Now().Add(pasetoExpirationHours)
	signedBy := envconfig.EnvVars.PASETO_SIGNED_BY
	return CustomClaimsWeb2User{
		Email,
		signedBy,
		expiration,
	}
}
