package dto

type StatesDTO struct {
	Id     uint   `json:"id" `
	Name   string `json:"name" `
	CityId uint   `json:"city_id"`
	Active bool   `json:"active"`
}
