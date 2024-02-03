package plantcontroller

import (
	config "github.com/paulmuenzner/powerplantmanager/config"
	model "github.com/paulmuenzner/powerplantmanager/models"
	errHandler "github.com/paulmuenzner/powerplantmanager/services/errorHandler"
	responsehandler "github.com/paulmuenzner/powerplantmanager/services/responseHandler"
	ip "github.com/paulmuenzner/powerplantmanager/utils/ip"
	logger "github.com/paulmuenzner/powerplantmanager/utils/logs"
	mongodb "github.com/paulmuenzner/powerplantmanager/utils/mongoDB"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
)

func SetPlantConfig(mongoDBInterface *mongodb.MethodInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//////////////////////////////////////////////////////
		///////// SETUP //////////////////////////////////////
		//
		neutralResponseErr := "We appologize. Modifying plant configuration currently not possible due to github.com/paulmuenzner/powerplantmanager updates. Please, try again later."

		//////////////////////////////////////////////
		// REQUEST BODY VALIDATION ///////////////////
		//
		// Access the parsed JSON data from the context
		dataBody, ok := r.Context().Value("requestBody").(map[string]interface{})
		if !ok {
			logger.GetLogger().Error("Cannot parse and access JSON of requestBody in 'SetPlantConfig()'.")
			errHandler.HandleError(w, neutralResponseErr, errHandler.InternalServerError)
			return
		}

		intervalSecRaw := dataBody["intervalSec"].(float64)
		intervalSec := int(intervalSecRaw)
		ipWhiteList, ok := dataBody["ipWhiteList"].([]interface{})
		if !ok {
			logger.GetLogger().Errorf("Cannot parse and access ipWhiteList from dataBody in 'SetPlantConfig()'. dataBody[\"ipWhiteList\"]: %+v", dataBody["ipWhiteList"])
			errHandler.HandleError(w, neutralResponseErr, errHandler.InternalServerError)
			return
		}

		var ipWhiteListExtended []string
		for _, value := range ipWhiteList {
			if ipString, ok := value.(string); ok {
				ipNormalized, _ := ip.NormalizeIP(ipString)
				ipWhiteListExtended = append(ipWhiteListExtended, ipNormalized)
			}
		}

		//////////////////////////////////////////////
		// ACCESS REQUEST ATTACHMENT /////////////////
		//
		// Access the parsed JSON data from the context attached in SetPlantConfigValidation
		plantQuery, ok := r.Context().Value("plantRequest").(model.PhotovoltaicPlant)
		if !ok {
			logger.GetLogger().Errorf("Cannot parse and access JSON of plantQuery in 'SetPlantConfig()'. r.Context().Value(\"plantRequest\"): %+v", r.Context().Value("plantRequest"))
			errHandler.HandleError(w, neutralResponseErr, errHandler.InternalServerError)
			return
		}

		// Update PlantLoggerConfig finally
		filterUpdate := bson.M{"_id": plantQuery.ID}
		update := bson.M{"$set": bson.M{"interval_sec": intervalSec, "ip_whitelist": ipWhiteListExtended}}

		_, errUpdate := mongoDBInterface.RepositoryInterface.UpdateOneInMongo(config.DatabaseNamePlantLoggerConfig, filterUpdate, update, config.CollectionNamePlantLoggerConfig)
		if errUpdate != nil {
			logger.GetLogger().Errorf("Update of verified account not possible in 'SetPlantConfig()' using 'UpdateOneInMongo()': %v", errUpdate)
			errHandler.HandleError(w, neutralResponseErr, errHandler.InternalServerError)
			return
		}

		responsehandler.HandleSuccess(w, "Plant configuration updated.", responsehandler.OK)

	}
}
