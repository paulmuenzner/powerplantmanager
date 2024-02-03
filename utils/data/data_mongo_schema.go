package data

import (
	config "github.com/paulmuenzner/powerplantmanager/config"

	"gopkg.in/go-playground/validator.v9"
)

const (
	UserAuthDatabase          string = config.DatabaseNameUserAuth
	PlantDatabase             string = config.DatabaseNamePlants
	PlantLoggerDatabase       string = config.DatabaseNamePlantLogger
	PlantLoggerConfigDatabase string = config.DatabaseNamePlantLoggerConfig
)

// ////////////////////////////////////////
// VALIDATOR /////////////////////////////
func ValidateStruct(data interface{}) error {
	validate := validator.New()

	if err := validate.Struct(data); err != nil {
		return err
	}

	return nil
}
