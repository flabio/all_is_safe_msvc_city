package uicore

import "github.com/safe_msvc_city/insfratructure/entities"

type UIStatesCore interface {
	GetStatesFindAll() ([]entities.States, error)
	GetStatesFindById(id uint) (entities.States, error)
	GetStatesFindByIdOfCity(id uint) ([]entities.States, error)
	GetStatesFindByName(id uint, name string) (bool, error)
	CreateStates(states entities.States) (entities.States, error)
	UpdateStates(id uint, states entities.States) (entities.States, error)
	DeleteStates(id uint) (bool, error)
}
