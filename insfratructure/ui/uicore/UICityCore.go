package uicore

import "github.com/safe_msvc_city/insfratructure/entities"

type UICityCore interface {
	GetCityFindAll() ([]entities.City, error)
	GetCityFindById(id uint) (entities.City, error)
	GetCityFindByName(id uint, name string) (bool, error)
	CreateCity(city entities.City) (entities.City, error)
	UpdateCity(id uint, city entities.City) (entities.City, error)
	DeleteCity(id uint) (bool, error)
}
