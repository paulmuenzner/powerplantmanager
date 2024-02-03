package authcontroller

import (
	config "github.com/paulmuenzner/powerplantmanager/config"
	model "github.com/paulmuenzner/powerplantmanager/models"
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

func RegistrationVerify(emailInterface *email.RepositoryInterface, mongoDBInterface *mongodb.MethodInterface) http.HandlerFunc { // prepare email after successful verification
	return func(w http.ResponseWriter, r *http.Request) {
		// Prepare email
		//////////////////////////////////////////////////////
		///////// SETUP //////////////////////////////////////
		//
		neutralResponseErr := "We appologize. Verification currently not possible due to github.com/paulmuenzner/powerplantmanager update. Please, try again later."

		// Access the parsed JSON data from the context
		dataBody, ok := r.Context().Value("requestBody").(map[string]interface{})
		if !ok {
			logger.GetLogger().Error("Cannot parse and access JSON in RegistrationVerify.")
			errHandler.HandleError(w, neutralResponseErr, errHandler.InternalServerError)
			return
		}

		// Input data
		password := dataBody["password"].(string)
		verifyToken := dataBody["verifyToken"].(string)

		//////////////////////////////////////////////////////
		///////// TOKEN HANDLING /////////////////////////////
		//
		// Decrypt encrypted verifyToken
		// Get key to encrypt verifyToken
		key, err := env.GetEnvValue("KEY_VERIFY_TOKEN", "")
		if err != nil { // prepare error handlung
			logger.GetLogger().Error("Error loading .env file in 'RegistrationVerify()' using 'GetEnvValue()'. Error: ", err)
			errHandler.HandleError(w, neutralResponseErr, errHandler.InternalServerError)
			return
		}
		result, err := crypto.DecryptVerifyToken(verifyToken, key)
		if err != nil {
			logger.GetLogger().Error("Unable to decrypt encryptedVerifyToken in RegistrationVerify. Error: ", err)
			errHandler.HandleError(w, neutralResponseErr, errHandler.InternalServerError)
			return
		}

		//////////////////////////////////////////////////////
		///////// DATABASE QUERY /////////////////////////////
		//
		// Find user in database by decrypted token only
		var filter bson.M = bson.M{"verify_token": result}
		var user model.UserAuth
		var sort bson.D = bson.D{}
		foundOne, err := mongoDBInterface.RepositoryInterface.FindOneInMongo(config.DatabaseNameUserAuth, filter, config.UserAuthCollectionName, sort, &user)
		if err != nil { // prepare error handlung
			logger.GetLogger().Errorf("Error in 'RegistrationVerify()' using 'FindOneInMongo()' when querying collection '%s' part of database '%s'. Error: %v", config.UserAuthCollectionName, config.DatabaseNameUserAuth, err)
			errHandler.HandleError(w, neutralResponseErr, errHandler.InternalServerError)
			return
		}
		if !foundOne {
			logger.GetLogger().Warnf("Non valid verifyToken used for verifying registration in 'RegistrationVerify()' using 'FindOneInMongo()'when querying collection '%s' part of database '%s'. Error: %v", config.UserAuthCollectionName, config.DatabaseNameUserAuth, err)
			errHandler.HandleError(w, "URL not valid anymore. Please try to register again.", errHandler.NotAcceptable)
			return
		}

		//////////////////////////////////////////////////////
		///////// TOKEN VALIDATION ///////////////////////////
		//
		// Validate if verifyToken is still valid or expired
		dateVerifyToken := user.DateVerifyToken

		verifyTokenTimeAgo := dateVerifyToken.Add(+15 * time.Minute)
		if time.Now().After(verifyTokenTimeAgo) {
			errHandler.HandleError(w, "This URL has been expired. Please register again.", errHandler.NotAcceptable)
			return
		}

		//////////////////////////////////////////////////////////////////////////
		// If user in database save new password and set verify to true
		// Define a filter to match the document(s) you want to update
		filterUpdate := bson.M{"_id": user.ID}

		// Hash password
		passwordHash, err := crypto.Hash(password)
		if err != nil {
			logger.GetLogger().Error("Cannot hash password in 'RegistrationVerify()' using 'crypto.Hash()': ", err)
			errHandler.HandleError(w, neutralResponseErr, errHandler.InternalServerError)
			return
		}

		// Define an update to set a new value for a field
		update := bson.M{"$set": bson.M{"verified": true, "password": passwordHash}, "$unset": bson.M{"verify_token": "", "date_verify_token": ""}}
		_, err = mongoDBInterface.RepositoryInterface.UpdateOneInMongo(config.DatabaseNameUserAuth, filterUpdate, update, config.UserAuthCollectionName)

		if err != nil {
			logger.GetLogger().Errorf("Update of verified account not possible in 'RegistrationVerify()' using 'UpdateOneInMongo()': %v", err)
			errHandler.HandleError(w, neutralResponseErr, errHandler.InternalServerError)
			return
		}

		responsehandler.HandleSuccess(w, "Great, registration process completed. Feel free to login.", responsehandler.OK)

	}
}
