package plantcontroller

import (
	config "github.com/paulmuenzner/powerplantmanager/config"
	model "github.com/paulmuenzner/powerplantmanager/models"
	errHandler "github.com/paulmuenzner/powerplantmanager/services/errorHandler"
	responsehandler "github.com/paulmuenzner/powerplantmanager/services/responseHandler"
	crypto "github.com/paulmuenzner/powerplantmanager/utils/crypto"
	logger "github.com/paulmuenzner/powerplantmanager/utils/logs"
	mongodb "github.com/paulmuenzner/powerplantmanager/utils/mongoDB"
	stringHandler "github.com/paulmuenzner/powerplantmanager/utils/strings"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
)

func SetKeySecret(mongoDBInterface *mongodb.MethodInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//////////////////////////////////////////////////////
		///////// SETUP //////////////////////////////////////
		//
		neutralResponseErr := "We appologize. Adding a plant is currently not possible due to github.com/paulmuenzner/powerplantmanager update. Please, try again later."

		/////////////////////////////////////////////////////////////
		// Generate Keys
		secret, err := crypto.ByteSize18.GenerateKey()
		if err != nil {
			logger.GetLogger().Error("Error creating secret in SetKeySecret. Error: ", err)
			// Neutral message
			errHandler.HandleError(w, neutralResponseErr, errHandler.InternalServerError)
			return
		}
		hashedSecret, err := crypto.Hash(secret)
		if err != nil {
			logger.GetLogger().Error("Cannot hash password in SetKeySecret: ", err)
			errHandler.HandleError(w, neutralResponseErr, errHandler.InternalServerError)
			return
		}
		plantKey := stringHandler.GenerateRandomNumericString(40)

		//////////////////////////////////////////////
		// ACCESS REQUEST ATTACHMENT /////////////////
		//
		// Access the parsed JSON data from the context attached in SetPlantConfigValidation
		plantQuery, ok := r.Context().Value("plantRequest").(model.PhotovoltaicPlant)
		if !ok {
			logger.GetLogger().Error("Cannot parse and access JSON of plantQuery in SetPlantConfig.")
			errHandler.HandleError(w, neutralResponseErr, errHandler.InternalServerError)
			return
		}

		// Update PlantLoggerConfig finally
		filterUpdate := bson.M{"_id": plantQuery.ID}
		update := bson.M{"$set": bson.M{"key": plantKey, "secret": hashedSecret}}

		_, errUpdate := mongoDBInterface.RepositoryInterface.UpdateOneInMongo(config.DatabaseNamePlantLoggerConfig, filterUpdate, update, config.CollectionNamePlantLoggerConfig)
		if errUpdate != nil {
			logger.GetLogger().Error("Update of verified account not possible RegistrationVerify: ", errUpdate)
			errHandler.HandleError(w, neutralResponseErr, errHandler.InternalServerError)
			return
		}

		//////////////////////////////////////////////
		// POSITIVE RESPONSE /////////////////////////
		//
		type Data struct {
			Key    string
			Secret string
		}

		// Create a struct literal with the data
		data := Data{
			Key:    plantKey,
			Secret: secret,
		}

		responsehandler.HandleSuccess(w, "New key and secret created. Please note them in a safe place. We cannot retrieve them again. If they get lost, you must create a new key and secret here.", responsehandler.OK, data)

	}
}
