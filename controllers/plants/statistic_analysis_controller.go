package plantcontroller

import (
	config "github.com/paulmuenzner/powerplantmanager/config"
	model "github.com/paulmuenzner/powerplantmanager/models"
	errHandler "github.com/paulmuenzner/powerplantmanager/services/errorHandler"
	responsehandler "github.com/paulmuenzner/powerplantmanager/services/responseHandler"
	cookie "github.com/paulmuenzner/powerplantmanager/utils/cookies"
	logger "github.com/paulmuenzner/powerplantmanager/utils/logs"
	mongodb "github.com/paulmuenzner/powerplantmanager/utils/mongoDB"
	"github.com/paulmuenzner/powerplantmanager/utils/statistic"
	"net/http"

	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func GetPlantStatistics(mongoDBInterface *mongodb.MethodInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//////////////////////////////////////////////////////
		///////// SETUP //////////////////////////////////////
		//
		neutralResponseErr := "We appologize. Statistic evaluation currently not possible due to github.com/paulmuenzner/powerplantmanager updates. Please, try again later."

		// Extracting the userId value from the claim data and converting its type to string
		claimData, err := cookie.GetCookieData(r, config.AuthCookieName)
		if err != nil {
			logger.GetLogger().Error("Cannot extract data/claim from cookie in validator 'GetPlantStatisticsValidation'. Error: ", err, "Cookie name: ", config.AuthCookieName)
			// Neutral message
			errHandler.HandleError(w, neutralResponseErr, errHandler.InternalServerError)
			return
		}
		userID, ok := claimData["data"].(map[string]interface{})["userId"].(string)
		if !ok {
			logger.GetLogger().Error("Cannot parse and access userId from 'claimData[\"data\"]' in 'GetPlantStatistics()'.")
			errHandler.HandleError(w, neutralResponseErr, errHandler.InternalServerError)
			return
		}

		//////////////////////////////////////////////
		// REQUEST BODY VALIDATION ///////////////////
		//
		// Access the parsed JSON data from the context, here req
		dataBody, ok := r.Context().Value("requestBody").(map[string]interface{})
		if !ok {
			logger.GetLogger().Error("Cannot parse and access JSON of requestBody in 'GetPlantStatistics'.")
			errHandler.HandleError(w, neutralResponseErr, errHandler.InternalServerError)
			return
		}

		///////////////////////////////////////////////////////
		// Get start and end date for period to analyze
		//////////

		// Start time period
		dateStartString, dateStartValid := dataBody["dateStart"].(string)
		if !dateStartValid {
			errHandler.HandleError(w, neutralResponseErr, errHandler.InternalServerError)
			return
		}

		// End time period
		dateEndString, dateEndValid := dataBody["dateEnd"].(string)
		if !dateEndValid {
			errHandler.HandleError(w, neutralResponseErr, errHandler.InternalServerError)
			return
		}

		// Parse the string into a time.Time object
		dateStart, err := time.Parse(time.RFC3339Nano, dateStartString)
		if err != nil {
			logger.GetLogger().Error("Cannot parse dateStartString in 'GetPlantStatistics'. Error: ", err)
			errHandler.HandleError(w, neutralResponseErr, errHandler.InternalServerError)
			return
		}

		dateEnd, err := time.Parse(time.RFC3339Nano, dateEndString)
		if err != nil {
			logger.GetLogger().Error("Cannot parse dateEndString in 'GetPlantStatistics'. Error: ", err)
			errHandler.HandleError(w, neutralResponseErr, errHandler.InternalServerError)
			return
		}

		// Get plant-related collection name
		collectionNameLogger, ok := r.Context().Value("plantCollectionNameLogger").(string)
		if !ok {
			logger.GetLogger().Error("Cannot parse and access JSON of collectionNameLogger in 'GetPlantStatistics'.")
			errHandler.HandleError(w, neutralResponseErr, errHandler.InternalServerError)
			return
		}

		//////////////////////////////////////////////
		// REQUEST BODY VALIDATION ///////////////////
		//
		// Query plant logs
		filter := bson.M{
			"created_at": bson.M{
				"$gte": dateStart,
				"$lt":  dateEnd,
			},
		}

		var plantLogger []model.PlantLogger
		var sortCriteria bson.D = bson.D{}
		err = mongoDBInterface.RepositoryInterface.FindManyInMongo(config.DatabaseNamePlantLogger, filter, collectionNameLogger, sortCriteria, &plantLogger)
		if err != nil {
			log := "User with _id " + userID + " requested non-existing plant logger for collection name " + collectionNameLogger + " in controller 'GetPlantStatistics'."
			logger.GetLogger().Error(log, err)
			errHandler.HandleError(w, "Requested plant not available.", errHandler.BadRequest)
			return
		}

		// Restructure retrieved data for statistical analysis
		solarRadiation := []float64{}
		powerOutputs := []float64{}
		for _, plant := range plantLogger {
			solarRadiation = append(solarRadiation, plant.SolarRadiation)
			powerOutputs = append(powerOutputs, plant.PowerOutput)
		}

		// Mean
		meanPowerOutputs, _ := statistic.Mean(powerOutputs)
		meanSolarRadiation, _ := statistic.Mean(solarRadiation)

		// Median
		medianPowerOutputs, _ := statistic.Median(powerOutputs, 0.5, nil)
		medianSolarRadiation, _ := statistic.Median(solarRadiation, 0.5, nil)

		// Variance
		variancePowerOutputs, _ := statistic.Variance(powerOutputs, nil)
		varianceSolarRadiation, _ := statistic.Variance(solarRadiation, nil)

		// Variance
		standardDeviationPowerOutputs, _ := statistic.StandardDeviation(powerOutputs, nil)
		standardDeviationSolarRadiation, _ := statistic.StandardDeviation(solarRadiation, nil)

		// Skewness
		skewnessPowerOutputs, _ := statistic.Skewness(powerOutputs, nil)
		skewnessSolarRadiation, _ := statistic.Skewness(solarRadiation, nil)

		// Quantile
		quantile25PowerOutputs, quantile75PowerOutputs, iqrPowerOutputs, lowerBoundPowerOutputs, upperBoundPowerOutputs, outliersPowerOutputs, quantile90PowerOutputs, quantile95PowerOutputs := statistic.Quantile(powerOutputs)
		quantile25SolarRadiation, quantile75SolarRadiation, iqrSolarRadiation, lowerBoundSolarRadiation, upperBoundSolarRadiation, outliersSolarRadiation, quantile90SolarRadiation, quantile95SolarRadiation := statistic.Quantile(solarRadiation)

		// Correlation
		correlationPowerSolar, _ := statistic.Correlation(powerOutputs, solarRadiation)

		data := map[string]interface{}{
			"powerOutput": map[string]interface{}{
				"mean":               meanPowerOutputs,
				"variance":           variancePowerOutputs,
				"median":             medianPowerOutputs,
				"standardDeviation":  standardDeviationPowerOutputs,
				"skewness":           skewnessPowerOutputs,
				"quantile25":         quantile25PowerOutputs,
				"quantile75":         quantile75PowerOutputs,
				"quantile90":         quantile90PowerOutputs,
				"quantile95":         quantile95PowerOutputs,
				"interquartileRange": iqrPowerOutputs,
				"lowerBound":         lowerBoundPowerOutputs,
				"upperBound":         upperBoundPowerOutputs,
				"outliers":           outliersPowerOutputs,
			},
			"solarRadiation": map[string]interface{}{
				"mean":               meanSolarRadiation,
				"variance":           varianceSolarRadiation,
				"median":             medianSolarRadiation,
				"standardDeviation":  standardDeviationSolarRadiation,
				"skewness":           skewnessSolarRadiation,
				"quantile25":         quantile25SolarRadiation,
				"quantile75":         quantile75SolarRadiation,
				"quantile90":         quantile90SolarRadiation,
				"quantile95":         quantile95SolarRadiation,
				"interquartileRange": iqrSolarRadiation,
				"lowerBound":         lowerBoundSolarRadiation,
				"upperBound":         upperBoundSolarRadiation,
				"outliers":           outliersSolarRadiation,
			},
			"correlationPowerSolar": correlationPowerSolar,
		}

		responsehandler.HandleSuccess(w, "Requested statistical data retrieved.", responsehandler.OK, data)

	}
}
