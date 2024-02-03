package cookie

import (
	"encoding/hex"
	"errors"
	"github.com/paulmuenzner/powerplantmanager/config"
	logger "github.com/paulmuenzner/powerplantmanager/utils/logs"
	"net/http"
	"os"
)

func HasCookieExpired(r *http.Request, cookieName string) bool {

	authData, found := GetCookie(r, cookieName)
	if !found {
		return true
	}

	keyByteSlice, err := GetKeyByCookieName(cookieName)
	if err != nil {
		logger.GetLogger().Errorf("Cannot grab key for cookie with getKeyByCookieName in 'HasCookieExpired()' using 'GetKeyByCookieName()'. Cookie name: %s. Error: %v", cookieName, err)
		return true

	}

	_, err = VerifyAndExtractClaims(authData, keyByteSlice)
	if err != nil {
		logger.GetLogger().Errorf("Cannot verify cookie and extract claims in 'HasCookieExpired()' using 'VerifyAndExtractClaims()'. Cookie name: %s. Error: %v", cookieName, err)
		return true

	}

	return false
}

// ///////////////////////////////////////////
// Select key to unleash cookie claim

func GetKeyByCookieName(cookieName string) ([]byte, error) {

	////////// LIST ///////////////////////
	// Grab keys from .env file. Add more as needed
	jwtKeyAuthCookie, exists := os.LookupEnv("JWT_SECRET_KEY")
	if !exists {
		// Handle the case where jwtKey is not a string
		logger.GetLogger().Errorf("Cannot grab cookie key from .env in 'GetKeyByCookieName()' using 'os.LookupEnv()' for cookie name: %s", cookieName)
		return nil, errors.New("JwtKey is not a string")
	}

	//////////// LIST /////////////////////////////////////
	// Convert keys into type of []byte / byte slice
	jwtKeyAuthCookieSecret, err := hex.DecodeString(jwtKeyAuthCookie)
	if err != nil {
		// Handle the case where jwtKey is not a string
		logger.GetLogger().Errorf("Error converting string key into []slice type in 'GetKeyByCookieName()' using 'DecodeString()'. Error: %v", err)
		return nil, err
	}

	switch cookieName {
	case config.AuthCookieName:
		return jwtKeyAuthCookieSecret, nil
	default:
		return nil, errors.New("no valid key found in 'GetKeyByCookieName()'")
	}
}
