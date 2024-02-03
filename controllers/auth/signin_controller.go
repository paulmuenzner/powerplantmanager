package authcontroller

import (
	config "github.com/paulmuenzner/powerplantmanager/config"
	model "github.com/paulmuenzner/powerplantmanager/models"
	cookiehandler "github.com/paulmuenzner/powerplantmanager/services/cookieHandler"
	errHandler "github.com/paulmuenzner/powerplantmanager/services/errorHandler"
	responsehandler "github.com/paulmuenzner/powerplantmanager/services/responseHandler"
	crypto "github.com/paulmuenzner/powerplantmanager/utils/crypto"
	"github.com/paulmuenzner/powerplantmanager/utils/email"
	env "github.com/paulmuenzner/powerplantmanager/utils/env"
	logger "github.com/paulmuenzner/powerplantmanager/utils/logs"
	mongodb "github.com/paulmuenzner/powerplantmanager/utils/mongoDB"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func Signin(emailInterface *email.RepositoryInterface, mongoDBInterface *mongodb.MethodInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//////////////////////////////////////////////////////
		///////// SETUP //////////////////////////////////////
		//
		neutralResponseErr := "We appologize. Signin currently not possible due to github.com/paulmuenzner/powerplantmanager update. Please, try again later."

		// Access the parsed JSON data from the context
		dataBody, ok := r.Context().Value("requestBody").(map[string]interface{})
		if !ok {
			logger.GetLogger().Error("Cannot parse and access JSON in RegistrationVerify. Error: ", ok)
			errHandler.HandleError(w, neutralResponseErr, errHandler.InternalServerError)
			return
		}

		// Input data
		password := dataBody["password"].(string)
		email := dataBody["email"].(string)
		timeNow := time.Now()

		//////////////////////////////////////////////////////
		///////// DATABASE QUERY /////////////////////////////
		//
		// Find user in database using email
		var filter bson.M = bson.M{"email": email}
		var user model.UserAuth
		var sort bson.D = bson.D{}
		findOne, err := mongoDBInterface.RepositoryInterface.FindOneInMongo(config.DatabaseNameUserAuth, filter, config.UserAuthCollectionName, sort, &user)
		if err != nil {
			logger.GetLogger().Errorf("Error in 'Signin()' using 'FindOneInMongo()' when querying collection '%s' part of database '%s'. Error: %v", config.UserAuthCollectionName, config.DatabaseNameUserAuth, err)
			errHandler.HandleError(w, neutralResponseErr, errHandler.InternalServerError)
			return
		}
		// Neutral response if email not in database
		if !findOne {
			errHandler.HandleError(w, "Email and/or password not correct.", errHandler.BadRequest)
			return
		}

		//////////////////////////////////////////////////////
		///////// PASSWORD VALIDATION ////////////////////////
		//
		passwordHash := user.Password

		validatePassword := crypto.IsHashValid(password, passwordHash)
		if !validatePassword { //////////////////////////////////////////////////////
			///////// WRONG PASSWORD /////////////////////////////
			//
			errHandler.HandleError(w, "Email and/or password not correct.", errHandler.BadRequest)
			//
			// For safety
			// Inform user that somebody tried to login with non matching password
			err = emailInterface.RepositoryInterface.EmailInformUserFailedLogin(timeNow, email)
			if err != nil {
				logger.GetLogger().Error("Unable to send email in 'Signin()'using 'EmailInformUserFailedLogin()'. Error: ", err)
				return
			}
			return
		}

		//////////////////////////////////////////////////////
		///////// SET AUTH COOKIE ////////////////////////////
		//
		dataCookie := map[string]interface{}{"userId": user.ID, "email": email}
		jwtKey, err := env.GetEnvValue("JWT_SECRET_KEY", "")
		if err != nil {
			logger.GetLogger().Errorf("Error loading env value for 'JWT_SECRET_KEY' in 'Login()' using 'GetEnvValue()'. Error: %v", err)
		}

		err = cookiehandler.SetAuthCookie(w, neutralResponseErr, jwtKey, true, true, dataCookie)
		if err != nil {
			logger.GetLogger().Error("Error with setting a cookie in 'Login'. Error: ", err)
			errHandler.HandleError(w, neutralResponseErr, errHandler.InternalServerError)
			return
		}
		responsehandler.HandleSuccess(w, "Welcome!", responsehandler.Accepted)

	}
}
