package models

type RequestGetServiceCost struct {
	Pickup       Pickup   `json:"pickup"`
	Delivery     Delivery `json:"delivery"`
	SelfPickup   bool     `json:"selfPickup" example:"false"`  // самовывоз отправителя
	SelfDelivery bool     `json:"selfDelivery" example:"true"` // самовывоз получателя
	Parcel       Parcel   `json:"parcel"`
}

// @Description	Откуда
type Pickup struct {
	CityId      string `json:"cityId" example:"49694102"`
	Index       string `json:"index" example:"140012"`
	CityName    string `json:"cityName" example:"Москва"`
	RegionCode  string `json:"regionCode" example:"77"`
	CountryCode string `json:"countryCode" example:"RU"`
}

// @Description	Куда
type Delivery struct {
	CityId      string `json:"cityId" example:"49265227"`
	Index       string `json:"index" example:"140012"`
	CityName    string `json:"cityName" example:"Челябинск"`
	RegionCode  string `json:"regionCode" example:"74"`
	CountryCode string `json:"countryCode" example:"RU"`
}

// @Description	Посылка
type Parcel struct {
	Weight   int `json:"weight" example:"5"`
	Length   int `json:"length" example:"20"`
	Width    int `json:"width" example:"20"`
	Height   int `json:"height" example:"20"`
	Quantity int `json:"quantity" example:"1"`
}
