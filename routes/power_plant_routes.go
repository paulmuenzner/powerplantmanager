package routes

import (
	plantcontroller "github.com/paulmuenzner/powerplantmanager/controllers/plants"
	error "github.com/paulmuenzner/powerplantmanager/services/errorHandler"
	v "github.com/paulmuenzner/powerplantmanager/services/routevalidation"
	"github.com/paulmuenzner/powerplantmanager/utils/aws"
	"github.com/paulmuenzner/powerplantmanager/utils/email"
	mongodb "github.com/paulmuenzner/powerplantmanager/utils/mongoDB"
	"net/http"

	"github.com/gorilla/mux"
)

func CreatePlantsSubrouter(awsInterface *aws.MethodInterface, emailInterface *email.RepositoryInterface, mongoDBInterface *mongodb.MethodInterface) *mux.Router {
	plantRouter := mux.NewRouter()

	// Sub-routes
	plantRouter.HandleFunc("/add", v.AddPlantValidation(plantcontroller.AddPlant(mongoDBInterface), mongoDBInterface)).Methods("POST").Name("AddPlant")
	plantRouter.HandleFunc("/log/{apiID:[0-9]+}", v.AddPlantLogValidation(plantcontroller.AddLogEntry(mongoDBInterface), mongoDBInterface)).Methods("POST").Name("AddLog")
	plantRouter.HandleFunc("/setconfig", v.SetPlantConfigValidation(plantcontroller.SetPlantConfig(mongoDBInterface), mongoDBInterface)).Methods("PUT").Name("SetPlantConfig")
	plantRouter.HandleFunc("/keysecret", v.SetKeySecretValidation(plantcontroller.SetKeySecret(mongoDBInterface), mongoDBInterface)).Methods("PUT").Name("SetKeySecret")
	plantRouter.HandleFunc("/delete", v.DeletePlantValidation(plantcontroller.DeletePlant(mongoDBInterface), mongoDBInterface)).Methods("DELETE").Name("DeletePlant")
	plantRouter.HandleFunc("/statistics", v.GetPlantStatisticsValidation(plantcontroller.GetPlantStatistics(mongoDBInterface), mongoDBInterface)).Methods("Get").Name("GetStatistics")

	// Set a custom NotFoundHandler
	plantRouter.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		error.HandleError(w, "Not found!", error.NotFound)
	})

	return plantRouter
}
