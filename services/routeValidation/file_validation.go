package routevalidation

import (
	"context"
	config "github.com/paulmuenzner/powerplantmanager/config"
	model "github.com/paulmuenzner/powerplantmanager/models"
	errHandler "github.com/paulmuenzner/powerplantmanager/services/errorHandler"
	"github.com/paulmuenzner/powerplantmanager/utils/aws"
	cookie "github.com/paulmuenzner/powerplantmanager/utils/cookies"
	logger "github.com/paulmuenzner/powerplantmanager/utils/logs"
	mongodb "github.com/paulmuenzner/powerplantmanager/utils/mongoDB"
	stringHandler "github.com/paulmuenzner/powerplantmanager/utils/strings"
	"net/http"
	"os"

	"github.com/disintegration/imaging"
	"go.mongodb.org/mongo-driver/bson"
)

// /////////////////////////////////////////////////////////////////////////////////////////////
// UPLOAD VALIDATION
// ///////////////////////
func UploadImageValidation(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//////////////////////////////////////////////
		// VALIDATE AUTH STATUS //////////////////////
		//
		// Validate if already logged in
		expired := cookie.HasCookieExpired(r, config.AuthCookieName)
		if expired {
			errHandler.HandleError(w, "You must sign into your account before uploading an image.", errHandler.InternalServerError)
			return
		}

		// Set limit to 0 to parse the entire request body
		err := r.ParseMultipartForm(0)
		if err != nil {
			// handle error
			http.Error(w, "Unable to parse form", http.StatusBadRequest)
			return
		}

		// Extracting the multipart form data from the HTTP request.
		formdata := r.MultipartForm

		// Extracting a slice of multipart.FileHeader objects representing the uploaded files.
		// The files slice will contain information about each uploaded file, such as its name, size, and content type.
		files := formdata.File["files"]

		// Maximum upload size 5 megabytes per file
		maxSize := int64(4 << 20)

		// Define the maximum number of files which can be uploaded at once
		if len(files) > 3 {
			errHandler.HandleError(w, "Maximum number of files in one upload is 3. If needed, please upload separatly.", errHandler.BadRequest)
			return
		}

		// Validate each file separatly
		for _, f := range files {

			// Validate permitted file type png and jpeg
			typeImg := f.Header.Get("Content-Type")
			if typeImg != "image/jpeg" && typeImg != "image/png" {
				errHandler.HandleError(w, "You can only upload png or jpeg images.", errHandler.BadRequest)
				return
			}

			// Open the uploaded file
			file, err := f.Open()
			if err != nil {
				http.Error(w, "Error opening uploaded file", http.StatusInternalServerError)
				return
			}
			defer file.Close()

			// Decode the image
			img, err := imaging.Decode(file)
			if err != nil {
				http.Error(w, "Error decoding image", http.StatusInternalServerError)
				return
			}

			/////////////////// VALIDATION ////////////////////////
			// Validate minimum width of 700 pixels per image
			if img.Bounds().Dx() < 700 {
				errHandler.HandleError(w, "Each image must have a minimum width of 700 pixels.", errHandler.BadRequest)
				return
			}

			// Validate that each file does not surpass maximum allowed file size defined in maxSize
			if f.Size > maxSize {
				errHandler.HandleError(w, "Maximum file size is 4 megabytes for each. If needed, please upload separatly.", errHandler.BadRequest)
				return
			}

		}

		// Call the next handler if validation passes
		next.ServeHTTP(w, r)
	}
}

// /////////////////////////////////////////////////////////////////////////////////////////////
// DELETE FILE VALIDATION
// ////////////////
func DeleteImageValidation(next http.HandlerFunc, awsInterface *aws.MethodInterface, mongoDBInterface *mongodb.MethodInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		neutralResponseErr := "We appologize. Deletion currently not possible due to github.com/paulmuenzner/powerplantmanager update. Please, try again later."

		//////////////////////////////////////////////
		// VALIDATE AUTH STATUS //////////////////////
		//
		// Validate if already logged in
		expired := cookie.HasCookieExpired(r, config.AuthCookieName)
		if expired {
			errHandler.HandleError(w, "You must sign into your account before deleting an image.", errHandler.InternalServerError)
			return
		}

		//////////////////////////////////////////////
		// REQUEST BODY VALIDATION ///////////////////
		//
		// Access the parsed JSON data from the context
		data, ok := r.Context().Value("requestBody").(map[string]interface{})
		if !ok {
			logger.GetLogger().Error("Error in 'DeleteImageValidation'. Cannot parse requestBody.", " Request: ", r)
			errHandler.HandleError(w, neutralResponseErr, errHandler.InternalServerError)
			return
		}

		// Verify number of request values
		if len(data) != 1 {
			logger.GetLogger().Warn("Not correct number of request values in 'DeleteImageValidation'. Number: ", len(data), "Content: ", data)
			errHandler.HandleError(w, neutralResponseErr, errHandler.InternalServerError)
			return
		}

		// Define and check publicFileID
		publicFileID, publicFileIdValid := data["publicFileID"].(string)
		if !publicFileIdValid {
			logger.GetLogger().Error("Cannot convert publicFileID to string in 'DeleteImageValidation'. Value: ", publicFileID)
			errHandler.HandleError(w, neutralResponseErr, errHandler.InternalServerError)
			return
		}

		////////////////////////////////////////////////////////////////////////////////
		// Validate if file with publicFileID exists and if requesting user is authorized to access
		// Extract data from JWT in cookie
		claimData, err := cookie.GetCookieData(r, config.AuthCookieName)
		if err != nil {
			logger.GetLogger().Error("Cannot extract data/claim from cookie in validator 'DeleteImageValidation'. Error: ", err, "Cookie name: ", config.AuthCookieName)
			// Neutral message
			errHandler.HandleError(w, neutralResponseErr, errHandler.InternalServerError)
			return
		}

		// Extracting the userId value from the claim data and converting its type
		userIDRaw := claimData["data"].(map[string]interface{})["userId"]
		userID := stringHandler.InterfaceToString(userIDRaw)

		// Find file by provided public file id to validate if owner of file (file.User) equals userID in cookie
		var filter bson.M = bson.M{"public_file_id": publicFileID}
		var file model.File
		var sort bson.D = bson.D{}
		findOne, err := mongoDBInterface.RepositoryInterface.FindOneInMongo(config.DatabaseNameFiles, filter, config.CollectionNameFiles, sort, &file)
		if err != nil {
			logger.GetLogger().Errorf("Error in 'DeleteImageValidation()' using 'FindOneInMongo()' when querying collection '%s' part of database '%s'. Error: %v", config.CollectionNameFiles, config.DatabaseNameFiles, err)
			errHandler.HandleError(w, neutralResponseErr, errHandler.InternalServerError)
			return
		}
		if !findOne {
			logger.GetLogger().Errorf("User with id '%s' requested non-existing file with public file id '%s' in validator 'DeleteImageValidation()' using 'FindOneInMongo()'.", userID, publicFileID)
			errHandler.HandleError(w, "Requested file not found.", errHandler.BadRequest)
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), "file", file))

		// Validate if user is owner of this requested file
		if file.User.Hex() != userID {
			logger.GetLogger().Errorf("User with id '%s' requested file id '%s' without ownership in validator 'DeleteImageValidation'.", userID, publicFileID)
			errHandler.HandleError(w, "You don't own any file with your provided ID.", errHandler.BadRequest)
			return
		}

		/////////////////////////////////////
		// Validate if file exists on S3
		var key = file.Slug
		fileExists, err := awsInterface.RepositoryInterfaceS3.S3ObjectExists(os.Getenv("BUCKET_NAME"), key)
		if !fileExists {
			logger.GetLogger().Error("Requested file not found on S3 in 'DeleteImageValidation'.", " Public file ID: ", file.PublicFileID)
			errHandler.HandleError(w, neutralResponseErr, errHandler.InternalServerError)
			return
		}
		if err != nil {
			logger.GetLogger().Error("Cannot validate if file exists on S3 in 'DeleteImageValidation' using 'CheckFileExistsS3'. Error: ", err, " Public file ID: ", file.PublicFileID)
			errHandler.HandleError(w, neutralResponseErr, errHandler.InternalServerError)
			return
		}

		//////////////////////////////////////////////
		// PASS VALID REQUEST ////////////////////////
		//
		// Call the next handler if validation passes
		next.ServeHTTP(w, r)
	}
}
