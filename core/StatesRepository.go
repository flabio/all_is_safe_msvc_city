package core

import (
	"sync"

	var_db "github.com/flabio/safe_var_db"
	"github.com/safe_msvc_city/insfratructure/database"
	"github.com/safe_msvc_city/insfratructure/entities"
	"github.com/safe_msvc_city/insfratructure/ui/uicore"
)

func GetStatesInstance() uicore.UIStatesCore {
	var (
		_OPEN *OpenConnection
		_ONCE sync.Once
	)
	_ONCE.Do(func() {
		_OPEN = &OpenConnection{
			connection: database.DatabaseConnection(),
		}
	})
	return _OPEN
}

func (db *OpenConnection) GetStatesFindAll() ([]entities.States, error) {
	var states []entities.States
	db.mux.Lock()
	result := db.connection.Preload("City").Order(var_db.DB_ORDER_DESC).Find(&states)
	defer db.mux.Unlock()
	defer database.CloseConnection()
	return states, result.Error
}
func (db *OpenConnection) GetStatesFindById(id uint) (entities.States, error) {
	var state entities.States
	db.mux.Lock()
	result := db.connection.Where(var_db.DB_EQUAL_ID, id).Find(&state)
	defer db.mux.Unlock()
	defer database.CloseConnection()
	return state, result.Error
}
func (db *OpenConnection) GetStatesFindByIdOfCity(id uint) ([]entities.States, error) {
	var state []entities.States
	db.mux.Lock()
	result := db.connection.Preload("City").Where(var_db.DB_EQUAL_CITY_ID, id).Find(&state)
	defer db.mux.Unlock()
	defer database.CloseConnection()
	return state, result.Error
}

func (db *OpenConnection) CreateStates(state entities.States) (entities.States, error) {
	db.mux.Lock()
	err := db.connection.Create(&state).Error
	defer db.mux.Unlock()
	defer database.CloseConnection()
	return state, err
}
func (db *OpenConnection) UpdateStates(id uint, state entities.States) (entities.States, error) {
	db.mux.Lock()
	err := db.connection.Where(var_db.DB_EQUAL_ID, id).Updates(&state).Error
	defer db.mux.Unlock()
	defer database.CloseConnection()
	return state, err
}

func (db *OpenConnection) DeleteStates(id uint) (bool, error) {
	db.mux.Lock()
	var state entities.States
	query := db.connection.Where(var_db.DB_EQUAL_ID, id).Delete(&state)
	defer db.mux.Unlock()
	defer database.CloseConnection()
	return query.RowsAffected > 0, query.Error
}
func (db *OpenConnection) GetStatesFindByName(id uint, name string) (bool, error) {
	var state entities.States
	db.mux.Lock()
	query := db.connection.Where(var_db.DB_EQUAL_NAME, name)
	if id > 0 {
		query = query.Where(var_db.DB_DIFF_ID, id)
	}
	result := query.Find(&state)
	defer db.mux.Unlock()
	defer database.CloseConnection()
	return result.RowsAffected > 0, result.Error
}
