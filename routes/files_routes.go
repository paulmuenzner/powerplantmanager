package routes

import (
	filecontroller "github.com/paulmuenzner/powerplantmanager/controllers/files"
	error "github.com/paulmuenzner/powerplantmanager/services/errorHandler"
	v "github.com/paulmuenzner/powerplantmanager/services/routevalidation"
	"github.com/paulmuenzner/powerplantmanager/utils/aws"
	"github.com/paulmuenzner/powerplantmanager/utils/email"
	mongodb "github.com/paulmuenzner/powerplantmanager/utils/mongoDB"
	"net/http"

	"github.com/gorilla/mux"
)

func CreateFileSubrouter(awsInterface *aws.MethodInterface, emailInterface *email.RepositoryInterface, mongoDBInterface *mongodb.MethodInterface) *mux.Router {
	filesRouter := mux.NewRouter()

	// Sub-routes
	filesRouter.HandleFunc("/upload-file", v.UploadImageValidation(filecontroller.UploadImage(awsInterface, mongoDBInterface))).Methods("POST").Name("UploadFile")
	filesRouter.HandleFunc("/delete-file", v.DeleteImageValidation(filecontroller.DeleteImage(awsInterface, mongoDBInterface), awsInterface, mongoDBInterface)).Methods("DELETE").Name("DeleteFile")

	// Set a custom NotFoundHandler
	filesRouter.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		error.HandleError(w, "Not found!", error.NotFound)
	})

	return filesRouter
}
