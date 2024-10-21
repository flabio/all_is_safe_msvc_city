package dto

type StatesDTO struct {
	Id      uint   `json:"id" `
	Name    string `json:"name" `
	ZipCode string `json:"zip_code"`
	CityId  uint   `json:"city_id"`
	Active  bool   `json:"active"`
}
