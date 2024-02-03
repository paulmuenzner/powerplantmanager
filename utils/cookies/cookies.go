package cookie

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// Generate and set cookie
func SetCookie(w http.ResponseWriter, nameCookies, token string, expirationTime time.Duration, httpOnly, secure bool, sameSite SameSiteOption) error {
	// Convert SameSiteOption to http.SameSite
	sameSiteValue := convertSameSiteOption(sameSite)

	cookie := &http.Cookie{
		Name:     nameCookies,
		Value:    token,
		Path:     "/",
		MaxAge:   int(expirationTime.Seconds()), // Convert duration to seconds
		HttpOnly: httpOnly,
		Secure:   secure,
		SameSite: sameSiteValue,
	}

	http.SetCookie(w, cookie)
	return nil
}

func CreateJWTToken(data map[string]interface{}, secretKey []byte, expirationTime time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"data": data, // Your custom claims/data
		"exp":  time.Now().Add(expirationTime).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

// ///////////////////////////////////////////////////////////////////
// GET COOKIE
func GetCookie(r *http.Request, nameCookie string) (string, bool) {
	cookie, err := r.Cookie(nameCookie)
	if err != nil {
		return "", false
	}

	return cookie.Value, true
}

// ////////////////////////////////////////////////////////////////////////////////////////
// VerifyAndExtractClaims verifies a JWT token, extracts claims, and returns them.
func VerifyAndExtractClaims(tokenString string, secretKey []byte) (jwt.MapClaims, error) {
	claims := jwt.MapClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}

	return claims, nil
}

//////////////////////////////////

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
