package dto



type CityDTO struct {
    Id        uint   `json:"id" `
    Name      string `json:"name" validate:"required|min_len:3"  message:"required:{field} is required|min_len: min 3 caracters" label:"Name"`
    Active    bool   `json:"active" validate:"required" message:"required:{field} is required" label:"Active"`
}