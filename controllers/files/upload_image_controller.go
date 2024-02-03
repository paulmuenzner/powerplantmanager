package filecontroller

import (
	config "github.com/paulmuenzner/powerplantmanager/config"
	model "github.com/paulmuenzner/powerplantmanager/models"
	errHandler "github.com/paulmuenzner/powerplantmanager/services/errorHandler"
	responsehandler "github.com/paulmuenzner/powerplantmanager/services/responseHandler"
	"github.com/paulmuenzner/powerplantmanager/utils/aws"
	cookie "github.com/paulmuenzner/powerplantmanager/utils/cookies"
	data "github.com/paulmuenzner/powerplantmanager/utils/data"
	fileHandler "github.com/paulmuenzner/powerplantmanager/utils/files"
	logger "github.com/paulmuenzner/powerplantmanager/utils/logs"
	mongodb "github.com/paulmuenzner/powerplantmanager/utils/mongoDB"
	stringHandler "github.com/paulmuenzner/powerplantmanager/utils/strings"
	"net/http"
	"os"
	"time"

	"github.com/disintegration/imaging"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Image represents the image data and metadata
type Image struct {
	Name        string
	ContentType string
	Size        int64
	Data        []byte
}

func UploadImage(awsInterface *aws.MethodInterface, mongoDBInterface *mongodb.MethodInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		neutralResponseErr := "Access currently not possible due to internal github.com/paulmuenzner/powerplantmanager error. Our technical team is informed and working on it."

		// Extract user id from JWT in cookie
		claimData, err := cookie.GetCookieData(r, config.AuthCookieName)
		if err != nil {
			logger.GetLogger().Errorf("Cannot extract data/claim from cookie in 'UploadImage()' using 'GetCookieData()'. Cookie name: %s. Error: %v", config.AuthCookieName, err)
			// Neutral message
			errHandler.HandleError(w, neutralResponseErr, errHandler.InternalServerError)
			return
		}
		// Extracting the userId value from the claim data and converting its type
		userIDRaw := claimData["data"].(map[string]interface{})["userId"]
		userIDStr, ok := userIDRaw.(string)
		if !ok {
			logger.GetLogger().Errorf("Failed to extract the userId (userIDRaw) value from the claim data and converting its type in 'UploadImage()'. Raw claim data: %+v", userIDRaw)
			// Neutral message
			errHandler.HandleError(w, neutralResponseErr, errHandler.InternalServerError)
			return
		}

		userID, err := primitive.ObjectIDFromHex(userIDStr)
		if err != nil {
			logger.GetLogger().Errorf("Failed to convert hex value as string to ObjectID in 'UploadImage()' using 'ObjectIDFromHex()'. Hex value user id: %s. Error: %v", userIDStr, err)
			// Neutral message
			errHandler.HandleError(w, neutralResponseErr, errHandler.InternalServerError)
			return
		}

		// Parse form data
		err = r.ParseMultipartForm(13 << 20)
		if err != nil {
			logger.GetLogger().Errorf("Error occurred when parsing multipart form data in 'UploadImage()' with 'ParseMultipartForm()'. Error: %v. Request: %+v", err, r)
			errHandler.HandleError(w, neutralResponseErr, errHandler.InternalServerError)
			return
		}

		formdata := r.MultipartForm

		files := formdata.File["files"]

		for _, f := range files {
			// Open the uploaded file
			file, err := f.Open()
			if err != nil {
				logger.GetLogger().Errorf("Error opening uploaded file in 'UploadImage()' with 'f.Open()'. Error: %v. Request: %+v", err, r)
				errHandler.HandleError(w, neutralResponseErr, errHandler.InternalServerError)
				return
			}
			defer file.Close()

			// Decode the image
			img, err := imaging.Decode(file)
			if err != nil {
				logger.GetLogger().Errorf("Error decoding image in 'UploadImage()' using 'Decode()'. Error: %v. Request: %+v", err, r)
				errHandler.HandleError(w, neutralResponseErr, errHandler.InternalServerError)
				return
			}

			// Resize the image
			resizedImage := imaging.Resize(img, 700, 0, imaging.Lanczos)

			// Convert *image.NRGBA to byte slice
			resizedImageByteSlice, err := fileHandler.GolangImageToByteSlice(resizedImage, f.Header.Get("Content-Type"))
			if err != nil {
				logger.GetLogger().Errorf("Error converting golang image type to ByteSlice in 'UploadImage' using 'GolangImageToByteSlice()'. Error: %v. Request: %+v", err, r)
				errHandler.HandleError(w, neutralResponseErr, errHandler.InternalServerError)
				return
			}

			// Get image meta data from byte slice
			width, height, fileExtension, fileSize, _ := fileHandler.GetMetaFromImageByteSlice(resizedImageByteSlice)
			fileFolder := "golang"
			slug := fileFolder + "/" + f.Filename

			// Upload the file to S3
			err = awsInterface.RepositoryInterfaceS3.UploadFile(os.Getenv("BUCKET_NAME"), slug, resizedImageByteSlice)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				logger.GetLogger().Error("Error occurred when uploading image in 'UploadImage' using 'UploadFile()' to S3. Error: ", err, " Request: ", r)
				return
			}

			///////// ADD FILE TO DATABASE //////////
			// Prepare data to save new file document
			objectIDFile := primitive.NewObjectID() // Same Mongo Object ID for plant collection (PhotovoltaicPlant) and _id in PlantLoggerConfig
			publicFileID := stringHandler.GenerateRandomNumericString(15)
			fileDocument := model.File{
				ID:           objectIDFile,
				PublicFileID: publicFileID,
				Name:         f.Filename,
				Slug:         slug,
				Type:         fileExtension,
				User:         userID,
				Size:         fileSize,
				Width:        width,
				Height:       height,
				Folder:       fileFolder,
				CreatedAt:    time.Now(),
			}

			// Validate data against mongodb File model
			if err := data.ValidateStruct(fileDocument); err != nil {
				logger.GetLogger().Errorf("Data validation using 'ValidateStruct()' against mongodb file model failed in 'UploadImage()'. File document: %+v. Error:  %v", fileDocument, err)
				errHandler.HandleError(w, neutralResponseErr, errHandler.InternalServerError)
				return
			}

			// Save data
			_, err = mongoDBInterface.RepositoryInterface.InsertOneToMongo(config.DatabaseNameFiles, fileDocument, config.CollectionNameFiles)
			if err != nil {
				logger.GetLogger().Errorf("Unable to save file data in 'UploadImage()'. Database name: %s. File document: %+v. Collection name: %s. Error: %v", config.DatabaseNameFiles, fileDocument, config.CollectionNameFiles, err)
				errHandler.HandleError(w, neutralResponseErr, errHandler.InternalServerError)
				return
			}

		}

		responsehandler.HandleSuccess(w, "File successfully uploaded to S3", responsehandler.OK)

		return

	}
}
