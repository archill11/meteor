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
