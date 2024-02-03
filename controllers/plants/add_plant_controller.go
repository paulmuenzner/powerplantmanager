package plantcontroller

import (
	"context"
	config "github.com/paulmuenzner/powerplantmanager/config"
	model "github.com/paulmuenzner/powerplantmanager/models"
	errHandler "github.com/paulmuenzner/powerplantmanager/services/errorHandler"
	responsehandler "github.com/paulmuenzner/powerplantmanager/services/responseHandler"
	cookie "github.com/paulmuenzner/powerplantmanager/utils/cookies"
	"github.com/paulmuenzner/powerplantmanager/utils/data"
	"github.com/paulmuenzner/powerplantmanager/utils/date"
	logger "github.com/paulmuenzner/powerplantmanager/utils/logs"
	mongodb "github.com/paulmuenzner/powerplantmanager/utils/mongoDB"
	stringHandler "github.com/paulmuenzner/powerplantmanager/utils/strings"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func AddPlant(mongoDBInterface *mongodb.MethodInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//////////////////////////////////////////////////////
		///////// SETUP //////////////////////////////////////
		//
		neutralResponseErr := "We appologize. Adding a plant is currently not possible due to github.com/paulmuenzner/powerplantmanager update. Please, try again later."
		timeStamp := date.TimeStamp()

		// Access the parsed JSON data from the context
		dataBody, ok := r.Context().Value("requestBody").(map[string]interface{})
		if !ok {
			logger.GetLogger().Error("Cannot parse and access JSON in 'RegistrationVerify()'.")
			errHandler.HandleError(w, neutralResponseErr, errHandler.InternalServerError)
			return
		}

		// Input data
		name := dataBody["name"].(string)

		//////////////////////////////////////////////////////
		///////// TOKEN HANDLING /////////////////////////////
		//
		// Extract data from JWT in cookie
		claimData, err := cookie.GetCookieData(r, config.AuthCookieName)
		if err != nil {
			logger.GetLogger().Errorf("Cannot extract data/claim from cookie in controller 'AddPlant()' using 'GetCookieData()'. Cookie name: %s. Error: %v", config.AuthCookieName, err)
			// Neutral message
			errHandler.HandleError(w, neutralResponseErr, errHandler.InternalServerError)
			return
		}

		// Extracting the userId value from the claim data and converting its type
		// Extract
		userID, ok := claimData["data"].(map[string]interface{})["userId"].(string)
		if !ok {
			logger.GetLogger().Errorf("Failed to extract the userId value from the claim data and converting its type in 'AddPlant()'. Raw claim data: %+v", claimData["data"].(map[string]interface{})["userId"])
			// Neutral message
			errHandler.HandleError(w, neutralResponseErr, errHandler.InternalServerError)
			return
		}

		// Convert type
		userObjectID, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			logger.GetLogger().Errorf("Failed to convert hex value as string to ObjectID in 'AddPlant()' using 'ObjectIDFromHex()'. Hex value user id: %s. Error: %v", userID, err)
			// Neutral message
			errHandler.HandleError(w, neutralResponseErr, errHandler.InternalServerError)
			return

		}

		//////////////////////////////////////////////////////
		///////// STORE DATA  ////////////////////////////////
		// Start a session for the transaction
		session, err := mongoDBInterface.RepositoryInterface.StartSession()
		if err != nil {
			logger.GetLogger().Errorf("Unable to start session for transaction in 'AddPlant()'. Error: %v", err)
			errHandler.HandleError(w, neutralResponseErr, errHandler.InternalServerError)
			return
		}
		defer session.EndSession(context.Background())

		///////// Two new entries and one new and separate collection for logging related plant //////////
		// 1 ADD PLANT
		// Prepare data to save new plant document
		objectIDPlant := primitive.NewObjectID() // Same Mongo Object ID for plant collection (PhotovoltaicPlant) and _id in PlantLoggerConfig
		publicPlantID := stringHandler.GenerateRandomNumericString(15)
		dataToSaveNewPlant := model.PhotovoltaicPlant{
			ID:            objectIDPlant,
			PublicPlantID: publicPlantID,
			Name:          name,
			User:          userObjectID,
			CreatedAt:     timeStamp,
		}

		// Validate data against mongodb pv_plant_model
		if err := data.ValidateStruct(dataToSaveNewPlant); err != nil {
			logger.GetLogger().Errorf("Data validation against mongodb pv_plant_model failed in 'AddPlant()' using 'ValidateStruct()'. Data to save: %+v. Error: %v", dataToSaveNewPlant, err)
			// Handle the error appropriately
			return
		}

		// 2 ADD PLANT LOGGER CONFIG
		// Prepare data to save new plant logger document
		collectionIDLogger := stringHandler.GenerateRandomNumericString(12)
		urlID := stringHandler.GenerateRandomNumericString(20)
		collectionNamePlantLogger := "plant_logger_" + collectionIDLogger
		ips := make([]string, 0)
		dataToSaveNewPlantLoggerConfig := model.PlantLoggerConfig{
			ID:                   objectIDPlant,
			PublicPlantID:        publicPlantID,
			IntervalSec:          config.IntervalSecDefault,
			URLID:                urlID,
			CollectionNameLogger: collectionNamePlantLogger,
			IPWhitelist:          ips,
			CreatedAt:            timeStamp,
		}

		// Validate data entry against mongodb plant_logger
		if err := data.ValidateStruct(dataToSaveNewPlantLoggerConfig); err != nil {
			logger.GetLogger().Errorf("Data validation against mongodb plant_logger_config failed in 'AddPlant()' using 'ValidateStruct()'. Data to save: %+v. Error: %v", dataToSaveNewPlantLoggerConfig, err)
			errHandler.HandleError(w, neutralResponseErr, errHandler.InternalServerError)
			return
		}

		// 3 CREATE COLLECTION FOR LOGGER
		// Prepare data to save new plant
		// Set up the database and collection names
		databaseNamePlantLogger := config.DatabaseNamePlantLogger

		// Define the transaction function
		transactionFunc := func(sessionContext mongo.SessionContext) (interface{}, error) {
			/////////////////////////////////////////////////////////////////
			// NEW PLANT
			// Save document to PhotovoltaicPlant model
			_, err = mongoDBInterface.RepositoryInterface.InsertOneToMongo(config.DatabaseNamePlants, dataToSaveNewPlant, config.CollectionNamePhotovoltaicPlant)
			if err != nil {
				logger.GetLogger().Errorf("Unable to save new plant in 'AddPlant()' using 'InsertOneToMongo()'. Collection name: %s. Error: %v", config.CollectionNamePlantLoggerConfig, err)
				return nil, err
			}

			/////////////////////////////////////////////////////////////////
			// NEW PLANT LOGGER CONFIG
			// Save document to PlantLoggerConfig
			_, err = mongoDBInterface.RepositoryInterface.InsertOneToMongo(config.DatabaseNamePlantLoggerConfig, dataToSaveNewPlantLoggerConfig, config.CollectionNamePlantLoggerConfig)
			if err != nil {
				logger.GetLogger().Errorf("Unable to save new plant logger config in 'AddPlant()' using 'InsertOneToMongo()'. Collection name: %s. Error: %v", config.CollectionNamePlantLoggerConfig, err)
				return nil, err
			}

			/////////////////////////////////////////////////////////////////
			// NEW PLANT LOGGER COLLECTION
			// Create new plant logger collection including index setup
			err = mongoDBInterface.RepositoryInterface.CreateNewCollection(databaseNamePlantLogger, collectionNamePlantLogger)
			if err != nil {
				logger.GetLogger().Errorf("Unable to setup new plant logger collection in 'AddPlant()' using 'CreateNewCollection()'. Error: %v", err)
				return nil, err
			}

			return "Transaction completed successfully", nil
		}

		// Start the transaction
		_, err = session.WithTransaction(context.Background(), transactionFunc)
		if err != nil {
			logger.GetLogger().Error("Transaction error in 'AddPlant()' using 'WithTransaction()'. Cannot save two new documents and a new database. Error: ", err)
			errHandler.HandleError(w, neutralResponseErr, errHandler.InternalServerError)
			return
		}

		///////////////// INDEXES ////////////////////////////////////////////////////////////
		// After transaction succeeded, we need to setup indexes in our new collection
		// Create an index for public plant id in PhotovoltaicPlant model
		err = mongoDBInterface.RepositoryInterface.CreateUniqueIndex(config.CollectionNamePhotovoltaicPlant, config.DatabaseNamePlants, "public_plant_id", true)
		if err != nil {
			logger.GetLogger().Error("CreateUniqueIndex error in 'AddPlant()' using 'CreateUniqueIndex()'. Cannot create new index in PhotovoltaicPlant model on 'public_plant_id'. Error: ", err)
			errHandler.HandleError(w, neutralResponseErr, errHandler.InternalServerError)
			return
		}

		responsehandler.HandleSuccess(w, "New plant added to your account.", responsehandler.OK)

		return

	}
}
