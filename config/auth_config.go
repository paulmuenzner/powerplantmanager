package config

// Add more constants as needed

const (
	MinLengthPassword     int = 12
	MaxLengthPassword     int = 50
	MaxLengthEmailAddress int = 60
	// Cookie
	AuthCookieLifetimeSeconds    int    = 1 * 60 * 60 // One hour
	AuthCookieJWTLifetimeSeconds int    = 1 * 60 * 60 //1 * 60 * 60 // One hour
	AuthCookieName               string = "authCookie"
	//
	TimeValidVerifyTokenMinutes int = 15
)
