package server

import (
	"net/http"

	config "github.com/paulmuenzner/powerplantmanager/config"
	"github.com/paulmuenzner/powerplantmanager/utils/aws"
	"github.com/paulmuenzner/powerplantmanager/utils/email"
	mongodb "github.com/paulmuenzner/powerplantmanager/utils/mongoDB"

	"github.com/gorilla/mux"
)

// CreateSubrouter creates a subrouter for a given path and attaches it to the main router.
type AppConfig = config.AppConfig

func CreateSubrouter(router *mux.Router, path string, subrouterFunc func(awsInterface *aws.MethodInterface, emailInterface *email.RepositoryInterface, mongoDBInterface *mongodb.MethodInterface) *mux.Router, awsInterface *aws.MethodInterface, emailInterface *email.RepositoryInterface, mongoDBInterface *mongodb.MethodInterface) {
	subrouter := subrouterFunc(awsInterface, emailInterface, mongoDBInterface)
	router.PathPrefix(path).Handler(http.StripPrefix(path, subrouter))
}
