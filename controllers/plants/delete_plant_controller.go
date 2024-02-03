package plantcontroller

import (
	"context"
	config "github.com/paulmuenzner/powerplantmanager/config"
	errHandler "github.com/paulmuenzner/powerplantmanager/services/errorHandler"
	responsehandler "github.com/paulmuenzner/powerplantmanager/services/responseHandler"
	logger "github.com/paulmuenzner/powerplantmanager/utils/logs"
	mongodb "github.com/paulmuenzner/powerplantmanager/utils/mongoDB"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func DeletePlant(mongoDBInterface *mongodb.MethodInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Deletion of plant means deleting document in PhotovoltaicPlant and PlantLoggerConfig, and deletion of PlantLoggerCollection

		//////////////////////////////////////////////////////
		///////// SETUP //////////////////////////////////////
		//
		neutralResponseErr := "We appologize. Deletion currently not possible due to github.com/paulmuenzner/powerplantmanager update. Please, try again later."

		//////////////////////////////////////////////
		// REQUEST BODY VALIDATION ///////////////////
		//
		// Access the parsed JSON data from the context
		data, _ := r.Context().Value("requestBody").(map[string]interface{})
		publicPlantID, _ := data["publicPlantID"].(string)

		// Access the parsed JSON data from the context
		collectionNameLogger, ok := r.Context().Value("collectionNameLogger").(string)
		if !ok {
			logger.GetLogger().Error("Cannot parse and access 'collectionNameLogger' in 'DeletePlant()'.")
			errHandler.HandleError(w, neutralResponseErr, errHandler.InternalServerError)
			return
		}

		// Start a session for the transaction
		session, err := mongoDBInterface.RepositoryInterface.StartSession()
		if err != nil {
			logger.GetLogger().Error("Unable to start session for transaction in 'DeletePlant()'. Error: ", err)
			errHandler.HandleError(w, neutralResponseErr, errHandler.InternalServerError)
			return
		}
		defer session.EndSession(context.Background())

		// Input data
		var filterDeleteDocument bson.M = bson.M{"public_plant_id": publicPlantID}

		//////////////////////////////////////////////////////
		///////// DELETE TRANSACTION /////////////////////////
		//
		// Define the transaction function
		transactionFunc := func(sessionContext mongo.SessionContext) (interface{}, error) {
			/////////////////////////////////////////////////////////////////
			// DELETE DOCUMENT IN PlantLoggerConfig
			_, err := mongoDBInterface.RepositoryInterface.DeleteDocumentMongo(config.DatabaseNamePlantLoggerConfig, filterDeleteDocument, config.CollectionNamePlantLoggerConfig)
			if err != nil {
				logger.GetLogger().Error("Unable to delete PlantLoggerConfig document in 'DeletePlant()' using 'DeleteDocumentMongo()'. Error: ", err)
				return nil, err
			}

			/////////////////////////////////////////////////////////////////
			// DELETE DOCUMENT IN PhotovoltaicPlant
			_, err = mongoDBInterface.RepositoryInterface.DeleteDocumentMongo(config.DatabaseNamePlants, filterDeleteDocument, config.CollectionNamePhotovoltaicPlant)
			if err != nil {
				logger.GetLogger().Error("Unable to delete PhotovoltaicPlant document in 'DeletePlant()' using 'DeleteDocumentMongo()'. Error: ", err)
				return nil, err
			}

			/////////////////////////////////////////////////////////////////
			// DELETE LOGGER COLLECTION
			err = mongoDBInterface.RepositoryInterface.DeleteCollectionMongo(config.DatabaseNamePlantLogger, collectionNameLogger)
			if err != nil {
				logger.GetLogger().Error("Unable to delete Logger Collection document in 'DeletePlant()' using 'DeleteDocumentMongo()'. Error: ", err)
				return nil, err
			}

			return "Transaction completed successfully", nil
		}

		// Start the transaction
		_, err = session.WithTransaction(context.Background(), transactionFunc)
		if err != nil {
			logger.GetLogger().Error("Transaction error occurred in 'DeletePlant()' while using 'WithTransaction()'. Unable to delete two documents from the collection. Error: ", err)
			errHandler.HandleError(w, neutralResponseErr, errHandler.InternalServerError)
			return
		}

		responsehandler.HandleSuccess(w, "Deletion accomplished.", responsehandler.OK)

		return

	}
}
