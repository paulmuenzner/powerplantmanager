package plantcontroller

import (
	config "github.com/paulmuenzner/powerplantmanager/config"
	model "github.com/paulmuenzner/powerplantmanager/models"
	errHandler "github.com/paulmuenzner/powerplantmanager/services/errorHandler"
	responsehandler "github.com/paulmuenzner/powerplantmanager/services/responseHandler"
	"github.com/paulmuenzner/powerplantmanager/utils/data"
	logger "github.com/paulmuenzner/powerplantmanager/utils/logs"
	mongodb "github.com/paulmuenzner/powerplantmanager/utils/mongoDB"
	"net/http"

	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddLogEntry(mongoDBInterface *mongodb.MethodInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//////////////////////////////////////////////////////
		///////// SETUP //////////////////////////////////////
		//
		neutralResponseErr := "Access currently not possible due to internal github.com/paulmuenzner/powerplantmanager update. Our technical team is informed and working on it."

		// Access the parsed JSON data from the context
		dataBody, ok := r.Context().Value("requestBody").(map[string]interface{})
		if !ok {
			logger.GetLogger().Error("Cannot parse and access JSON of requestBody in 'AddLogEntry'.")
			errHandler.HandleError(w, neutralResponseErr, errHandler.InternalServerError)
			return
		}

		voltageOutput := dataBody["voltageOutput"].(float64)
		currentOutput := dataBody["currentOutput"].(float64)
		powerOutput := dataBody["powerOutput"].(float64)
		solarRadiation := dataBody["solarRadiation"].(float64)
		tAmbient := dataBody["tAmbient"].(float64)
		tModule := dataBody["tModule"].(float64)
		relativeHumidity := dataBody["relHumidity"].(float64)
		windSpeed := dataBody["windSpeed"].(float64)

		// Access the parsed JSON plantConfig from the context
		plantConfig, ok := r.Context().Value("collectionNameLogger").(interface{})
		if !ok {
			logger.GetLogger().Error("Cannot parse and access JSON of plantConfig in 'AddLogEntry'.")
			errHandler.HandleError(w, neutralResponseErr, errHandler.InternalServerError)
			return
		}
		collectionName := plantConfig.(string)

		// PLANT LOG
		// Prepare data to save new plant document
		objectIDPlant := primitive.NewObjectID()
		dataToSaveNewPlantLog := model.PlantLogger{
			ID:                 objectIDPlant,
			VoltageOutput:      voltageOutput,
			CurrentOutput:      currentOutput,
			PowerOutput:        powerOutput,
			SolarRadiation:     solarRadiation,
			AmbientTemperature: tAmbient,
			ModuleTemperature:  tModule,
			RelativeHumidity:   relativeHumidity,
			WindSpeed:          windSpeed,
			CreatedAt:          time.Now(), // Prepare Use function
		}

		// Validate data against mongodb pv_plant_model
		if err := data.ValidateStruct(dataToSaveNewPlantLog); err != nil {
			logger.GetLogger().Error("Data validation against mongodb pv_plant_model failed in 'AddLogEntry'. Error: ", err)
			// Neutral response
			errHandler.HandleError(w, neutralResponseErr, errHandler.InternalServerError)
			return
		}

		// Save document to PhotovoltaicPlant model
		_, err := mongoDBInterface.RepositoryInterface.InsertOneToMongo(config.DatabaseNamePlantLogger, dataToSaveNewPlantLog, collectionName)
		if err != nil {
			logger.GetLogger().Error("Unable to save new plant in 'AddLogEntry'. Error: ", err)
			// Neutral response
			errHandler.HandleError(w, neutralResponseErr, errHandler.InternalServerError)
			return
		}

		responsehandler.HandleSuccess(w, "New log added.", responsehandler.OK)

	}
}
