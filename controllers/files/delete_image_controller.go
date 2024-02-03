package filecontroller

import (
	"net/http"
	"os"

	config "github.com/paulmuenzner/powerplantmanager/config"
	model "github.com/paulmuenzner/powerplantmanager/models"
	errHandler "github.com/paulmuenzner/powerplantmanager/services/errorHandler"
	responsehandler "github.com/paulmuenzner/powerplantmanager/services/responseHandler"
	"github.com/paulmuenzner/powerplantmanager/utils/aws"
	logger "github.com/paulmuenzner/powerplantmanager/utils/logs"
	mongodb "github.com/paulmuenzner/powerplantmanager/utils/mongoDB"

	"go.mongodb.org/mongo-driver/bson"
)

func DeleteImage(awsInterface *aws.MethodInterface, mongoDBInterface *mongodb.MethodInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//////////////////////////////////////////////////////
		///////// SETUP //////////////////////////////////////
		//
		neutralResponseErr := "We appologize. Deletion currently not possible due to github.com/paulmuenzner/powerplantmanager update. Please, try again later."

		// Access the parsed JSON data from the context
		file, ok := r.Context().Value("file").(model.File)
		if !ok {
			logger.GetLogger().Error("Cannot parse and access JSON in DeleteImage.")
			errHandler.HandleError(w, neutralResponseErr, errHandler.InternalServerError)
			return
		}

		/////////////////////////////////////////////////////////////////
		// Delete file on S3
		var keys = []string{file.Slug}
		err := awsInterface.RepositoryInterfaceS3.DeleteObjects(os.Getenv("BUCKET_NAME"), keys) // Prepare dependency injection
		if err != nil {
			errHandler.HandleError(w, neutralResponseErr, errHandler.InternalServerError)
			logger.GetLogger().Error("Error deleting S3 image in 'DeleteImage()' using 'DeleteFileFromS3()'. Error: ", err, " Public file ID: ", file.PublicFileID)
			return
		}

		/////////////////////////////////////////////////////////////////
		// DELETE DOCUMENT IN File collection
		var filterDeleteDocument bson.M = bson.M{"public_file_id": file.PublicFileID}
		_, err = mongoDBInterface.RepositoryInterface.DeleteDocumentMongo(config.DatabaseNameFiles, filterDeleteDocument, config.CollectionNameFiles)
		if err != nil {
			logger.GetLogger().Error("Unable to delete document in File collection in 'DeleteImage'. Error: ", err, " Public file ID: ", file.PublicFileID)
			errHandler.HandleError(w, neutralResponseErr, errHandler.InternalServerError)
			return
		}

		responsehandler.HandleSuccess(w, "Image deleted.", responsehandler.OK)
		return

	}
}
