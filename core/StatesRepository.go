package core

import (
	"sync"

	var_db "github.com/flabio/safe_var_db"
	"github.com/safe_msvc_city/insfratructure/database"
	"github.com/safe_msvc_city/insfratructure/entities"
	"github.com/safe_msvc_city/insfratructure/ui/uicore"

	"gorm.io/gorm"
)

// OpenConnection representa una conexión abierta a la base de datos con un mutex para sincronización
type openConnection struct {
	connection *gorm.DB
	mux        sync.Mutex
}

// Variables globales para implementar el patrón singleton
var (
	_OPEN *openConnection
	_ONCE sync.Once
)

// GetStatesInstance devuelve una instancia única de OpenConnection que implementa UIStatesCore
func GetStatesInstance() uicore.UIStatesCore {
	_ONCE.Do(func() {
		_OPEN = &openConnection{
			connection: database.GetDatabaseInstance(),
		}
	})
	return _OPEN
}

// GetStatesFindAll obtiene todos los estados
func (db *openConnection) GetStatesFindAll() ([]entities.States, error) {
	var states []entities.States
	db.mux.Lock()
	defer db.mux.Unlock()

	result := db.connection.Preload("City").Order(var_db.DB_ORDER_DESC).Find(&states)
	return states, result.Error
}

// GetStatesFindById obtiene un estado por su ID
func (db *openConnection) GetStatesFindById(id uint) (entities.States, error) {
	var state entities.States
	db.mux.Lock()
	defer db.mux.Unlock()

	result := db.connection.Where(var_db.DB_EQUAL_ID, id).First(&state)
	return state, result.Error
}

// GetStatesFindByIdOfCity obtiene estados por el ID de la ciudad
func (db *openConnection) GetStatesFindByIdOfCity(id uint) ([]entities.States, error) {
	var states []entities.States
	db.mux.Lock()
	defer db.mux.Unlock()

	result := db.connection.Preload("City").Where(var_db.DB_EQUAL_CITY_ID, id).Find(&states)
	return states, result.Error
}

// CreateStates crea un nuevo estado
func (db *openConnection) CreateStates(state entities.States) (entities.States, error) {
	db.mux.Lock()
	defer db.mux.Unlock()

	err := db.connection.Create(&state).Error
	return state, err
}

// UpdateStates actualiza un estado existente
func (db *openConnection) UpdateStates(id uint, state entities.States) (entities.States, error) {
	db.mux.Lock()
	defer db.mux.Unlock()

	err := db.connection.Where(var_db.DB_EQUAL_ID, id).Updates(&state).Error
	return state, err
}

// DeleteStates elimina un estado por su ID
func (db *openConnection) DeleteStates(id uint) (bool, error) {
	db.mux.Lock()
	defer db.mux.Unlock()

	result := db.connection.Where(var_db.DB_EQUAL_ID, id).Delete(&entities.States{})
	return result.RowsAffected > 0, result.Error
}

// GetStatesFindByName verifica si existe un estado por nombre, excluyendo un ID específico si se proporciona
func (db *openConnection) GetStatesFindByName(id uint, name string) (bool, error) {
	var state entities.States
	db.mux.Lock()
	defer db.mux.Unlock()

	query := db.connection.Where(var_db.DB_EQUAL_NAME, name)
	if id > 0 {
		query = query.Where(var_db.DB_DIFF_ID, id)
	}
	result := query.First(&state)
	return result.RowsAffected > 0, result.Error
}
