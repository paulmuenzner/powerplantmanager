package authcontroller

import (
	"net/http"
	"strconv"
	"time"

	config "github.com/paulmuenzner/powerplantmanager/config"

	model "github.com/paulmuenzner/powerplantmanager/models"
	errHandler "github.com/paulmuenzner/powerplantmanager/services/errorHandler"

	responsehandler "github.com/paulmuenzner/powerplantmanager/services/responseHandler"
	crypto "github.com/paulmuenzner/powerplantmanager/utils/crypto"
	data "github.com/paulmuenzner/powerplantmanager/utils/data"
	"github.com/paulmuenzner/powerplantmanager/utils/date"
	"github.com/paulmuenzner/powerplantmanager/utils/email"
	env "github.com/paulmuenzner/powerplantmanager/utils/env"
	logger "github.com/paulmuenzner/powerplantmanager/utils/logs"
	mongodb "github.com/paulmuenzner/powerplantmanager/utils/mongoDB"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Registration(emailInterface *email.RepositoryInterface, mongoDBInterface *mongodb.MethodInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//////////////////////////////////////////////////////
		///////// SETUP //////////////////////////////////////
		//
		timeStamp := date.TimeStamp()
		neutralResponseErr := "We appologize. Registration currently not possible due to github.com/paulmuenzner/powerplantmanager update. Please, try again later."

		// Access the parsed JSON data from the context
		dataBody, ok := r.Context().Value("requestBody").(map[string]interface{})
		if !ok {
			logger.GetLogger().Error("Cannot parse requestBody json in 'Registration()'.")
			// Neutral message
			errHandler.HandleError(w, neutralResponseErr, errHandler.InternalServerError)
			return
		}
		email := dataBody["email"].(string)
		verifyLinkValidMinutes := strconv.Itoa(config.TimeValidVerifyTokenMinutes)

		/////////////////////////////////////////////////
		// Generate Validation Token
		// Get key to encrypt verifyToken
		key, err := env.GetEnvValue("KEY_VERIFY_TOKEN", "")
		if err != nil {
			logger.GetLogger().Error("Error loading KEY_VERIFY_TOKEN from .env in 'Registration()' using 'GetEnvValue()'. Error: ", err)
			// Neutral message
			errHandler.HandleError(w, neutralResponseErr, errHandler.InternalServerError)
			return
		}

		verifyToken, encryptedVerifyToken, err := crypto.GenerateVerifyToken(key)
		if err != nil {
			logger.GetLogger().Errorf("Error in 'Registration()' using 'GenerateVerifyToken()'. Error: %v", err)
			// Neutral message
			errHandler.HandleError(w, neutralResponseErr, errHandler.InternalServerError)
			return
		}

		// Validate if email already in database
		// Find user in database using email
		var filter bson.M = bson.M{"email": email}
		var sort bson.D = bson.D{}
		var user model.UserAuth
		foundOne, err := mongoDBInterface.RepositoryInterface.FindOneInMongo(config.DatabaseNameUserAuth, filter, config.UserAuthCollectionName, sort, &user)
		// Handle error
		if err != nil {
			logger.GetLogger().Errorf("Error in 'Registration()' using 'FindOneInMongo()' when querying email address %s collection '%s' part of database '%s'. Error: %v", email, config.UserAuthCollectionName, config.DatabaseNameUserAuth, err)
			// Neutral message
			errHandler.HandleError(w, neutralResponseErr, errHandler.InternalServerError)
			return
		}

		// Registration with provided email address found
		if foundOne {
			verified := user.Verified

			// If somebody tries to register an email already existing and verified, a neutral respond is send to the user
			// ...providing no indication to potential attackers whether this email might be already taken or not.
			// An email is send to the users inbox with a notification
			if verified {
				err := emailInterface.RepositoryInterface.EmailRegistrationVerifiedAccount(timeStamp, email)
				if err != nil {
					logger.GetLogger().Error("Unable to send email in 'Registration()' using 'EmailRegistrationVerifiedAccount()'. Error: ", err)
					// Neutral response
					errHandler.HandleError(w, neutralResponseErr, errHandler.InternalServerError)
					return
				}
				responsehandler.HandleSuccess(w, "Great, please verify your mailbox to verify your account.", responsehandler.Accepted)
				return
			}

			// Existing account without valid verification
			if !verified {

				// Update database
				filter := bson.M{"email": email}
				update := bson.M{"$set": bson.M{"verify_token": verifyToken, "date_verify_token": time.Now()}}
				_, errUpdate := mongoDBInterface.RepositoryInterface.UpdateOneInMongo(config.DatabaseNameUserAuth, filter, update, config.UserAuthCollectionName)
				if errUpdate != nil {
					logger.GetLogger().Error("Update of verified account not possible in 'Registration()' using 'UpdateOneInMongo()': ", err)
					errHandler.HandleError(w, neutralResponseErr, errHandler.InternalServerError)
					return
				}

				// Send email with new token
				err := emailInterface.RepositoryInterface.EmailNewRegistration(timeStamp, email, verifyLinkValidMinutes, encryptedVerifyToken)
				if err != nil {
					logger.GetLogger().Error("Unable to send email 'EmailNewRegistration()' in in 'Registration()'. Error: ", err)
					// Neutral respond message
					errHandler.HandleError(w, neutralResponseErr, errHandler.InternalServerError)
					return
				}

				// Response
				responsehandler.HandleSuccess(w, "Great, please verify your mailbox to verify your account.", responsehandler.Accepted)
				return
			}
		}

		// Continue with not existing user (email address)
		///////////////////////////////////////
		// Prepare data to save new user
		objectID := primitive.NewObjectID()
		dataToSave := model.UserAuth{
			ID:              objectID,
			Email:           email,
			VerifyToken:     verifyToken,
			DateVerifyToken: time.Now(),
			CreatedAt:       time.Now(),
		}

		// Validate data against mongodb user_auth_model
		if err := data.ValidateStruct(dataToSave); err != nil {
			logger.GetLogger().Error("Data validation against mongodb user_auth_model failed in 'Registration()'. Error: ", err)
			errHandler.HandleError(w, neutralResponseErr, errHandler.InternalServerError)
			return
		}

		// Save data
		_, err = mongoDBInterface.RepositoryInterface.InsertOneToMongo(config.DatabaseNameUserAuth, dataToSave, config.UserAuthCollectionName)
		if err != nil {
			logger.GetLogger().Error("Unable to save registration data in 'Registration'. Error: ", err)
			errHandler.HandleError(w, neutralResponseErr, errHandler.InternalServerError)
			return
		}

		// Send email with new token
		err = emailInterface.RepositoryInterface.EmailNewRegistration(timeStamp, email, verifyLinkValidMinutes, encryptedVerifyToken)
		if err != nil {
			logger.GetLogger().Error("Unable to send email in 'Registration'. Error: ", err)
			errHandler.HandleError(w, neutralResponseErr, errHandler.InternalServerError)
			return
		}

		responsehandler.HandleSuccess(w, "Great, please verify your mailbox to verify your account.", responsehandler.OK)

	}
}
