package service

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"meteor/internal/models"
	"net/http"
	"strings"
	"time"

	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

type DpdCfg struct {
	ClientNumber string `default:"" envconfig:"CLIENT_NUMBER"`
	ClientKey    string `default:"" envconfig:"CLIENT_KEY"`
}

type EnvelopeGetServiceCost struct {
	XMLName xml.Name `xml:"Envelope"`
	Text    string   `xml:",chardata"`
	S       string   `xml:"S,attr"`
	Body    struct {
		Text                             string `xml:",chardata"`
		GetServiceCostByParcels2Response struct {
			Text   string `xml:",chardata"`
			Ns2    string `xml:"ns2,attr"`
			Return []struct {
				Text        string `xml:",chardata"`
				ServiceCode string `xml:"serviceCode" json:"serviceCode"`
				ServiceName string `xml:"serviceName" json:"serviceName"`
				Cost        string `xml:"cost" json:"cost"`
				Days        string `xml:"days" json:"days"`
			} `xml:"return" json:"return"`
		} `xml:"getServiceCostByParcels2Response" json:"getServiceCostByParcels2Response"`
	} `xml:"Body" json:"Body"`
}

func (s *Service) GetServiceCost(ctx *fasthttp.RequestCtx, body models.RequestGetServiceCost) ([]byte, error) {
	xmlbody := fmt.Sprintf(`
		<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/" xmlns:ns1="http://dpd.ru/ws/calculator/2012-03-20">
		<SOAP-ENV:Body>
		<ns1:getServiceCostByParcels2>
			<request>
				<auth>
					<clientNumber>%s</clientNumber>
					<clientKey>%s</clientKey>
				</auth>
				<pickup>
					<cityId>%s</cityId>
					<cityName>%s</cityName>
					<regionCode>%s</regionCode>
					<countryCode>%s</countryCode>
				</pickup>
				<delivery>
					<cityId>%s</cityId>
					<cityName>%s</cityName>
					<regionCode>%s</regionCode>
					<countryCode>%s</countryCode>
				</delivery>
				<selfPickup>%v</selfPickup>
				<selfDelivery>%v</selfDelivery>
				<parcel>
					<weight>%d</weight>
					<length>%d</length>
					<width>%d</width>
					<height>%d</height>
					<quantity>%d</quantity>
				</parcel>
			</request>
		</ns1:getServiceCostByParcels2>
		</SOAP-ENV:Body>
		</SOAP-ENV:Envelope>`,

		s.cfg.Dpd.ClientNumber, s.cfg.Dpd.ClientKey,
		body.Pickup.CityId, body.Pickup.CityName, body.Pickup.RegionCode, body.Pickup.CountryCode,
		body.Delivery.CityId, body.Delivery.CityName, body.Delivery.RegionCode, body.Delivery.CountryCode,
		body.SelfPickup, body.SelfDelivery,
		body.Parcel.Weight, body.Parcel.Length, body.Parcel.Width, body.Parcel.Height, body.Parcel.Quantity,
	)

	s.Logger.Debug("GetServiceCost", zap.Any("RequestGetServiceCost", body))

	url := "https://ws.dpd.ru/services/calculator2?wsdl"
	method := "POST"
	client := &http.Client{
		Timeout: time.Second * 60,
	}
	req, err := http.NewRequest(method, url, strings.NewReader(xmlbody))
	if err != nil {
		return nil, fmt.Errorf("GetServiceCost http.NewRequest: %v", err)
	}

	// resp, err := http.Post("", "text/xml", )
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("GetServiceCost http.Post: %v", err)
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("GetServiceCost io.ReadAll: %v", err)
	}
	var cAny EnvelopeGetServiceCost
	if err := xml.Unmarshal(data, &cAny); err != nil {
		return nil, fmt.Errorf("GetServiceCost Error unmarshalling to XML: %v", err)
	}

	var res models.ResponseGetServiceCost
	for _, v := range cAny.Body.GetServiceCostByParcels2Response.Return {
		service := models.Service{
			ServiceCode: v.ServiceCode,
			ServiceName: v.ServiceName,
			Cost:        v.Cost,
			Days:        v.Days,
		}
		res.Services = append(res.Services, service)
	}

	result, err := json.Marshal(res)
	if nil != err {
		return nil, fmt.Errorf("GetServiceCost Error marshalling to JSON: %v", err)
	}
	return result, nil
}

type EnvelopeGetCitiesCashPay struct {
	XMLName xml.Name `xml:"Envelope"`
	Text    string   `xml:",chardata"`
	S       string   `xml:"S,attr"`
	Body    struct {
		Text                     string `xml:",chardata"`
		GetCitiesCashPayResponse struct {
			Text   string `xml:",chardata"`
			Ns2    string `xml:"ns2,attr"`
			Return []struct {
				Text         string `xml:",chardata"`
				CityId       string `xml:"cityId" json:"cityId"`
				CountryCode  string `xml:"countryCode" json:"countryCode"`
				CountryName  string `xml:"countryName" json:"countryName"`
				RegionCode   string `xml:"regionCode" json:"regionCode"`
				RegionName   string `xml:"regionName" json:"regionName"`
				CityCode     string `xml:"cityCode" json:"cityCode"`
				CityName     string `xml:"cityName" json:"cityName"`
				Abbreviation string `xml:"abbreviation" json:"abbreviation"`
				IndexMin     string `xml:"indexMin" json:"indexMin"`
				IndexMax     string `xml:"indexMax" json:"indexMax"`
			} `xml:"return"`
		} `xml:"getCitiesCashPayResponse"`
	} `xml:"Body"`
}

func (s *Service) GetCitiesCashPay() ([]byte, error) {
	xmlbody := fmt.Sprintf(`
		<SOAP-ENV:Envelope
			xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/" xmlns:ns1="http://dpd.ru/ws/geography/2015-05-20">
			<SOAP-ENV:Body>
				<ns1:getCitiesCashPay>
					<request>
						<auth>
						<clientNumber>` + s.cfg.Dpd.ClientNumber + `</clientNumber>
						<clientKey>` + s.cfg.Dpd.ClientKey + `</clientKey>
						</auth>
					</request>
				</ns1:getCitiesCashPay>
			</SOAP-ENV:Body>
		</SOAP-ENV:Envelope>`,
	)

	url := "https://ws.dpd.ru/services/geography2?wsdl"
	method := "POST"
	client := &http.Client{
		Timeout: time.Second * 60,
	}
	req, err := http.NewRequest(method, url, strings.NewReader(xmlbody))
	if err != nil {
		return nil, fmt.Errorf("GetServiceCost http.NewRequest: %v", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("GetServiceCost http.Post: %v", err)
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("GetServiceCost io.ReadAll: %v", err)
	}
	var cAny EnvelopeGetCitiesCashPay
	if err := xml.Unmarshal(data, &cAny); err != nil {
		return nil, fmt.Errorf("GetServiceCost Error unmarshalling to XML: %v", err)
	}

	result, err := json.Marshal(cAny.Body.GetCitiesCashPayResponse.Return)
	if nil != err {
		return nil, fmt.Errorf("GetServiceCost Error marshalling to JSON: %v", err)
	}

	return result, nil
}

func (s *Service) GetCityByName(sity string) (models.City, error) {
	var si models.City

	cities, err := s.GetCitiesCashPay()
	if nil != err {
		return si, fmt.Errorf("GetCityByName: %v", err)
	}

	var cAny []models.City
	if err := json.Unmarshal(cities, &cAny); err != nil {
		return si, fmt.Errorf("GetCityByName Error unmarshalling to JSON: %v", err)
	}

	for _, v := range cAny {
		if strings.EqualFold(v.CityName, sity) {
			si.CityId = v.CityId
			si.CountryCode = v.CountryCode
			si.RegionCode = v.RegionCode
			return si, nil
		}
	}

	return si, fmt.Errorf("sity not found: %s", sity)
}
