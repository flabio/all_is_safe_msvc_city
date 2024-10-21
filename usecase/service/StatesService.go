package service

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	utils "github.com/flabio/safe_constants"
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
			utils.STATUS:  fiber.StatusBadRequest,
			utils.MESSAGE: utils.ERROR_QUERY,
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		utils.STATUS: fiber.StatusOK,
		utils.DATA:   result,
	})
}

func (s *statesService) GetStatesFindById(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params(utils.ID))
	result, err := s.states.GetStatesFindById(uint(id))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			utils.STATUS:  fiber.StatusBadRequest,
			utils.MESSAGE: utils.ERROR_QUERY,
		})
	}
	if result.Id == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			utils.STATUS: fiber.StatusNotFound,
		})
	}
	return c.Status(http.StatusOK).JSON(result)
}
func (s *statesService) GetStatesFindByIdOfCity(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params(utils.ID))

	result, err := s.states.GetStatesFindByIdOfCity(uint(id))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			utils.STATUS:  fiber.StatusBadRequest,
			utils.MESSAGE: utils.ERROR_QUERY,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		utils.STATUS: fiber.StatusOK,
		utils.DATA:   result,
	})
}
func (s *statesService) CreateState(c *fiber.Ctx) error {
	var states entities.States
	stateDto, msgError := validateState(0, s, c)
	if msgError != "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			utils.STATUS:  fiber.StatusBadRequest,
			utils.MESSAGE: msgError,
		})
	}
	deepcopier.Copy(stateDto).To(&states)
	result, err := s.states.CreateStates(states)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			utils.STATUS:  fiber.StatusInternalServerError,
			utils.MESSAGE: utils.ERROR_CREATE,
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		utils.STATUS:  fiber.StatusCreated,
		utils.DATA:    result,
		utils.MESSAGE: utils.CREATED,
	})
}

func (s *statesService) UpdateState(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params(utils.ID))
	stateDto, msgError := validateState(uint(id), s, c)
	if msgError != "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			utils.STATUS:  http.StatusBadRequest,
			utils.MESSAGE: msgError,
		})
	}
	state, _ := s.states.GetStatesFindById(uint(id))
	if state.Id == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			utils.STATUS:  http.StatusNotFound,
			utils.MESSAGE: utils.ID_NO_EXIST,
		})
	}
	deepcopier.Copy(stateDto).To(&state)
	result, err := s.states.UpdateStates(uint(id), state)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			utils.STATUS:  fiber.StatusInternalServerError,
			utils.MESSAGE: err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		utils.STATUS:  fiber.StatusOK,
		utils.DATA:    result,
		utils.MESSAGE: utils.UPDATED,
	})
}
func (s *statesService) DeleteState(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params(utils.ID))
	state, _ := s.states.GetStatesFindById(uint(id))
	if state.Id == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			utils.STATUS:  http.StatusNotFound,
			utils.MESSAGE: utils.ID_NO_EXIST,
		})
	}
	result, err := s.states.DeleteStates(uint(id))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			utils.STATUS:  fiber.StatusInternalServerError,
			utils.MESSAGE: utils.ERROR_DELETE,
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		utils.STATUS:  fiber.StatusOK,
		utils.DATA:    result,
		utils.MESSAGE: utils.REMOVED,
	})
}

func validateState(id uint, s *statesService, c *fiber.Ctx) (dto.StatesDTO, string) {
	defer func() {
		if r := recover(); r != nil {
			log.Println(utils.RECOVER_PANIC, r)
		}
	}()
	var msg string = ""
	stateDto := dto.StatesDTO{}
	body := c.Body()
	//dataMap := make(map[string]string)

	// fields := []string{
	// 	utils.NAME,
	// 	utils.CITY_ID,
	// 	utils.ACTIVE,
	// 	utils.ZIP_CODE,
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
		msg = utils.NAME_ALREADY_EXIST
	}
	return stateDto, msg
}

// func MapToStructStates(stateDto *dto.StatesDTO, dataMap map[string]string) {
// 	cityId, _ := strconv.Atoi(dataMap[utils.CITY_ID])
// 	isActive, _ := strconv.ParseBool(dataMap[utils.ACTIVE])

// 	state := dto.StatesDTO{
// 		Name:    dataMap[utils.NAME],
// 		CityId:  uint(cityId),
// 		ZipCode: dataMap[utils.ZIP_CODE],
// 		Active:  isActive,
// 	}
// 	*stateDto = state
// }

func validateField(dataMap map[string]interface{}) string {
	msg := ""
	// for field, value := range dataMap {
	// 	switch field {
	// 	case utils.NAME:
	// 		if value == "" {
	// 			msg = utils.NAME_FIELD_IS_REQUIRED
	// 		}
	// 	case utils.CITY_ID:
	// 		cityId, err := strconv.Atoi(value)
	// 		if err != nil || cityId == 0 {
	// 			msg = utils.CITY_ID_FIELD_IS_REQUIRED
	// 		}
	// 	case utils.ACTIVE:
	// 		isActive, err := strconv.ParseBool(value)
	// 		if err != nil || !isActive && !(!isActive) {
	// 			msg = utils.ACTIVE_FIELD_IS_REQUIRED
	// 		}
	// 	case utils.ZIP_CODE:
	// 		if value == "" {
	// 			msg = utils.ZIP_CODE_IS_FIELD_REQUIRED
	// 		}
	// 	}
	// }
	if dataMap[utils.NAME] == "" {
		msg = utils.NAME_FIELD_IS_REQUIRED
	}
	if dataMap[utils.ACTIVE] == nil {
		msg = utils.ACTIVE_FIELD_IS_REQUIRED
	}
	if dataMap[utils.CITY_ID] == nil {
		msg = utils.CITY_ID_FIELD_IS_REQUIRED
	}
	if dataMap[utils.ZIP_CODE] == "" {
		msg = utils.ZIP_CODE_IS_FIELD_REQUIRED
	}

	return msg
}
func validateRequired(stateDto dto.StatesDTO) string {

	var msg string = utils.EMPTY
	if stateDto.Name == utils.EMPTY {
		msg = utils.NAME_IS_REQUIRED
	}
	if stateDto.CityId == 0 {
		msg = utils.CITY_ID_IS_REQUIRED
	}
	if stateDto.ZipCode == utils.EMPTY {
		msg = utils.ZIP_CODE_IS_REQUIRED
	}
	return msg
}
