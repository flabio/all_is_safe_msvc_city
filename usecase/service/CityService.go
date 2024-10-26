package service

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"

	"github.com/ulule/deepcopier"

	constants "github.com/flabio/safe_constants"
	"github.com/safe_msvc_city/core"
	"github.com/safe_msvc_city/insfratructure/entities"
	"github.com/safe_msvc_city/insfratructure/helpers"

	"github.com/safe_msvc_city/insfratructure/ui/global"
	"github.com/safe_msvc_city/insfratructure/ui/uicore"
	"github.com/safe_msvc_city/usecase/dto"
)

type cityService struct {
	cityRepository uicore.UICityCore
}

func NewCityService() global.UICity {
	return &cityService{
		cityRepository: core.GetCityInstance(),
	}
}

func (s *cityService) GetCityFindAll(c *fiber.Ctx) error {
	result, err := s.cityRepository.GetCityFindAll()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			constants.STATUS:  fiber.StatusBadRequest,
			constants.MESSAGE: constants.ERROR_QUERY,
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		constants.STATUS: http.StatusOK,
		constants.DATA:   result,
	})
}
func (s *cityService) GetCityFindById(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params(constants.ID))
	result, err := s.cityRepository.GetCityFindById(uint(id))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			constants.STATUS:  fiber.StatusBadRequest,
			constants.MESSAGE: constants.ERROR_QUERY,
		})
	}
	if result.Id == 0 {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			constants.STATUS:  http.StatusNotFound,
			constants.MESSAGE: constants.ID_NO_EXIST,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		constants.STATUS: http.StatusOK,
		constants.DATA:   result,
	})
}
func (s *cityService) CreateCity(c *fiber.Ctx) error {
	var cityCreate entities.City

	cityDto, msgError := validateCity(0, s, c)
	if msgError != "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			constants.STATUS:  http.StatusBadRequest,
			constants.MESSAGE: msgError,
		})
	}
	deepcopier.Copy(cityDto).To(&cityCreate)
	result, err := s.cityRepository.CreateCity(cityCreate)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			constants.STATUS:  fiber.StatusInternalServerError,
			constants.MESSAGE: constants.ERROR_CREATE,
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		constants.STATUS:  http.StatusCreated,
		constants.DATA:    result,
		constants.MESSAGE: constants.CREATED,
	})
}

func (s *cityService) UpdateCity(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params(constants.ID))
	cityDto, msgError := validateCity(uint(id), s, c)
	if msgError != "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			constants.STATUS:  http.StatusBadRequest,
			constants.MESSAGE: msgError,
		})
	}
	city, err := s.cityRepository.GetCityFindById(uint(id))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			constants.STATUS:  fiber.StatusInternalServerError,
			constants.MESSAGE: constants.ERROR_QUERY,
		})
	}
	if city.Id == 0 {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			constants.STATUS:  http.StatusNotFound,
			constants.MESSAGE: constants.ID_NO_EXIST,
		})
	}
	deepcopier.Copy(cityDto).To(&city)
	result, err := s.cityRepository.UpdateCity(uint(id), city)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			constants.STATUS:  fiber.StatusInternalServerError,
			constants.MESSAGE: constants.ERROR_UPDATE,
		})
	}
	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		constants.STATUS:  http.StatusAccepted,
		constants.DATA:    result,
		constants.MESSAGE: constants.UPDATED,
	})
}
func (s *cityService) DeleteCity(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params(constants.ID))
	city, err := s.cityRepository.GetCityFindById(uint(id))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			constants.STATUS:  fiber.StatusInternalServerError,
			constants.MESSAGE: constants.ERROR_QUERY,
		})
	}
	if city.Id == 0 {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			constants.STATUS:  http.StatusNotFound,
			constants.MESSAGE: constants.ID_NO_EXIST,
		})
	}
	result, err := s.cityRepository.DeleteCity(uint(id))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			constants.STATUS:  fiber.StatusInternalServerError,
			constants.MESSAGE: constants.ERROR_DELETE,
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		constants.STATUS:  http.StatusOK,
		constants.MESSAGE: constants.REMOVED,
		constants.DATA:    result,
	})
}

func validateCity(id uint, s *cityService, c *fiber.Ctx) (dto.CityDTO, string) {
	var cityDto dto.CityDTO
	var msg string = ""
	b := c.Body()

	var dataMap map[string]interface{}
	errJson := json.Unmarshal([]byte(b), &dataMap)

	if errJson != nil {
		msg = errJson.Error()
	}
	msgValid := helpers.ValidateFieldCity(dataMap)
	if msgValid != constants.EMPTY {
		return dto.CityDTO{}, msgValid
	}

	helpers.MapToStruct(&cityDto, dataMap)
	msgRequired := helpers.ValidateRequiredCity(cityDto)
	if msgRequired != constants.EMPTY {
		return dto.CityDTO{}, msgRequired
	}
	existName, _ := s.cityRepository.GetCityFindByName(id, cityDto.Name)
	if existName {
		msg = constants.NAME_ALREADY_EXIST
	}
	return cityDto, msg
}
