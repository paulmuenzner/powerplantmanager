package routes

import (
	authcontroller "github.com/paulmuenzner/powerplantmanager/controllers/auth"
	error "github.com/paulmuenzner/powerplantmanager/services/errorHandler"
	v "github.com/paulmuenzner/powerplantmanager/services/routevalidation"
	"github.com/paulmuenzner/powerplantmanager/utils/aws"
	"github.com/paulmuenzner/powerplantmanager/utils/email"
	mongodb "github.com/paulmuenzner/powerplantmanager/utils/mongoDB"
	"net/http"

	"github.com/gorilla/mux"
)

func CreateAuthSubrouter(awsInterface *aws.MethodInterface, emailInterface *email.RepositoryInterface, mongoDBInterface *mongodb.MethodInterface) *mux.Router {
	authRouter := mux.NewRouter()

	// Sub-routes
	authRouter.HandleFunc("/registration", v.RegistrationValidation(authcontroller.Registration(emailInterface, mongoDBInterface), mongoDBInterface)).Methods("POST").Name("Registration")
	authRouter.HandleFunc("/verify", v.VerificationValidation(authcontroller.RegistrationVerify(emailInterface, mongoDBInterface))).Methods("POST").Name("RegistrationVerify")
	authRouter.HandleFunc("/signin", v.SigninValidation(authcontroller.Signin(emailInterface, mongoDBInterface))).Methods("POST").Name("Signin")
	authRouter.HandleFunc("/signout", v.SignoutValidation(authcontroller.Signout())).Methods("POST").Name("Signout")

	// Set a custom NotFoundHandler
	authRouter.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		error.HandleError(w, "Not found!", error.NotFound)
	})

	return authRouter
}
