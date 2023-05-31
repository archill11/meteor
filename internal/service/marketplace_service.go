package service

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"meteor/internal/models"
	"net/http"
	"strings"

	"github.com/valyala/fasthttp"
)

type DpdCfg struct {
	ClientNumber string `default:"" envconfig:"CLIENT_NUMBER"`
	ClientKey    string `default:"" envconfig:"CLIENT_KEY"`
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
					<index>%s</index>
					<cityName>%s</cityName>
					<regionCode>%s</regionCode>
					<countryCode>%s</countryCode>
				</pickup>
				<delivery>
					<cityId>%s</cityId>
					<index>%s</index>
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
		body.Pickup.CityId, body.Pickup.Index, body.Pickup.CityName, body.Pickup.RegionCode, body.Pickup.CountryCode,
		body.Delivery.CityId, body.Delivery.Index, body.Delivery.CityName, body.Delivery.RegionCode, body.Delivery.CountryCode,
		body.SelfPickup, body.SelfDelivery,
		body.Parcel.Weight, body.Parcel.Length, body.Parcel.Width, body.Parcel.Height, body.Parcel.Quantity,
	)

	resp, err := http.Post("http://ws.dpd.ru/services/calculator2?wsdl", "text/xml", strings.NewReader(xmlbody))
	if err != nil {
		return nil, fmt.Errorf("GetServiceCost http.Post: %v", err)
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("GetServiceCost io.ReadAll: %v", err)
	}
	var cAny Envelope
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

type Envelope struct {
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
