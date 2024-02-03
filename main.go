package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	config "github.com/paulmuenzner/powerplantmanager/config"
	routes "github.com/paulmuenzner/powerplantmanager/routes"
	errorHandler "github.com/paulmuenzner/powerplantmanager/services/errorHandler"
	aws "github.com/paulmuenzner/powerplantmanager/utils/aws"
	emailHandler "github.com/paulmuenzner/powerplantmanager/utils/email"
	env "github.com/paulmuenzner/powerplantmanager/utils/env"
	ip "github.com/paulmuenzner/powerplantmanager/utils/ip"
	logger "github.com/paulmuenzner/powerplantmanager/utils/logs"
	mongoDB "github.com/paulmuenzner/powerplantmanager/utils/mongoDB"
	mongodb "github.com/paulmuenzner/powerplantmanager/utils/mongoDB"
	serverConfig "github.com/paulmuenzner/powerplantmanager/utils/server"

	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
)

type AppConfig struct {
	MongoClient *mongo.Client
}

func main() {
	router := mux.NewRouter()
	router.Use(ip.RealIP)

	///////////////////////////////////////////////
	// ENV VARIABLES //////////////////////////////
	///////////////////////////////////////////////

	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		logger.GetLogger().Errorf("Cannot load environment variables from .env file in 'main.go'. Default value used. Error: %v", err)
		return
	}

	// Initialize environment variables with helper function
	port, err := env.GetEnvValue("PORT", "8080")
	if err != nil {
		logger.GetLogger().Warnf("Cannot retrieve .env value for PORT in 'main.go'. Default value used. Error: %v", err)
	}

	///////////////////////////////////////////////
	// END ENV VARIABLES //////////////////////////
	///////////////////////////////////////////////

	///////////////////////////////////////////////
	// PROTECTION /////////////////////////////////
	///////////////////////////////////////////////

	// Create a rate limiter with a limit of 3 requests per second
	lim := tollbooth.NewLimiter(3, &limiter.ExpirableOptions{DefaultExpirationTTL: time.Hour})
	router.Use(serverConfig.TollboothMiddleware(lim))

	// Apply the custom middleware to check request methods with the allowed methods parameter
	allowedMethods := map[string]bool{"GET": true, "POST": true, "PUT": true, "DELETE": true}
	router.Use(serverConfig.CheckAllowedMethodsMiddleware(allowedMethods))

	///////////////////////////////////////////////
	// END PROTECTION /////////////////////////////
	///////////////////////////////////////////////

	///////////////////////////////////////////////
	// PARSER COOKIE LOGGER MIDDLEWARE ////////////
	///////////////////////////////////////////////

	// Body parser
	router.Use(serverConfig.BodyRequestParser)

	// Cookie parser
	router.Use(serverConfig.CookieParser)

	// Logger
	logFileName := logger.GetLogFileName()
	logger.Init(logFileName)

	///////////////////////////////////////////////
	// END PARSER COOKIE LOGGER MIDDLEWARE ////////
	///////////////////////////////////////////////

	///////////////////////////////////////////////
	// CONNECT DATABASE MONGODB ///////////////////
	///////////////////////////////////////////////

	// Get Uniform Resource Identifier
	mongodbURI, err := mongoDB.ClientConfig()
	if err != nil {
		logger.GetLogger().Warnf("Cannot retrieve .env value for Mongo URI in 'main.go'. Default value used. Error: %v", err)
	}

	// Connect to database
	client, err := mongoDB.ConnectToMongoDB(mongodbURI)
	if err != nil {
		logger.GetLogger().Errorf("Connection with MongoDB failed due to following error: %v", err)
		return
	}

	// Provide interface for database repository
	mongoDBInterface := mongodb.NewMongoDBMethodInterface(client)

	// Disconnect from MongoDB
	defer func() {
		if err := client.MongoDB.Disconnect(context.TODO()); err != nil {
			logger.GetLogger().Errorf("Disconnection MongoDB failed due to following error: %v", err)
			return
		}
	}()

	///////////////////////////////////////////////
	// END CONNECT DATABASE MONGODB ///////////////
	///////////////////////////////////////////////

	///////////////////////////////////////////////
	// PRODUCTION CONFIG //////////////////////////
	///////////////////////////////////////////////

	// Email client config production
	emailClientConfig, err := emailHandler.ProductionConfig()
	if err != nil {
		logger.GetLogger().Error("Error in 'main()' utilizing 'ProductionConfig()' retrieving emailClientConfig. Error: ", err)
		return
	}
	emailInterface, err := emailHandler.GetEmailRepositoryInterface(emailClientConfig)
	if err != nil {
		logger.GetLogger().Error("Error in 'main()' utilizing 'GetEmailRepositoryInterface()'. Cannot create 'emailInterface'. Error: ", err)
		return
	}

	// AWS client config production
	awsClientConfig, _, err := aws.S3ProductionConfig()
	if err != nil {
		logger.GetLogger().Error("Error in 'main()' utilizing 'S3ProductionConfig()' retrieving awsClientConfig. Error: ", err)
		return
	}
	awsInterface, err := aws.GetAwsMethods(awsClientConfig)
	if err != nil {
		logger.GetLogger().Error("Error in 'main()' utilizing 'GetAwsMethods()'. Cannot create 'awsInterface'. Error: ", err)
		return
	}

	///////////////////////////////////////////////
	// END PRODUCTION CONFIG //////////////////////
	///////////////////////////////////////////////

	///////////////////////////////////////////////
	// ROUTING ////////////////////////////////////
	///////////////////////////////////////////////

	// SUBROUTER
	// Auth
	serverConfig.CreateSubrouter(router, "/auth", routes.CreateAuthSubrouter, awsInterface, emailInterface, mongoDBInterface)

	// Files
	serverConfig.CreateSubrouter(router, "/files", routes.CreateFileSubrouter, awsInterface, emailInterface, mongoDBInterface)

	// Power plants
	serverConfig.CreateSubrouter(router, "/plants", routes.CreatePlantsSubrouter, awsInterface, emailInterface, mongoDBInterface)

	// Set a custom NotFoundHandler
	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Requested route: %s", r.URL.Path)
		errorHandler.HandleError(w, "URL not found!", errorHandler.NotFound)
	})

	// Start the github.com/paulmuenzner/powerplantmanager
	http.Handle("/", router)
	///////////////////////////////////////////////
	// END ROUTING ////////////////////////////////
	///////////////////////////////////////////////

	///////////////////////////////////////////////
	// SERVER CONFIG & START //////////////////////
	///////////////////////////////////////////////
	server := &http.Server{
		Addr:         "0.0.0.0:" + port,                                // Specifies the network address on which the github.com/paulmuenzner/powerplantmanager should listen. In this case, it's set to listen on all available network interfaces on port 8080.
		WriteTimeout: time.Second * time.Duration(config.WriteTimeout), // The maximum duration allowed for writing response headers to the client.
		ReadTimeout:  time.Second * time.Duration(config.ReadTimeout),  // The maximum duration allowed for reading the entire request, including the body.
		IdleTimeout:  time.Second * time.Duration(config.IdleTimeout),  // The maximum amount of time to wait for the next request when keep-alives are enabled.
		Handler:      router,                                           // The handler to invoke for each incoming request. In this case, it's set to the Gorilla Mux router (`router`).
	}

	fmt.Printf("Server started! Open http://localhost:%s\n", port)
	err = server.ListenAndServe()
	if err != nil {
		logger.GetLogger().Error("Servere error: ", err)
		log.Fatal(err)
	}

}
