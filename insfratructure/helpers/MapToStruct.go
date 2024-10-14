package helpers

import (
	utils "github.com/flabio/safe_constants"
	"github.com/safe_msvc_city/usecase/dto"
)

func MapToStruct(dataDto *dto.CityDTO, dataMap map[string]interface{}) {
	city := dto.CityDTO{
		Name:   dataMap[utils.NAME].(string),
		Active: dataMap[utils.ACTIVE].(bool),
	}
	*dataDto = city
}

func ValidateFieldCity(value map[string]interface{}) string {
	var msg string = utils.EMPTY
	if value[utils.NAME] == nil {
		msg = utils.NAME_FIELD_IS_REQUIRED
	}
	if value[utils.ACTIVE] == nil {
		msg = utils.ACTIVE_FIELD_IS_REQUIRED
	}
	return msg
}

func ValidateRequiredCity(field dto.CityDTO) string {
	var msg string = utils.EMPTY
	if field.Name == utils.EMPTY {
		msg = utils.NAME_IS_REQUIRED
	}
	return msg
}
