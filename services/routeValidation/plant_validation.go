package routevalidation

import (
	"context"
	"fmt"
	config "github.com/paulmuenzner/powerplantmanager/config"
	model "github.com/paulmuenzner/powerplantmanager/models"
	errHandler "github.com/paulmuenzner/powerplantmanager/services/errorHandler"
	arrayhandler "github.com/paulmuenzner/powerplantmanager/utils/array"
	"github.com/paulmuenzner/powerplantmanager/utils/convert"
	cookie "github.com/paulmuenzner/powerplantmanager/utils/cookies"
	crypto "github.com/paulmuenzner/powerplantmanager/utils/crypto"
	ip "github.com/paulmuenzner/powerplantmanager/utils/ip"
	logger "github.com/paulmuenzner/powerplantmanager/utils/logs"
	mongodb "github.com/paulmuenzner/powerplantmanager/utils/mongoDB"
	stringHandler "github.com/paulmuenzner/powerplantmanager/utils/strings"
	typepackage "github.com/paulmuenzner/powerplantmanager/utils/type"
	v "github.com/paulmuenzner/powerplantmanager/utils/validate"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
)

// ///////////////////////////////////////////////
// Validation middleware for registration endpoint
type Plant struct {
	Email string
}

// /////////////////////////////////////////////////////////////////////////////////////////////
// ADD ENERGY PLANT VALIDATION
// ///////////////////////
func AddPlantValidation(next http.HandlerFunc, mongoDBInterface *mongodb.MethodInterface) http.HandlerFunc {
	// Validate if plant with name alreday existing
	return func(w http.ResponseWriter, r *http.Request) {
		neutralResponseErr := "We appologize. Login currently not possible due to github.com/paulmuenzner/powerplantmanager update. Please, try again later."
		//////////////////////////////////////////////
		// VALIDATE AUTH STATUS //////////////////////
		//
		// Validate if logged in
		expired := cookie.HasCookieExpired(r, config.AuthCookieName)
		if expired {
			errHandler.HandleError(w, "You are not authenticated. Please signin.", errHandler.Unauthorized)
			return
		}

		//////////////////////////////////////////////
		// REQUEST BODY VALIDATION ///////////////////
		//
		// Access the parsed JSON data from the context
		data, ok := r.Context().Value("requestBody").(map[string]interface{})
		if !ok {
			logger.GetLogger().Error("Error in 'AddPlantValidation'. Cannot parse requestBody.", " Request: ", r)
			errHandler.HandleError(w, "Internal Server Error", errHandler.InternalServerError)
			return
		}

		// Verify number of request values
		if len(data) != 1 {
			logger.GetLogger().Warn("Not exact number of request values in 'AddPlantValidation'. Number: ", len(data), "Content: ", data)
			errHandler.HandleError(w, neutralResponseErr, errHandler.InternalServerError)
			return
		}

		// Verify and Define variables
		name, nameValid := data["name"].(string)
		if !nameValid {
			errHandler.HandleError(w, neutralResponseErr, errHandler.InternalServerError)
			return
		}

		// Validate if request body exactly contains number and names of expected keys
		expectedKeys := []string{"name"}
		validateKeys := v.Validate(data).
			HasMapExactKeys(expectedKeys).
			GetResult()

		if len(validateKeys) > 0 {
			errHandler.HandleError(w, validateKeys[0], errHandler.BadRequest)
			return
		}

		// Validate plant name
		plantNameLengthString := convert.IntToString(config.PlantNameLength)
		errorMsgName := fmt.Sprintf("Maximum number of characters for plant name cannot exceed %s.", plantNameLengthString)
		validateName := v.Validate(name).
			MaxLength(config.PlantNameLength, errorMsgName).
			GetResult()

		if len(validateName) > 0 {
			errHandler.HandleError(w, validateName[0], errHandler.BadRequest)
			return
		}

		// Check if plant name already in database. Each plant name must be unique
		exists, _ := mongoDBInterface.RepositoryInterface.IsValueInCollection(config.DatabaseNameUserAuth, config.CollectionNamePhotovoltaicPlant, "name", data["name"].(string))
		if exists {
			errHandler.HandleError(w, "Please choose another plant name.", errHandler.BadRequest)
			return
		}

		// Call the next handler if validation passes
		next.ServeHTTP(w, r)
	}
}

// /////////////////////////////////////////////////////////////////////////////////////////////
// CREATE NEW KEY AND SECRET
// /////////////////////////
func SetKeySecretValidation(next http.HandlerFunc, mongoDBInterface *mongodb.MethodInterface) http.HandlerFunc {
	// Validate if plant with name alreday existing
	return func(w http.ResponseWriter, r *http.Request) {
		neutralResponseErr := "We appologize. Configuration update currently not possible due to github.com/paulmuenzner/powerplantmanager update. Please, try again later."

		//////////////////////////////////////////////
		// VALIDATE SOURCE IP ////////////////////////
		//
		// As this is a sensitive route, it makes sense to log the IP enabling to identify and block suspicious requests and potential attackers
		clientIP, err := ip.ExtractIP(r)
		if err != nil {
			logger.GetLogger().Error("Error clientIP in 'SetKeySecretValidation()' using 'ExtractIP()'. Error: ", err)
			errHandler.HandleError(w, neutralResponseErr, errHandler.InternalServerError)
			return
		}

		//////////////////////////////////////////////
		// VALIDATE AUTH STATUS //////////////////////
		//
		// Validate if logged in
		expired := cookie.HasCookieExpired(r, config.AuthCookieName)
		if expired {
			logger.GetLogger().Errorf("Received request for key and secret update without valid cookie in 'SetKeySecretValidation()' using 'HasCookieExpired()' from IP: %s", clientIP)
			errHandler.HandleError(w, "You are not authenticated. Please signin.", errHandler.Unauthorized)
			return
		}

		//////////////////////////////////////////////
		// REQUEST BODY VALIDATION ///////////////////
		//
		// Access the parsed JSON data from the context
		data, ok := r.Context().Value("requestBody").(map[string]interface{})
		if !ok {
			logger.GetLogger().Error("Error in 'SetKeySecretValidation'. Cannot parse requestBody.", " Request: ", r)
			errHandler.HandleError(w, "Internal Server Error", errHandler.InternalServerError)
			return
		}

		// Verify number of request values
		if len(data) != 1 {
			logger.GetLogger().Warn("Not exact number of request values in 'SetKeySecretValidation'. Number: ", len(data), "Content: ", data)
			errHandler.HandleError(w, neutralResponseErr, errHandler.InternalServerError)
			return
		}

		// Define and check publicPlantID
		publicPlantID, publicPlantIdValid := data["publicPlantID"].(string)
		if !publicPlantIdValid {
			logger.GetLogger().Error("Cannot convert publicPlantID in 'SetKeySecretValidation' to string. Value: ", publicPlantID)
			errHandler.HandleError(w, neutralResponseErr, errHandler.InternalServerError)
			return
		}

		// Validate if request body exactly contains number and names of expected keys
		expectedKeys := []string{"publicPlantID"}
		validateKeys := v.Validate(data).
			HasMapExactKeys(expectedKeys).
			GetResult()

		if len(validateKeys) > 0 {
			errHandler.HandleError(w, neutralResponseErr, errHandler.BadRequest)
			return
		}

		////////////////////////////////////////////////////////////////////////////////
		// Validate if plant with publicPlantID exists and if requesting user is authorized to access and update its config
		// Extract data from JWT in cookie
		claimData, err := cookie.GetCookieData(r, config.AuthCookieName)
		if err != nil {
			logger.GetLogger().Error("Cannot extract data/claim from cookie in validator 'SetKeySecretValidation'. Error: ", err, "Cookie name: ", config.AuthCookieName)
			// Neutral message
			errHandler.HandleError(w, neutralResponseErr, errHandler.InternalServerError)
			return
		}

		// Extracting the userId value from the claim data and converting its type
		userIDRaw := claimData["data"].(map[string]interface{})["userId"]
		userID := stringHandler.InterfaceToString(userIDRaw)

		// Find plant by provided public plant id to validate if owner of plant (plantDoc["userId"]) equals _id in cookie
		var filter bson.M = bson.M{"public_plant_id": publicPlantID}
		var plant model.PhotovoltaicPlant
		var sort bson.D = bson.D{}
		database := config.DatabaseNamePlants
		collection := config.CollectionNamePhotovoltaicPlant
		findOne, err := mongoDBInterface.RepositoryInterface.FindOneInMongo(database, filter, collection, sort, &plant)
		if err != nil {
			logger.GetLogger().Errorf("Error in 'SetKeySecretValidation()' using 'FindOneInMongo()' in collection '%s' part of database '%s' finding plant with id %s. Error: %v", collection, database, publicPlantID, err)
			errHandler.HandleError(w, neutralResponseErr, errHandler.InternalServerError)
			return
		}
		if !findOne {
			logger.GetLogger().Errorf("User with id '%s' requested not existing plant with public plant id '%s' in validator 'SetKeySecretValidation'.", userID, publicPlantID)
			errHandler.HandleError(w, "Requested plant not found.", errHandler.BadRequest)
			return
		}

		// type plantRequestKey model.PhotovoltaicPlant
		// var plant plantRequestKey
		r = r.WithContext(context.WithValue(r.Context(), "plantRequest", plant))

		if plant.User.Hex() != userID {
			logger.GetLogger().Errorf("User with id '%s' requested plant id '%s' without ownership in validator 'SetKeySecretValidation'.", userID, publicPlantID)
			errHandler.HandleError(w, "You don't own any plant with your provided ID.", errHandler.BadRequest)
			return
		}

		// Call the next handler if validation passes
		next.ServeHTTP(w, r)
	}
}

// /////////////////////////////////////////////////////////////////////////////////////////////
// PLANT CONFIGURATION UPDATE
// //////////////////////////
func SetPlantConfigValidation(next http.HandlerFunc, mongoDBInterface *mongodb.MethodInterface) http.HandlerFunc {
	// Validate if plant with name alreday existing
	return func(w http.ResponseWriter, r *http.Request) {
		neutralResponseErr := "We appologize. Configuration update currently not possible due to github.com/paulmuenzner/powerplantmanager update. Please, try again later."
		//////////////////////////////////////////////
		// VALIDATE AUTH STATUS //////////////////////
		//
		// Validate if logged in
		value, hasAuthCookie := cookie.GetCookie(r, config.AuthCookieName)
		if hasAuthCookie && len(value) > 0 {
			expired := cookie.HasCookieExpired(r, config.AuthCookieName)
			if expired {
				errHandler.HandleError(w, "You are not authenticated. Please sign in.", errHandler.Unauthorized)
				return
			}
		}

		//////////////////////////////////////////////
		// REQUEST BODY VALIDATION ///////////////////
		//
		// Access the parsed JSON data from the context
		data, ok := r.Context().Value("requestBody").(map[string]interface{})
		if !ok {
			logger.GetLogger().Error("Error in 'SetPlantConfigValidation()'. Cannot parse requestBody.", " Request: ", r)
			errHandler.HandleError(w, neutralResponseErr, errHandler.InternalServerError)
			return
		}

		// Verify number of request values
		if len(data) != 3 {
			logger.GetLogger().Warn("To many request values in 'SetPlantConfigValidation()'. Number: ", len(data), "Content: ", data)
			errHandler.HandleError(w, neutralResponseErr, errHandler.InternalServerError)
			return
		}

		// Define and check publicPlantID
		publicPlantID, publicPlantIdValid := data["publicPlantID"].(string)
		if !publicPlantIdValid {
			logger.GetLogger().Error("Cannot convert publicPlantID in 'SetPlantConfigValidation()' to storing. Value: ", publicPlantID)
			errHandler.HandleError(w, neutralResponseErr, errHandler.InternalServerError)
			return
		}

		// Define and check ipWhiteList
		ipWhiteList, ipWhiteListValid := data["ipWhiteList"]
		if !ipWhiteListValid && !arrayhandler.IsEmptyArray(ipWhiteList) {
			errHandler.HandleError(w, neutralResponseErr, errHandler.InternalServerError)
			return
		}

		// Validate time interval for logging. Compare past time with minimum time which needs to past to be allowed to submit new log >> Security measure (rate limit)
		intervalSecRaw, intervalSecValid := data["intervalSec"].(float64)
		if !intervalSecValid {
			logger.GetLogger().Error("Cannot convert intervalSec in 'SetPlantConfigValidation' to float64. Value: ", intervalSecRaw)
			errHandler.HandleError(w, neutralResponseErr, errHandler.InternalServerError)
			return
		}
		intervalSec := int(intervalSecRaw)

		// Validate if request body exactly contains number and names of expected keys
		expectedKeys := []string{"publicPlantID", "ipWhiteList", "intervalSec"}
		validateKeys := v.Validate(data).
			HasMapExactKeys(expectedKeys).
			GetResult()

		if len(validateKeys) > 0 {
			errHandler.HandleError(w, neutralResponseErr, errHandler.BadRequest)
			return
		}

		// Validate interval for logging
		errorMsgName := "Interval cannot be larger than one day. Logging once a day is a minimum requirement."
		validateIntervalSec := v.Validate(intervalSec).
			MaxIntValue(60*60*24+1, errorMsgName).
			GetResult()

		if len(validateIntervalSec) > 0 {
			errHandler.HandleError(w, validateIntervalSec[0], errHandler.BadRequest)
			return
		}

		// Validate ipWhitelist
		errorMsgIPWhiteList := "Only valid IPv4 or IPv6 addresses allowed for IP white list."
		validateIPWhiteList := v.Validate(ipWhiteList).
			IsValidIPList(errorMsgIPWhiteList).
			MaxStringArrayLength(20, "To many IP addresses for white list. Max is 20.").
			GetResult()

		if len(validateIPWhiteList) > 0 {
			errHandler.HandleError(w, validateIPWhiteList[0], errHandler.BadRequest)
			return
		}

		////////////////////////////////////////////////////////////////////////////////
		// Validate if plant with publicPlantID exists and if requesting user is authorized to access and update its config
		// Extract data from JWT in cookie
		claimData, err := cookie.GetCookieData(r, config.AuthCookieName)
		if err != nil {
			logger.GetLogger().Error("Cannot extract data/claim from cookie in validator 'SetPlantConfigValidation'. Error: ", err, "Cookie name: ", config.AuthCookieName)
			// Neutral message
			errHandler.HandleError(w, neutralResponseErr, errHandler.InternalServerError)
			return
		}

		// Extracting the userId value from the claim data and converting its type
		userIDRaw := claimData["data"].(map[string]interface{})["userId"]
		userID := stringHandler.InterfaceToString(userIDRaw)

		// Find plant by provided public plant id to validate if owner of plant (plantDoc["userId"]) equals _id in cookie
		var filter bson.M = bson.M{"public_plant_id": publicPlantID}
		var plant model.PhotovoltaicPlant
		var sort bson.D = bson.D{}
		collection := config.CollectionNamePhotovoltaicPlant
		database := config.DatabaseNamePlants
		findOne, err := mongoDBInterface.RepositoryInterface.FindOneInMongo(database, filter, collection, sort, &plant)
		if err != nil {
			logger.GetLogger().Errorf("Error in 'SetPlantConfigValidation()' using 'FindOneInMongo()' in collection '%s' part of database '%s' finding plant with id %s. Error: %v", collection, database, publicPlantID, err)
			errHandler.HandleError(w, neutralResponseErr, errHandler.InternalServerError)
			return
		}
		if !findOne {
			logger.GetLogger().Errorf("User with id %s requested not existing plant with public plant id %s from collection %s and database %s in validator 'SetPlantConfigValidation()' using 'FindOneInMongo()'.", userID, publicPlantID, collection, database)
			errHandler.HandleError(w, "Requested plant not found.", errHandler.BadRequest)
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), "plantRequest", plant))

		if plant.User.Hex() != userID {
			logger.GetLogger().Errorf("User with id '%s' requested plant id '%s' without ownership in validator 'SetPlantConfigValidation'.", userID, publicPlantID)
			errHandler.HandleError(w, "You don't own any plant with your provided ID.", errHandler.BadRequest)
			return
		}

		// Call the next handler if validation passes
		next.ServeHTTP(w, r)
	}
}

// /////////////////////////////////////////////////////////////////////////////////////////////
// ADD PLANT LOG
// /////////////
func AddPlantLogValidation(next http.HandlerFunc, mongoDBInterface *mongodb.MethodInterface) http.HandlerFunc {
	// Validate if plant with name alreday existing
	return func(w http.ResponseWriter, r *http.Request) {
		neutralResponseErr := "Access is currently unavailable due to an internal github.com/paulmuenzner/powerplantmanager error. Our technical team has been notified and is actively addressing the issue."

		// No cookie validation for this route needed. Validation is realized via ip whitelist, url id and as part of request body: key and secret

		//////////////////////////////////////////////
		// REQUEST BODY VALIDATION ///////////////////
		//
		// Access the parsed JSON data from the context
		data, ok := r.Context().Value("requestBody").(map[string]interface{})
		if !ok {
			logger.GetLogger().Errorf("Error in 'AddPlantLogValidation()'. Cannot parse requestBody. Request: %+v", r)
			errHandler.HandleError(w, neutralResponseErr, errHandler.InternalServerError)
			return
		}

		//////////////////////////////////////////////
		// VALIDATE URL ID ///////////////////////////
		//
		// Access URL ID called 'apiID'
		vars := mux.Vars(r)
		apiID := vars["apiID"]

		//////////////////////////////////////////////
		// VALIDATE SOURCE IP ////////////////////////
		//
		// IP is needed to validate if it is white listed (see plant config)
		// Only IPs in white list can submit plant log
		clientIP, err := ip.ExtractIP(r)
		if err != nil {
			logger.GetLogger().Error("Error clientIP in 'AddPlantLogValidation()'. Error: ", err)
			errHandler.HandleError(w, neutralResponseErr, errHandler.Unauthorized)
			return
		}

		// Normalize compressed IP address to enable faultless comparisons of requesting IP addresses with ipWhitelist
		normalizedIP, err := ip.NormalizeIP(clientIP)
		if err != nil {
			logger.GetLogger().Error("Error normalizing clientIP in 'AddPlantLogValidation()' using 'NormalizeIP()'. Error: ", err, " IP address: ", clientIP)
			errHandler.HandleError(w, neutralResponseErr, errHandler.Unauthorized)
			return
		}

		// Verify number of request values
		if len(data) != 10 {
			logger.GetLogger().Warn("Not correct number of request values in 'AddPlantLogValidation'. Number: ", len(data), " Content: ", data)
			errHandler.HandleError(w, neutralResponseErr, errHandler.InternalServerError)
			return
		}

		// Validate if request body exactly contains number and names of expected keys
		expectedKeys := []string{"key", "secret", "voltageOutput", "currentOutput", "powerOutput", "solarRadiation", "tAmbient", "tModule", "relHumidity", "windSpeed"}
		validateKeys := v.Validate(data).
			HasMapExactKeys(expectedKeys).
			GetResult()

		if len(validateKeys) > 0 {
			errHandler.HandleError(w, neutralResponseErr, errHandler.BadRequest)
			return
		}

		// Define and check key
		key, keyValid := data["key"].(string)
		if !keyValid {
			logger.GetLogger().Error("Cannot convert key in 'AddPlantLogValidation' to string. Value: ", key)
			errHandler.HandleError(w, neutralResponseErr, errHandler.InternalServerError)
			return
		}

		// Define and check secret
		secret, secretValid := data["secret"].(string)
		if !secretValid {
			logger.GetLogger().Error("Cannot convert secret in 'AddPlantLogValidation' to string. Value: ", secret)
			errHandler.HandleError(w, neutralResponseErr, errHandler.InternalServerError)
			return
		}

		// Define and check voltageOutput
		voltageOutput, voltageOutputValid := data["voltageOutput"].(float64)
		if !voltageOutputValid {
			logger.GetLogger().Error("Cannot convert voltageOutput in 'AddPlantLogValidation' to float64. Value: ", voltageOutput)
			errHandler.HandleError(w, neutralResponseErr, errHandler.InternalServerError)
			return
		}

		// Define and check currentOutput
		currentOutput, currentOutputValid := data["currentOutput"].(float64)
		if !currentOutputValid {
			logger.GetLogger().Error("Cannot convert currentOutput in 'AddPlantLogValidation' to float64. Value: ", currentOutput)
			errHandler.HandleError(w, neutralResponseErr, errHandler.InternalServerError)
			return
		}

		// Define and check power output (powerOutput)
		powerOutput, powerOutputValid := data["powerOutput"].(float64)
		if !powerOutputValid {
			logger.GetLogger().Error("Cannot convert powerOutput in 'AddPlantLogValidation' to float64. Value: ", powerOutput)
			errHandler.HandleError(w, neutralResponseErr, errHandler.InternalServerError)
			return
		}

		// Define and check solar radiation (solarRadiation)
		solarRadiation, solarRadiationValid := data["solarRadiation"].(float64)
		if !solarRadiationValid {
			logger.GetLogger().Error("Cannot convert solarRadiation in 'AddPlantLogValidation' to float64. Value: ", solarRadiation)
			errHandler.HandleError(w, neutralResponseErr, errHandler.InternalServerError)
			return
		}

		// Define and check ambient temperature (tAmbient)
		tAmbient, tAmbientValid := data["tAmbient"].(float64)
		if !tAmbientValid {
			logger.GetLogger().Error("Cannot convert tAmbient in 'AddPlantLogValidation' to float64. Value: ", tAmbient)
			errHandler.HandleError(w, neutralResponseErr, errHandler.InternalServerError)
			return
		}

		// Define and check Temperature of the module (tModule)
		tModule, tModuleValid := data["tModule"].(float64)
		if !tModuleValid {
			logger.GetLogger().Error("Cannot convert tModule in 'AddPlantLogValidation' to float64. Value: ", tModule)
			errHandler.HandleError(w, neutralResponseErr, errHandler.InternalServerError)
			return
		}

		// Define and check relative humidity (relHumidity)
		relativeHumidity, relativeHumidityValid := data["relHumidity"].(float64)
		if !relativeHumidityValid {
			logger.GetLogger().Error("Cannot convert relHumidity in 'AddPlantLogValidation()' to float64. Value: ", relativeHumidity)
			errHandler.HandleError(w, neutralResponseErr, errHandler.InternalServerError)
			return
		}

		// Define and check wind speed (windSpeed)
		windSpeed, windSpeedValid := data["windSpeed"].(float64)
		if !windSpeedValid {
			logger.GetLogger().Error("Cannot convert windSpeed in 'AddPlantLogValidation' to float64. Value: ", windSpeed)
			errHandler.HandleError(w, neutralResponseErr, errHandler.InternalServerError)
			return
		}

		////////////////////////////////////////////////////////////////////////////////
		// Validate existence of public_plant_id and access permission
		//
		// Find plant by provided key and validate permission with ip whitelist and secret
		var filter bson.M = bson.M{"key": key}
		var plantConfig model.PlantLoggerConfig
		var sort bson.D = bson.D{}
		findOne, err := mongoDBInterface.RepositoryInterface.FindOneInMongo(config.DatabaseNamePlantLoggerConfig, filter, config.CollectionNamePlantLoggerConfig, sort, &plantConfig)
		if err != nil {
			logger.GetLogger().Errorf("Error in 'AddPlantLogValidation()' using 'FindOneInMongo()' when querying collection '%s' part of database '%s' for plant logging key %s and provided secret %s. Error: %v", config.CollectionNamePlantLoggerConfig, config.DatabaseNamePlantLoggerConfig, key, secret, err)
			errHandler.HandleError(w, neutralResponseErr, errHandler.InternalServerError)
			return
		}
		if !findOne {
			log := "User with ip " + normalizedIP + " requested non-existing plant with key " + key + " in validator 'AddPlantLogValidation'."
			logger.GetLogger().Error(log)
			errHandler.HandleError(w, "Requested plant not found or no permission.", errHandler.BadRequest)
			return
		}

		// Validate url id
		isURLIDValid := plantConfig.URLID == apiID
		if !isURLIDValid {
			logger.GetLogger().Warnf("Not valid url_id '%s' detected in 'AddPlantLogValidation()' logging plant data. Public plant id: %v", apiID, plantConfig.PublicPlantID)
			errHandler.HandleError(w, "URL not valid.", errHandler.BadRequest)
			return
		}

		// Validate secret
		isSecretValid := crypto.IsHashValid(secret, plantConfig.Secret)
		if !isSecretValid {
			log := "Request with ip " + normalizedIP + " requested existing plant with key " + key + " in validator 'AddPlantLogValidation' by providing wrong secret."
			logger.GetLogger().Error(log, err)
			errHandler.HandleError(w, "Requested plant not found or no permission.", errHandler.BadRequest)
			return
		}

		// Validate ip against white list
		isIPValid := false
		for _, ip := range plantConfig.IPWhitelist {
			if ip == normalizedIP {
				isIPValid = true
			}
		}

		if !isIPValid {
			log := "Request with not whitelisted ip " + normalizedIP + " requested existing plant with key " + key + " in validator 'AddPlantLogValidation'."
			logger.GetLogger().Error(log)
			errHandler.HandleError(w, "Requested plant not found or no permission.", errHandler.BadRequest)
			return
		}

		////////////////////////////////////////////////////////////////////////////////
		// RATE LIMIT
		//
		// Validate time interval to prevent spamming (rate limiter)
		// Get latest entry from logger
		filter = bson.M{}
		var sort2 bson.D = bson.D{{Key: "created_at", Value: -1}}
		var plantLogger model.PlantLogger
		collectionNameLogger := plantConfig.CollectionNameLogger
		findOne, err = mongoDBInterface.RepositoryInterface.FindOneInMongo(config.DatabaseNamePlantLogger, filter, collectionNameLogger, sort2, &plantLogger)
		if err != nil {
			logger.GetLogger().Errorf("Error in 'AddPlantLogValidation()' using 'FindOneInMongo()' retrieving latest entry in collection '%s' part of database '%s' for plant logging key %s and provided secret %s. Error: %v", collectionNameLogger, config.DatabaseNamePlantLogger, key, secret, err)
			errHandler.HandleError(w, neutralResponseErr, errHandler.InternalServerError)
			return
		}
		if !findOne {
			logger.GetLogger().Errorf("User with ip %s requested non-existing plant with key %s and secret %s in validator 'AddPlantLogValidation()'.", normalizedIP, key, secret)
			errHandler.HandleError(w, "Requested plant not found or no permission.", errHandler.BadRequest)
			return
		}

		logInterval := plantConfig.IntervalSec - 60 // Deduct 60 seconds as security buffer
		dateLatestEntry := plantLogger.CreatedAt
		datePastMinimum := time.Now().Add(-time.Second * time.Duration(logInterval))

		// Check if the duration is within allowed seconds
		if datePastMinimum.Before(dateLatestEntry) {
			log := "User with ip " + normalizedIP + " requested plant with key " + key + " in validator 'AddPlantLogValidation' and tried to log to often. Minimum required interval (rate limit): " + strconv.Itoa(plantConfig.IntervalSec)
			logger.GetLogger().Error(log, err)
			errHandler.HandleError(w, "No permission to save new log. Minimum time difference between logs in seconds: "+strconv.Itoa(plantConfig.IntervalSec), errHandler.BadRequest)
			return
		}

		// Attach plantConfig to context
		r = r.WithContext(context.WithValue(r.Context(), "collectionNameLogger", plantConfig.CollectionNameLogger))

		// Call the next handler if validation passes
		next.ServeHTTP(w, r)
	}
}

// /////////////////////////////////////////////////////////////////////////////////////////////
// DELETE ENERGY PLANT VALIDATION
// ///////////////////////
func DeletePlantValidation(next http.HandlerFunc, mongoDBInterface *mongodb.MethodInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		neutralResponseErr := "We appologize. Deletion currently not possible due to github.com/paulmuenzner/powerplantmanager update. Please, try again later."
		//////////////////////////////////////////////
		// VALIDATE AUTH STATUS //////////////////////
		//
		// Validate if logged in
		expired := cookie.HasCookieExpired(r, config.AuthCookieName)
		if expired {
			logger.GetLogger().Warnf("In 'DeletePlantValidation()', requesting deletion with invalid cookie.")
			errHandler.HandleError(w, "You are not authenticated. Please signin.", errHandler.Unauthorized)
			return
		}

		//////////////////////////////////////////////
		// REQUEST BODY VALIDATION ///////////////////
		//
		// Access the parsed JSON data from the context
		data, ok := r.Context().Value("requestBody").(map[string]interface{})
		if !ok {
			logger.GetLogger().Error("Error in 'DeletePlantValidation'. Cannot parse requestBody.", " Request: ", r)
			errHandler.HandleError(w, "Internal Server Error", errHandler.InternalServerError)
			return
		}

		// Verify number of request values
		if len(data) != 1 {
			logger.GetLogger().Warn("Not correct number of values in 'DeletePlantValidation'. Number: ", len(data))
			errHandler.HandleError(w, neutralResponseErr, errHandler.InternalServerError)
			return
		}

		// Verify and Define variables
		publicPlantID, publicPlantIdValid := data["publicPlantID"].(string)
		if !publicPlantIdValid {
			valueType := typepackage.GetType(data["publicPlantID"])
			logger.GetLogger().Warnf("In 'DeletePlantValidation()', not able to define valid publicPlantID. publicPlantID: %s. Type: %+v", data["publicPlantID"], valueType)
			errHandler.HandleError(w, neutralResponseErr, errHandler.InternalServerError)
			return
		}

		// Validate if request body exactly contains number and names of expected keys
		expectedKeys := []string{"publicPlantID"}
		validateKeys := v.Validate(data).HasMapExactKeys(expectedKeys).
			GetResult()

		if len(validateKeys) > 0 {
			errHandler.HandleError(w, neutralResponseErr, errHandler.BadRequest)
			return
		}

		////////////////////////////////////////////////////////////////////////////////
		// Validate if plant with publicPlantID exists and if requesting user is authorized to access and update its config
		// Extract data from JWT in cookie
		claimData, err := cookie.GetCookieData(r, config.AuthCookieName)
		if err != nil {
			logger.GetLogger().Error("Cannot extract data/claim from cookie in validator 'DeletePlantValidation'. Error: ", err, "Cookie name: ", config.AuthCookieName)
			// Neutral message
			errHandler.HandleError(w, neutralResponseErr, errHandler.InternalServerError)
			return
		}

		// Extracting the userId value from the claim data and converting its type
		userID, ok := claimData["data"].(map[string]interface{})["userId"].(string)
		if !ok {
			logger.GetLogger().Errorf("Failed to claim userId value from claim data and converting its type to string in 'DeletePlantValidation()'. Raw claim data: %+v", claimData["data"].(map[string]interface{})["userId"])
			// Neutral message
			errHandler.HandleError(w, neutralResponseErr, errHandler.InternalServerError)
			return
		}

		// Find plant logger config collection
		var filter bson.M = bson.M{"public_plant_id": publicPlantID}
		var plantLoggerConfig model.PlantLoggerConfig
		var sort bson.D = bson.D{}
		findOne, err := mongoDBInterface.RepositoryInterface.FindOneInMongo(config.DatabaseNamePlantLoggerConfig, filter, config.CollectionNamePlantLoggerConfig, sort, &plantLoggerConfig)
		if err != nil {
			logger.GetLogger().Errorf("Error in 'DeletePlantValidation()' using 'FindOneInMongo()' when querying collection '%s' part of database '%s' finding plant logger for plant ID %s by user ID %s. Error: %v", config.CollectionNamePlantLoggerConfig, config.DatabaseNamePlantLoggerConfig, publicPlantID, userID, err)
			errHandler.HandleError(w, neutralResponseErr, errHandler.InternalServerError)
			return
		}
		if !findOne {
			logger.GetLogger().Errorf("User with id '%s' requested non-existing plant logger config with public plant id '%s' in validator 'DeletePlantValidation()' using 'FindOneInMongo()'.", userID, publicPlantID)
			errHandler.HandleError(w, "Requested plant not found.", errHandler.BadRequest)
			return
		}

		// Find user_id in Plant to compare with userId in cookie
		filter = bson.M{"_id": plantLoggerConfig.ID} // prepare correctness
		var plant model.PhotovoltaicPlant
		findOne, err = mongoDBInterface.RepositoryInterface.FindOneInMongo(config.DatabaseNamePlants, filter, config.CollectionNamePhotovoltaicPlant, sort, &plant)
		if err != nil {
			logger.GetLogger().Errorf("Error in 'DeletePlantValidation()' using 'FindOneInMongo()' when querying collection '%s' part of database '%s' for plant ID %s by user ID %s. Error: %v", config.CollectionNamePhotovoltaicPlant, config.DatabaseNamePlants, publicPlantID, userID, err)
			errHandler.HandleError(w, neutralResponseErr, errHandler.InternalServerError)
			return
		}
		if !findOne {
			log := "User with id " + userID + " requested non-existing plant with public plant id " + publicPlantID + " in validator 'DeletePlantValidation'."
			logger.GetLogger().Error(log)
			errHandler.HandleError(w, "Requested plant not found.", errHandler.BadRequest)
			return
		}

		// Validate if user is owner of this requested plant
		if plant.User.Hex() != userID {
			logger.GetLogger().Errorf("User with id %s requested plant id %s without ownership in validator 'DeletePlantValidation()'. Requested plant belongs to user id %s.", userID, publicPlantID, plant.User.Hex())
			errHandler.HandleError(w, "You don't own any plant with your provided ID.", errHandler.BadRequest)
			return
		}

		// Attach JSON data of queried plant to request context needed in delete_plant_controller
		collectionNameLogger := plantLoggerConfig.CollectionNameLogger
		r = r.WithContext(context.WithValue(r.Context(), "collectionNameLogger", collectionNameLogger))

		// Call the next handler if validation passes
		next.ServeHTTP(w, r)
	}
}

// /////////////////////////////////////////////////////////////////////////////////////////////
// GET PLANT STATISTICS
// ///////////////////////
func GetPlantStatisticsValidation(next http.HandlerFunc, mongoDBInterface *mongodb.MethodInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		neutralResponseErr := "We appologize. Statistics currently not available due to github.com/paulmuenzner/powerplantmanager update. Please, try again later."
		//////////////////////////////////////////////
		// VALIDATE AUTH STATUS //////////////////////
		//
		// Validate if logged in
		expired := cookie.HasCookieExpired(r, config.AuthCookieName)
		if expired {
			errHandler.HandleError(w, "You are not authenticated. Please signin.", errHandler.Unauthorized)
			return
		}

		//////////////////////////////////////////////
		// REQUEST BODY VALIDATION ///////////////////
		//
		// Access the parsed JSON data from the context
		data, ok := r.Context().Value("requestBody").(map[string]interface{})
		if !ok {
			logger.GetLogger().Error("Error in 'GetPlantStatisticsValidation'. Cannot parse requestBody.", " Request: ", r)
			errHandler.HandleError(w, neutralResponseErr, errHandler.InternalServerError)
			return
		}

		// Verify number of request values
		if len(data) != 3 {
			logger.GetLogger().Warn("To many request values in 'GetPlantStatisticsValidation'. Number: ", len(data))
			errHandler.HandleError(w, neutralResponseErr, errHandler.InternalServerError)
			return
		}

		// Verify and Define variables
		publicPlantID, publicPlantIdValid := data["publicPlantID"].(string)
		if !publicPlantIdValid {
			errHandler.HandleError(w, neutralResponseErr, errHandler.InternalServerError)
			return
		}

		// Validate if request body exactly contains number and names of expected keys
		expectedKeys := []string{"publicPlantID", "dateStart", "dateEnd"}
		validateKeys := v.Validate(data).HasMapExactKeys(expectedKeys).
			GetResult()

		if len(validateKeys) > 0 {
			errHandler.HandleError(w, neutralResponseErr, errHandler.BadRequest)
			return
		}

		////////////////////////////////////////////////////////////////////////////////
		// Validate if plant with publicPlantID exists and if requesting user is authorized to access statistical data
		// Extract data from JWT in cookie
		claimData, err := cookie.GetCookieData(r, config.AuthCookieName)
		if err != nil {
			logger.GetLogger().Error("Cannot extract data/claim from cookie in validator 'GetPlantStatisticsValidation'. Error: ", err, "Cookie name: ", config.AuthCookieName)
			// Neutral message
			errHandler.HandleError(w, neutralResponseErr, errHandler.InternalServerError)
			return
		}

		// Extracting the userId value from the claim data and converting its type
		userID, ok := claimData["data"].(map[string]interface{})["userId"].(string)
		if !ok {
			logger.GetLogger().Warn("Unable to extract userID from claim data in 'GetPlantStatisticsValidation'. Claim data: ", claimData["data"])
			errHandler.HandleError(w, neutralResponseErr, errHandler.InternalServerError)
			return
		}

		// Find plant logger config collection
		var filterA bson.M = bson.M{"public_plant_id": publicPlantID}
		var plantLoggerConfig model.PlantLoggerConfig
		var sort bson.D = bson.D{}
		findOne, err := mongoDBInterface.RepositoryInterface.FindOneInMongo(config.DatabaseNamePlantLoggerConfig, filterA, config.CollectionNamePlantLoggerConfig, sort, &plantLoggerConfig)
		if err != nil {
			logger.GetLogger().Errorf("Error in 'GetPlantStatisticsValidation()' using 'FindOneInMongo()' when querying collection '%s' part of database '%s' for plant ID %s by user ID %s. Error: %v", config.CollectionNamePlantLoggerConfig, config.DatabaseNamePlantLoggerConfig, publicPlantID, userID, err)
			errHandler.HandleError(w, neutralResponseErr, errHandler.InternalServerError)
			return
		}
		if !findOne {
			logger.GetLogger().Errorf("User with id '%s' requested non-existing plant logger config with public plant id '%s' in validator 'GetPlantStatisticsValidation()' using 'FindOneInMongo()'.", userID, publicPlantID)
			errHandler.HandleError(w, "Requested plant not available.", errHandler.BadRequest)
			return
		}

		// Find user_id in Plant to compare with userId in cookie
		var filter bson.M = bson.M{"_id": plantLoggerConfig.ID}
		var plant model.PhotovoltaicPlant
		findOne, err = mongoDBInterface.RepositoryInterface.FindOneInMongo(config.DatabaseNamePlants, filter, config.CollectionNamePhotovoltaicPlant, sort, &plant)
		if err != nil {
			logger.GetLogger().Errorf("Error in 'GetPlantStatisticsValidation()' using 'FindOneInMongo()' when querying collection '%s' part of database '%s' for plant logger ID %s by user ID %s. Error: %v", config.CollectionNamePhotovoltaicPlant, config.DatabaseNamePlants, plantLoggerConfig.ID, userID, err)
			errHandler.HandleError(w, neutralResponseErr, errHandler.InternalServerError)
			return
		}
		if !findOne {
			log := "User with id " + userID + " requested non-existing plant with public plant id " + publicPlantID + " in validator 'GetPlantStatisticsValidation'."
			logger.GetLogger().Error(log)
			errHandler.HandleError(w, "Requested plant not found.", errHandler.BadRequest)
			return
		}

		// Validate if user is owner of this requested plant
		if plant.User.Hex() != userID {
			log := "User with id " + userID + " requested plant id " + publicPlantID + " without ownership in validator 'GetPlantStatisticsValidation'."
			logger.GetLogger().Error(log, err)
			errHandler.HandleError(w, "You don't own any plant with your provided plant id.", errHandler.BadRequest)
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), "plantCollectionNameLogger", plantLoggerConfig.CollectionNameLogger))

		// Call the next handler if validation passes
		next.ServeHTTP(w, r)
	}
}
