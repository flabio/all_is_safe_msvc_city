package core

import (
	"sync"

	utils "github.com/flabio/safe_constants"
	"github.com/safe_msvc_city/insfratructure/database"
	"github.com/safe_msvc_city/insfratructure/entities"
	"github.com/safe_msvc_city/insfratructure/ui/uicore"
)

func GetCityInstance() uicore.UICityCore {
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

func (db *OpenConnection) GetCityFindAll() ([]entities.City, error) {
	var cities []entities.City
	db.mux.Lock()
	result := db.connection.Order(utils.DB_ORDER_DESC).Find(&cities)
	defer db.mux.Unlock()
	defer database.CloseConnection()
	return cities, result.Error
}
func (db *OpenConnection) GetCityFindById(id uint) (entities.City, error) {
	var city entities.City
	db.mux.Lock()
	result := db.connection.Where(utils.DB_EQUAL_ID, id).Find(&city)
	defer db.mux.Unlock()
	defer database.CloseConnection()
	return city, result.Error
}
func (db *OpenConnection) CreateCity(city entities.City) (entities.City, error) {
	db.mux.Lock()
	err := db.connection.Create(&city).Error
	defer db.mux.Unlock()
	defer database.CloseConnection()
	return city, err
}
func (db *OpenConnection) UpdateCity(id uint, city entities.City) (entities.City, error) {
	db.mux.Lock()
	err := db.connection.Where(utils.DB_EQUAL_ID, id).Updates(&city).Error
	defer db.mux.Unlock()
	defer database.CloseConnection()
	return city, err
}
func (db *OpenConnection) DeleteCity(id uint) (bool, error) {
	db.mux.Lock()
	var city entities.City
	err := db.connection.Where(utils.DB_EQUAL_ID, id).Delete(&city).Error
	defer db.mux.Unlock()
	defer database.CloseConnection()
	return err == nil, err

}
func (db *OpenConnection) GetCityFindByName(id uint, name string) (bool, error) {
	db.mux.Lock()
	var city entities.City
	query := db.connection.Where(utils.DB_EQUAL_NAME, name)
	if id > 0 {
		query = query.Where(utils.DB_DIFF_ID, id)
	}
	query = query.Find(&city)
	defer db.mux.Unlock()
	defer database.CloseConnection()
	return query.RowsAffected > 0, query.Error

}
func (db *OpenConnection) GetCityFindByCode(id uint, name string) (bool, error) {
	db.mux.Lock()
	var city entities.City
	query := db.connection.Where(utils.DB_EQUAL_NAME, name)
	if id > 0 {
		query = query.Where(utils.DB_DIFF_ID, id)
	}
	query = query.Find(&city)

	defer db.mux.Unlock()
	defer database.CloseConnection()
	return city.Id > 0, query.Error
}
