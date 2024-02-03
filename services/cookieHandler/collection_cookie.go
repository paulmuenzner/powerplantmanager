package cookiehandler

import (
	"encoding/hex"
	"errors"
	"github.com/paulmuenzner/powerplantmanager/config"
	errHandler "github.com/paulmuenzner/powerplantmanager/services/errorHandler"
	cookie "github.com/paulmuenzner/powerplantmanager/utils/cookies"
	logger "github.com/paulmuenzner/powerplantmanager/utils/logs"
	"net/http"
	"time"
)

/////////////////////////////////////////////////////////////////////
// SET AUTH COOKIE //////////////////////////////////////////////////

type SameSiteOption string

const (
	SameSiteDefault SameSiteOption = "SameSite"
	SameSiteStrict  SameSiteOption = "Strict"
	SameSiteLax     SameSiteOption = "Lax"
	SameSiteNone    SameSiteOption = "None"
)

// convertSameSiteOption converts a SameSiteOption to http.SameSite.
func convertSameSiteOption(sameSite SameSiteOption) http.SameSite {
	switch sameSite {
	case SameSiteDefault:
		return http.SameSiteDefaultMode
	case SameSiteStrict:
		return http.SameSiteStrictMode
	case SameSiteLax:
		return http.SameSiteLaxMode
	case SameSiteNone:
		return http.SameSiteNoneMode
	default:
		return http.SameSiteDefaultMode
	}
}

// SetAuthCookie generates a JWT token, sets it as a cookie, and returns any error encountered.
func SetAuthCookie(w http.ResponseWriter, neutralResponseErr string, jwtKey interface{}, httpOnly, secure bool, data map[string]interface{}) error {

	// Check if jwtKey is a string
	secretKey, errKey := jwtKey.(string)
	if !errKey {
		// Handle the case where jwtKey is not a string
		logger.GetLogger().Error("Error loading .env file in 'SetAuthCookie'. JwtKey is not a string. Error: ", errKey)

		// Neutral message
		errHandler.HandleError(w, neutralResponseErr, errHandler.InternalServerError)
		return errors.New("JwtKey is not a string")
	}
	secret, _ := hex.DecodeString(secretKey)

	expirationTimeJWT := time.Second * time.Duration(config.AuthCookieJWTLifetimeSeconds) // Token will expire in 24 hours
	expirationTimeCookie := time.Second * time.Duration(config.AuthCookieLifetimeSeconds) // Token will expire in 24 hours

	// Create JWT token
	token, err := cookie.CreateJWTToken(data, secret, expirationTimeJWT)
	if err != nil {
		return err
	}

	// Set JWT token as a cookie
	cookieName := config.AuthCookieName
	err = cookie.SetCookie(w, cookieName, token, expirationTimeCookie, httpOnly, secure, cookie.SameSiteDefault)
	if err != nil {
		return err
	}

	return nil
}

// ///////////////////////////////////////////////////////////////////
// SIGNOUT ///////////////////////////////////////////////////////////
//
// Signout
func Signout(w http.ResponseWriter) error {

	// Cookie name
	cookieName := config.AuthCookieName

	err := cookie.SetCookie(w, cookieName, "", -3600, true, true, cookie.SameSiteDefault)
	if err != nil {
		return err
	}

	return nil
}
