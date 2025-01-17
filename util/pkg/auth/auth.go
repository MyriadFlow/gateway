package auth

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	gopaseto "aidanwoods.dev/go-paseto"
	"github.com/MyriadFlow/gateway/config/envconfig"
	"github.com/MyriadFlow/gateway/models/claims"
	"github.com/MyriadFlow/gateway/util/pkg/logwrapper"
)

var PublicKey gopaseto.V4AsymmetricPublicKey
var secretKey gopaseto.V4AsymmetricSecretKey
var initialized bool

func Init() {
	secretKey = gopaseto.NewV4AsymmetricSecretKey()
	PublicKey = secretKey.Public()
	initialized = true

}

func GenerateTokenPaseto(claim claims.CustomClaims) (string, error) {
	if !initialized {
		secretKey = gopaseto.NewV4AsymmetricSecretKey()
		PublicKey = secretKey.Public()
	}

	footer := envconfig.EnvVars.PASETO_FOOTER
	claimbyte, _ := json.Marshal(claim)
	token, err := gopaseto.NewTokenFromClaimsJSON(claimbyte, []byte(footer))
	if err != nil {
		return "", err
	}

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
	token.SetExpiration(expiration)
	signed := token.V4Sign(secretKey, nil)
	return signed, nil
}

func GenerateTokenPasetoWeb2User(claim claims.CustomClaimsWeb2User) (string, error) {
	if !initialized {
		secretKey = gopaseto.NewV4AsymmetricSecretKey()
		PublicKey = secretKey.Public()
	}

	footer := envconfig.EnvVars.PASETO_FOOTER
	claimbyte, _ := json.Marshal(claim)
	token, err := gopaseto.NewTokenFromClaimsJSON(claimbyte, []byte(footer))
	if err != nil {
		return "", err
	}

	pasetoExpirationInHours := envconfig.EnvVars.PASETO_EXPIRATION_IN_HOURS
	ok := false
	if pasetoExpirationInHours != "" {
		ok = true
	}
	fmt.Println("ok value token", ok)
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
	token.SetExpiration(expiration)
	signed := token.V4Sign(secretKey, nil)
	return signed, nil
}

func Getpublickey() gopaseto.V4AsymmetricPublicKey {
	publickey := PublicKey
	return publickey
}

func GetTestpublickey() gopaseto.V4AsymmetricPublicKey {
	publickey := PublicKey
	return publickey
}

func GenerateExpiredTokenPaseto(claim claims.CustomClaims) (string, error) {
	if !initialized {
		secretKey = gopaseto.NewV4AsymmetricSecretKey()
		PublicKey = secretKey.Public()
	}

	footer := envconfig.EnvVars.PASETO_FOOTER
	claimbyte, _ := json.Marshal(claim)
	token, err := gopaseto.NewTokenFromClaimsJSON(claimbyte, []byte(footer))
	if err != nil {
		return "", err
	}
	pasetoExpirationInHoursInt := time.Duration(-1)
	pasetoExpirationHours := pasetoExpirationInHoursInt * time.Hour
	expiration := time.Now().Add(pasetoExpirationHours)
	token.SetExpiration(expiration)
	signed := token.V4Sign(secretKey, nil)
	return signed, nil
}
