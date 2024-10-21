package service

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"

	"github.com/ulule/deepcopier"

	utils "github.com/flabio/safe_constants"
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
			utils.STATUS:  fiber.StatusBadRequest,
			utils.MESSAGE: utils.ERROR_QUERY,
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		utils.STATUS: http.StatusOK,
		utils.DATA:   result,
	})
}
func (s *cityService) GetCityFindById(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params(utils.ID))
	result, err := s.cityRepository.GetCityFindById(uint(id))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			utils.STATUS:  fiber.StatusBadRequest,
			utils.MESSAGE: utils.ERROR_QUERY,
		})
	}
	if result.Id == 0 {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			utils.STATUS:  http.StatusNotFound,
			utils.MESSAGE: utils.ID_NO_EXIST,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		utils.STATUS: http.StatusOK,
		utils.DATA:   result,
	})
}
func (s *cityService) CreateCity(c *fiber.Ctx) error {
	var cityCreate entities.City

	cityDto, msgError := validateCity(0, s, c)
	if msgError != "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			utils.STATUS:  http.StatusBadRequest,
			utils.MESSAGE: msgError,
		})
	}
	deepcopier.Copy(cityDto).To(&cityCreate)
	result, err := s.cityRepository.CreateCity(cityCreate)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			utils.STATUS:  fiber.StatusInternalServerError,
			utils.MESSAGE: utils.ERROR_CREATE,
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		utils.STATUS: http.StatusCreated,
		utils.DATA:   result,
		utils.MESSAGE: utils.CREATED,
	})
}

func (s *cityService) UpdateCity(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params(utils.ID))
	cityDto, msgError := validateCity(uint(id), s, c)
	if msgError != "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			utils.STATUS:  http.StatusBadRequest,
			utils.MESSAGE: msgError,
		})
	}
	city, err := s.cityRepository.GetCityFindById(uint(id))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			utils.STATUS:  fiber.StatusInternalServerError,
			utils.MESSAGE: utils.ERROR_QUERY,
		})
	}
	if city.Id == 0 {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			utils.STATUS:  http.StatusNotFound,
			utils.MESSAGE: utils.ID_NO_EXIST,
		})
	}
	deepcopier.Copy(cityDto).To(&city)
	result, err := s.cityRepository.UpdateCity(uint(id), city)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			utils.STATUS:  fiber.StatusInternalServerError,
			utils.MESSAGE: utils.ERROR_UPDATE,
		})
	}
	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		utils.STATUS: http.StatusAccepted,
		utils.DATA:   result,
		utils.MESSAGE: utils.UPDATED,
	})
}
func (s *cityService) DeleteCity(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params(utils.ID))
	city, err := s.cityRepository.GetCityFindById(uint(id))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			utils.STATUS:  fiber.StatusInternalServerError,
			utils.MESSAGE: utils.ERROR_QUERY,
		})
	}
	if city.Id == 0 {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			utils.STATUS:  http.StatusNotFound,
			utils.MESSAGE: utils.ID_NO_EXIST,
		})
	}
	result, err := s.cityRepository.DeleteCity(uint(id))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			utils.STATUS:  fiber.StatusInternalServerError,
			utils.MESSAGE: utils.ERROR_DELETE,
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		utils.STATUS:  http.StatusOK,
		utils.MESSAGE: utils.REMOVED,
		utils.DATA:    result,
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
	if msgValid != utils.EMPTY {
		return dto.CityDTO{}, msgValid
	}

	helpers.MapToStruct(&cityDto, dataMap)
	msgRequired := helpers.ValidateRequiredCity(cityDto)
	if msgRequired != utils.EMPTY {
		return dto.CityDTO{}, msgRequired
	}
	existName, _ := s.cityRepository.GetCityFindByName(id, cityDto.Name)
	if existName {
		msg = utils.NAME_ALREADY_EXIST
	}
	return cityDto, msg
}

