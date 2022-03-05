package utils

import (
	"encoding/xml"
	"fmt"
)

const (
	dollarCharCode = "USD"
	euroCharCode   = "EUR"
)

type Valute struct {
	XMLName  xml.Name `xml:"Valute"`
	Id       string   `xml:"ID,attr"`
	NumCode  string   `xml:"NumCode"`
	CharCode string   `xml:"CharCode"`
	Nominal  string   `xml:"Nominal"`
	Name     string   `xml:"Name"`
	Value    string   `xml:"Value"`
}

type ValuteResponse struct {
	XMLName xml.Name `xml:"ValCurs"`
	Date    string   `xml:"Date,attr"`
	Name    string   `xml:"name,atr"`
	Valute  []Valute `xml:"Valute"`
}

// GetDollarValue returns value for the dollar from the received ValuteResponse
func GetDollarValue(resp *ValuteResponse) (string, error) {
	return getValue(resp, dollarCharCode)
}

// GetEuroValue returns value for the euro from the received ValuteResponse
func GetEuroValue(resp *ValuteResponse) (string, error) {
	return getValue(resp, euroCharCode)
}

// getValue returns value for the received charCode from the received ValuteResponse
func getValue(resp *ValuteResponse, charCode string) (string, error) {
	for _, val := range resp.Valute {
		if val.CharCode == charCode {
			return val.Value, nil
		}
	}
	return "", fmt.Errorf("there is no value with code %s", charCode)
}
