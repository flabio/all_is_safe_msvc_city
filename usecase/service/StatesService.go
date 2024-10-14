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
		utils.STATUS: fiber.StatusCreated,
		utils.DATA:   result,
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
		utils.STATUS: fiber.StatusOK,
		utils.DATA:   result,
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
		utils.STATUS: fiber.StatusOK,
		utils.DATA:   result,
		utils.MESSAGE:utils.REMOVED,
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
	var dataMap map[string]interface{}
	errJson := json.Unmarshal([]byte(body), &dataMap)
	if errJson != nil {
		msg = errJson.Error()
	}

	msgValid := validateField(dataMap)
	if msgValid != "" {
		return dto.StatesDTO{}, msgValid
	}

	MapToStructStates(&stateDto, dataMap)
	msgRequierd := validateRequired(stateDto)
	if msgRequierd != "" {
		return dto.StatesDTO{}, msgRequierd
	}
	existName, _ := s.states.GetStatesFindByName(id, stateDto.Name)
	if existName {
		msg = utils.NAME_ALREADY_EXIST
	}
	return stateDto, msg
}

func MapToStructStates(stateDto *dto.StatesDTO, dataMap map[string]interface{}) {
	state := dto.StatesDTO{
		Name:   dataMap[utils.NAME].(string),
		CityId: uint(dataMap[utils.CITY_ID].(float64)),
		Active: dataMap[utils.ACTIVE].(bool),
	}
	*stateDto = state

}
func validateField(value map[string]interface{}) string {
	msg := ""
	if value[utils.NAME] == nil {
		msg = utils.NAME_FIELD_IS_REQUIRED
	}
	if value[utils.ACTIVE] == nil {
		msg = utils.ACTIVE_FIELD_IS_REQUIRED
	}
	if value[utils.CITY_ID] == nil {
		msg = utils.CITY_ID_FIELD_IS_REQUIRED
	}

	return msg
}
func validateRequired(value dto.StatesDTO) string {
	msg := ""
	if len(value.Name) == 0 || value.Name == "" {
		msg = utils.NAME_IS_REQUIRED
	}
	if value.CityId == 0 {
		msg = utils.CITY_ID_IS_REQUIRED
	}
	return msg
}
