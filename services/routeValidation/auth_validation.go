package routevalidation

import (
	config "github.com/paulmuenzner/powerplantmanager/config"
	errHandler "github.com/paulmuenzner/powerplantmanager/services/errorHandler"
	responsehandler "github.com/paulmuenzner/powerplantmanager/services/responseHandler"
	cookie "github.com/paulmuenzner/powerplantmanager/utils/cookies"
	logger "github.com/paulmuenzner/powerplantmanager/utils/logs"
	mongodb "github.com/paulmuenzner/powerplantmanager/utils/mongoDB"
	"net/http"

	v "github.com/paulmuenzner/powerplantmanager/utils/validate"
)

// /////////////////////////////////////////////////////////////////////////////////////////////
// REGISTRATION VALIDATION
// ///////////////////////
func RegistrationValidation(next http.HandlerFunc, mongoDBInterface *mongodb.MethodInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		neutralResponseErr := "We appologize. Registration currently not possible due to github.com/paulmuenzner/powerplantmanager update. Please, try again later."
		//////////////////////////////////////////////
		// VALIDATE AUTH STATUS //////////////////////
		//
		// Validate if already logged in
		value, hasAuthCookie := cookie.GetCookie(r, config.AuthCookieName)
		if hasAuthCookie && len(value) > 0 {
			expired := cookie.HasCookieExpired(r, config.AuthCookieName)
			if !expired {
				errHandler.HandleError(w, "No point to register; you are already signed in.", errHandler.InternalServerError)
				return
			}
		}

		//////////////////////////////////////////////
		// REQUEST BODY VALIDATION ///////////////////
		//
		// Access the parsed JSON data from the context
		data, ok := r.Context().Value("requestBody").(map[string]interface{})
		if !ok {
			errHandler.HandleError(w, "Internal Server Error", errHandler.InternalServerError)
			return
		}

		// Verify and Define variables
		email, emailValid := data["email"].(string)
		if !emailValid {
			errHandler.HandleError(w, neutralResponseErr, errHandler.InternalServerError)
			return
		}

		// Validate if request body exactly contains number and names of expected keys
		expectedKeys := []string{"email"}
		validateKeys := v.Validate(data).
			HasMapExactKeys(expectedKeys).
			GetResult()

		if len(validateKeys) > 0 {
			errHandler.HandleError(w, validateKeys[0], errHandler.BadRequest)
			return
		}

		// Validate email
		validateEmail := v.Validate(email).
			IsEmail().
			MaxLength(config.MaxLengthEmailAddress).
			GetResult()

		if len(validateEmail) > 0 {
			errHandler.HandleError(w, validateEmail[0], errHandler.BadRequest)
			return
		}

		// Check if email already in database. Each email must be unique
		exists, _ := mongoDBInterface.RepositoryInterface.IsValueInCollection(config.DatabaseNameUserAuth, config.UserAuthCollectionName, "email", data["email"].(string))
		if exists {
			// Prepare email notification
			errHandler.HandleError(w, "Please choose another email address.", errHandler.BadRequest)
			return
		}

		// Call the next handler if validation passes
		next.ServeHTTP(w, r)
	}
}

// /////////////////////////////////////////////////////////////////////////////////////////////
// LOGIN VALIDATION
// ////////////////
func SigninValidation(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		neutralResponseErr := "We appologize. Login currently not possible due to github.com/paulmuenzner/powerplantmanager update. Please, try again later."

		//////////////////////////////////////////////
		// VALIDATE AUTH STATUS //////////////////////
		//
		// Validate if already logged in
		value, hasAuthCookie := cookie.GetCookie(r, config.AuthCookieName)
		if hasAuthCookie && len(value) > 0 {
			expired := cookie.HasCookieExpired(r, config.AuthCookieName)
			if !expired {
				errHandler.HandleError(w, "You are already signed in.", errHandler.InternalServerError)
				return
			}
		}

		//////////////////////////////////////////////
		// REQUEST BODY VALIDATION ///////////////////
		//
		// Access the parsed JSON data from the context
		data, ok := r.Context().Value("requestBody").(map[string]interface{})
		if !ok {
			logger.GetLogger().Error("Cannot parse and access JSON in LoginValidation. Error: ", ok)
			errHandler.HandleError(w, "Internal Server Error", errHandler.InternalServerError)
			return
		}

		// Verify and Define variables
		password, passwordValid := data["password"].(string)
		if !passwordValid {
			errHandler.HandleError(w, neutralResponseErr, errHandler.InternalServerError)
			return
		}

		email, emailValid := data["email"].(string)
		if !emailValid {
			errHandler.HandleError(w, neutralResponseErr, errHandler.InternalServerError)
			return
		}

		// Validate if request body exactly contains number and names of expected keys
		expectedKeys := []string{"email", "password"}
		validateKeys := v.Validate(data).
			HasMapExactKeys(expectedKeys, neutralResponseErr).
			GetResult()

		if len(validateKeys) > 0 {
			errHandler.HandleError(w, validateKeys[0], errHandler.BadRequest)
			return
		}

		// Validate password. Only accept strings.
		// Only check max length in order to reduce unnecessary load. Don't use config.MinLengthPassword and pick a larger length (eg 100) as MinLengthPassword settings might have been different (longer) compared to now
		// No password regex here as this would only give hints to potential attackers
		validatePassword := v.Validate(password).
			MaxLength(100).
			GetResult()

		if len(validatePassword) > 0 {
			errHandler.HandleError(w, validatePassword[0], errHandler.BadRequest)
			return
		}

		// Validate email
		validateEmail := v.Validate(email).
			IsEmail().
			MaxLength(100).
			GetResult()

		if len(validateEmail) > 0 {
			errHandler.HandleError(w, validateEmail[0], errHandler.BadRequest)
			return
		}

		//////////////////////////////////////////////
		// PASS VALID REQUEST ////////////////////////
		//
		// Call the next handler if validation passes
		next.ServeHTTP(w, r)
	}
}

// /////////////////////////////////////////////////////////////////////////////////////////////
// REGISTRATION VERIFICATION VALIDATION
// ////////////////////////////////////
func VerificationValidation(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		neutralResponseErr := "We appologize. Verification currently not possible due to github.com/paulmuenzner/powerplantmanager update. Please, try again later."
		//////////////////////////////////////////////
		// VALIDATE AUTH STATUS //////////////////////
		//
		// Validate if already logged in
		value, hasAuthCookie := cookie.GetCookie(r, config.AuthCookieName)
		if hasAuthCookie && len(value) > 0 {
			expired := cookie.HasCookieExpired(r, config.AuthCookieName)
			if !expired {
				errHandler.HandleError(w, "You are already signed in. No point to validate registration.", errHandler.InternalServerError)
				return
			}
		}

		//////////////////////////////////////////////
		// REQUEST BODY VALIDATION ///////////////////
		//
		// Access the parsed JSON data from the context
		data, ok := r.Context().Value("requestBody").(map[string]interface{})
		if !ok {
			errHandler.HandleError(w, neutralResponseErr, errHandler.InternalServerError)
			return
		}

		// Verify and Define variables
		password, passwordValid := data["password"].(string)
		if !passwordValid {
			errHandler.HandleError(w, neutralResponseErr, errHandler.InternalServerError)
			return
		}

		verifyToken, verifyTokenValid := data["verifyToken"].(string)
		if !verifyTokenValid {
			errHandler.HandleError(w, neutralResponseErr, errHandler.InternalServerError)
			return
		}

		passwordVerify, passwordVerifyValid := data["passwordVerify"].(string)
		if !passwordVerifyValid {
			errHandler.HandleError(w, neutralResponseErr, errHandler.InternalServerError)
			return
		}

		// Validate if request body exactly contains number and names of expected keys
		expectedKeys := []string{"password", "passwordVerify", "verifyToken"}
		validateKeys := v.Validate(data).
			HasMapExactKeys(expectedKeys).
			GetResult()

		if len(validateKeys) > 0 {
			errHandler.HandleError(w, validateKeys[0], errHandler.BadRequest)
			return
		}

		// Validate password
		validatePassword := v.Validate(password).
			IsPasswordValid().
			MaxLength(config.MaxLengthPassword).
			MinLength(config.MinLengthPassword).
			GetResult()

		if len(validatePassword) > 0 {
			errHandler.HandleError(w, validatePassword[0], errHandler.BadRequest)
			return
		}

		// Validate verification key
		validateVerifyToken := v.Validate(verifyToken).
			MinLength(20, "No valid URL").
			MaxLength(150, "No valid URL").
			GetResult()

		if len(validateVerifyToken) > 0 {
			errHandler.HandleError(w, validateVerifyToken[0], errHandler.BadRequest)
			return
		}

		// Validate equality of password and passwordVerify
		if password != passwordVerify {
			errHandler.HandleError(w, "Passwords don't match.", errHandler.BadRequest)
			return
		}

		//////////////////////////////////////////////
		// PASS VALID REQUEST ////////////////////////
		//
		// Call the next handler if validation passes
		next.ServeHTTP(w, r)
	}
}

// /////////////////////////////////////////////////////////////////////////////////////////////
// SIGNOUT VERIFICATION VALIDATION
// ////////////////////////////////////
func SignoutValidation(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//////////////////////////////////////////////
		// VALIDATE AUTH STATUS //////////////////////
		//
		// Validate if already signed out
		value, hasAuthCookie := cookie.GetCookie(r, config.AuthCookieName)
		if !hasAuthCookie || len(value) < 1 {
			responsehandler.HandleSuccess(w, "Signed out successfully.", responsehandler.OK)
			return
		}

		//////////////////////////////////////////////
		// REQUEST BODY VALIDATION ///////////////////
		//
		// Access the parsed JSON data from the context
		data, ok := r.Context().Value("requestBody").(map[string]interface{})
		if !ok {
			errHandler.HandleError(w, "Internal Server Error", errHandler.InternalServerError)
			return
		}

		// Validate if request body empty as no data is needed for signout
		if len(data) != 0 {
			errHandler.HandleError(w, "Invalid request", errHandler.BadRequest)
			return
		}

		//////////////////////////////////////////////
		// PASS VALID REQUEST ////////////////////////
		//
		// Call the next handler if validation passes
		next.ServeHTTP(w, r)
	}
}
