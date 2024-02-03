package authcontroller

import (
	"net/http"

	cookiehandler "github.com/paulmuenzner/powerplantmanager/services/cookieHandler"
	errHandler "github.com/paulmuenzner/powerplantmanager/services/errorHandler"
	responsehandler "github.com/paulmuenzner/powerplantmanager/services/responseHandler"
	logger "github.com/paulmuenzner/powerplantmanager/utils/logs"
)

func Signout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		err := cookiehandler.Signout(w)
		if err != nil {
			logger.GetLogger().Error("Error with setting a cookie in 'Login'. Error: ", err)
			errHandler.HandleError(w, "Currently signout not possible.", errHandler.InternalServerError)
			return
		}
		responsehandler.HandleSuccess(w, "See you soon.", responsehandler.Accepted)

	}
}
