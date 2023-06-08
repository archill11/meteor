package models

type ResponseGetServiceCost struct {
	Services []Service `json:"services"`
}

type Service struct {
	ServiceCode string `json:"serviceCode"`
	ServiceName string `json:"serviceName"`
	Cost        string `json:"cost"`
	Days        string `json:"days"`
}

type ResponseGetCitiesCashPay struct {
	Cities []City `json:"cities"`
}

type City struct {
	CityId       string `json:"cityId" example:"49694102"`
	CountryCode  string `json:"countryCode" example:"RU"`
	CountryName  string `json:"countryName,omitempty"`
	RegionCode   string `json:"regionCode" example:"77"`
	RegionName   string `json:"regionName,omitempty"`
	CityCode     string `json:"cityCode,omitempty"`
	CityName     string `json:"cityName" example:"Москва"`
	Abbreviation string `json:"abbreviation,omitempty"`
	IndexMin     string `json:"indexMin,omitempty" example:"140012"`
	IndexMax     string `json:"indexMax,omitempty" example:"140012"`
}
