package service

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	constants "github.com/flabio/safe_constants"
	"github.com/gofiber/fiber/v2"
	"github.com/safe_msvc_city/core"
	"github.com/safe_msvc_city/insfratructure/entities"
	"github.com/safe_msvc_city/insfratructure/helpers"
	"github.com/safe_msvc_city/insfratructure/ui/global"
	"github.com/safe_msvc_city/insfratructure/ui/uicore"
	"github.com/safe_msvc_city/usecase/dto"
	"github.com/ulule/deepcopier"
)

type statesService struct {
	states uicore.UIStatesCore
}

func NewSatatesService() global.UIStates {
	return &statesService{
		states: core.GetStatesInstance(),
	}

}

func (s *statesService) GetStatesFindAll(c *fiber.Ctx) error {
	result, err := s.states.GetStatesFindAll()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			constants.STATUS:  fiber.StatusBadRequest,
			constants.MESSAGE: constants.ERROR_QUERY,
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		constants.STATUS: fiber.StatusOK,
		constants.DATA:   result,
	})
}

func (s *statesService) GetStatesFindById(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params(constants.ID))
	result, err := s.states.GetStatesFindById(uint(id))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			constants.STATUS:  fiber.StatusBadRequest,
			constants.MESSAGE: constants.ERROR_QUERY,
		})
	}
	if result.Id == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			constants.STATUS: fiber.StatusNotFound,
		})
	}
	return c.Status(http.StatusOK).JSON(result)
}
func (s *statesService) GetStatesFindByIdOfCity(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params(constants.ID))

	result, err := s.states.GetStatesFindByIdOfCity(uint(id))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			constants.STATUS:  fiber.StatusBadRequest,
			constants.MESSAGE: constants.ERROR_QUERY,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		constants.STATUS: fiber.StatusOK,
		constants.DATA:   result,
	})
}
func (s *statesService) CreateState(c *fiber.Ctx) error {
	var states entities.States
	stateDto, msgError := validateState(0, s, c)
	if msgError != "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			constants.STATUS:  fiber.StatusBadRequest,
			constants.MESSAGE: msgError,
		})
	}
	deepcopier.Copy(stateDto).To(&states)
	result, err := s.states.CreateStates(states)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			constants.STATUS:  fiber.StatusInternalServerError,
			constants.MESSAGE: constants.ERROR_CREATE,
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		constants.STATUS:  fiber.StatusCreated,
		constants.DATA:    result,
		constants.MESSAGE: constants.CREATED,
	})
}

func (s *statesService) UpdateState(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params(constants.ID))
	stateDto, msgError := validateState(uint(id), s, c)
	if msgError != "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			constants.STATUS:  http.StatusBadRequest,
			constants.MESSAGE: msgError,
		})
	}
	state, _ := s.states.GetStatesFindById(uint(id))
	if state.Id == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			constants.STATUS:  http.StatusNotFound,
			constants.MESSAGE: constants.ID_NO_EXIST,
		})
	}
	deepcopier.Copy(stateDto).To(&state)
	result, err := s.states.UpdateStates(uint(id), state)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			constants.STATUS:  fiber.StatusInternalServerError,
			constants.MESSAGE: err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		constants.STATUS:  fiber.StatusOK,
		constants.DATA:    result,
		constants.MESSAGE: constants.UPDATED,
	})
}
func (s *statesService) DeleteState(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params(constants.ID))
	state, _ := s.states.GetStatesFindById(uint(id))
	if state.Id == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			constants.STATUS:  http.StatusNotFound,
			constants.MESSAGE: constants.ID_NO_EXIST,
		})
	}
	result, err := s.states.DeleteStates(uint(id))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			constants.STATUS:  fiber.StatusInternalServerError,
			constants.MESSAGE: constants.ERROR_DELETE,
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		constants.STATUS:  fiber.StatusOK,
		constants.DATA:    result,
		constants.MESSAGE: constants.REMOVED,
	})
}

func validateState(id uint, s *statesService, c *fiber.Ctx) (dto.StatesDTO, string) {
	defer func() {
		if r := recover(); r != nil {
			log.Println(constants.RECOVER_PANIC, r)
		}
	}()
	var msg string = ""
	stateDto := dto.StatesDTO{}
	body := c.Body()
	//dataMap := make(map[string]string)

	// fields := []string{
	// 	constants.NAME,
	// 	constants.CITY_ID,
	// 	constants.ACTIVE,
	// 	constants.ZIP_CODE,
	// }

	// for _, field := range fields {
	// 	value := c.FormValue(field)

	// 	if value != "" {
	// 		dataMap[field] = value
	// 	} else {
	// 		dataMap[field] = ""
	// 	}
	// }
	var dataMap map[string]interface{}
	errJson := json.Unmarshal([]byte(body), &dataMap)
	if errJson != nil {
		msg = errJson.Error()
	}

	msgValid := validateField(dataMap)
	if msgValid != "" {
		return dto.StatesDTO{}, msgValid
	}

	helpers.MapToStructState(dataMap, &stateDto)
	msg = validateRequired(stateDto)
	if msg != "" {
		return dto.StatesDTO{}, msg
	}
	existName, _ := s.states.GetStatesFindByName(id, stateDto.Name)
	if existName {
		msg = constants.NAME_ALREADY_EXIST
	}
	return stateDto, msg
}

// func MapToStructStates(stateDto *dto.StatesDTO, dataMap map[string]string) {
// 	cityId, _ := strconv.Atoi(dataMap[constants.CITY_ID])
// 	isActive, _ := strconv.ParseBool(dataMap[constants.ACTIVE])

// 	state := dto.StatesDTO{
// 		Name:    dataMap[constants.NAME],
// 		CityId:  uint(cityId),
// 		ZipCode: dataMap[constants.ZIP_CODE],
// 		Active:  isActive,
// 	}
// 	*stateDto = state
// }

func validateField(dataMap map[string]interface{}) string {
	msg := ""
	// for field, value := range dataMap {
	// 	switch field {
	// 	case constants.NAME:
	// 		if value == "" {
	// 			msg = constants.NAME_FIELD_IS_REQUIRED
	// 		}
	// 	case constants.CITY_ID:
	// 		cityId, err := strconv.Atoi(value)
	// 		if err != nil || cityId == 0 {
	// 			msg = constants.CITY_ID_FIELD_IS_REQUIRED
	// 		}
	// 	case constants.ACTIVE:
	// 		isActive, err := strconv.ParseBool(value)
	// 		if err != nil || !isActive && !(!isActive) {
	// 			msg = constants.ACTIVE_FIELD_IS_REQUIRED
	// 		}
	// 	case constants.ZIP_CODE:
	// 		if value == "" {
	// 			msg = constants.ZIP_CODE_IS_FIELD_REQUIRED
	// 		}
	// 	}
	// }
	if dataMap[constants.NAME] == "" {
		msg = constants.NAME_FIELD_IS_REQUIRED
	}
	if dataMap[constants.ACTIVE] == nil {
		msg = constants.ACTIVE_FIELD_IS_REQUIRED
	}
	if dataMap[constants.CITY_ID] == nil {
		msg = constants.CITY_ID_FIELD_IS_REQUIRED
	}
	if dataMap[constants.ZIP_CODE] == "" {
		msg = constants.ZIP_CODE_IS_FIELD_REQUIRED
	}

	return msg
}
func validateRequired(stateDto dto.StatesDTO) string {

	var msg string = constants.EMPTY
	if stateDto.Name == constants.EMPTY {
		msg = constants.NAME_IS_REQUIRED
	}
	if stateDto.CityId == 0 {
		msg = constants.CITY_ID_IS_REQUIRED
	}
	if stateDto.ZipCode == constants.EMPTY {
		msg = constants.ZIP_CODE_IS_REQUIRED
	}
	return msg
}
