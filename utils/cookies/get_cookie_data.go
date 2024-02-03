package cookie

import (
	"errors"
	logger "github.com/paulmuenzner/powerplantmanager/utils/logs"
	"net/http"

	"github.com/dgrijalva/jwt-go"
)

func GetCookieData(r *http.Request, cookieName string) (jwt.MapClaims, error) {

	authData, found := GetCookie(r, cookieName)
	if !found {
		return nil, errors.New("")
	}

	keyByteSlice, err := GetKeyByCookieName(cookieName)
	if err != nil {
		logger.GetLogger().Error("Cannot grab key for cookie with getKeyByCookieName in 'HasCookieExpired'. Error: ", err, "Cookie name: ", cookieName)
		return nil, err

	}

	claim, err := VerifyAndExtractClaims(authData, keyByteSlice)
	if err != nil {
		logger.GetLogger().Error("Cannot verify cookie and extract claims in 'HasCookieExpired'. Error: ", err, "Cookie name: ", cookieName)
		return nil, err

	}

	return claim, nil
}
