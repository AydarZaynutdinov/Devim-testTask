package handler

import (
	"encoding/xml"
	"fmt"
	"golang.org/x/text/encoding/charmap"
	"io"
	"net/http"
	"strconv"
	"strings"
	"testTask/config"
	"testTask/valute"
	"time"
)

const (
	layout = "02/01/2006"
)

var url = "https://www.cbr.ru/scripts/XML_daily.asp?date_req=%s"

// GetValuteValueByCoordinate checks that url from the http.Request contains x and y coordinate and if this coordinate
//   is inside circle that is created by values form config file.
//  If coordinate is inside returns value of dollar from url by the current day.
//  If coordinate isn't inside returns value of euro from url by the current day.
func GetValuteValueByCoordinate(w http.ResponseWriter, r *http.Request) {
	coordinate, done := getParams(w, r)
	if done {
		return
	}

	updateUrl()

	valuteResponse := getValutes(w)
	if valuteResponse == nil {
		return
	}

	value, done := getValueByCoordinate(w, coordinate, valuteResponse)
	if done {
		return
	}

	_, _ = w.Write([]byte(value))
	w.WriteHeader(200)
}

// getValueByCoordinate returns value of valute by the coordinate
func getValueByCoordinate(w http.ResponseWriter, coordinate []int, valuteResponse *valute.ValuteResponse) (string, bool) {
	appConfig, err := config.NewConfig()
	if err != nil {
		panic(err)
	}
	var value string
	if appConfig.Check(coordinate[0], coordinate[1]) {
		value, err = valute.GetDollarValue(valuteResponse)
	} else {
		value, err = valute.GetEuroValue(valuteResponse)
	}
	if err != nil {
		if err != nil {
			_, _ = w.Write([]byte("Error during converting valute response"))
			w.WriteHeader(500)
			return "", true
		}
	}
	return value, false
}

// getValutes returns all valutes by the url
func getValutes(w http.ResponseWriter) *valute.ValuteResponse {
	resp, err := http.Get(url)
	if err != nil {
		_, _ = w.Write([]byte("Error during receiving valute"))
		w.WriteHeader(500)
		return nil
	}

	valuteResponse := &valute.ValuteResponse{}
	decoder := xml.NewDecoder(resp.Body)
	decoder.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
		switch charset {
		case "windows-1251":
			return charmap.Windows1251.NewDecoder().Reader(input), nil
		default:
			return nil, fmt.Errorf("unknown charset: %s", charset)
		}
	}
	if err := decoder.Decode(valuteResponse); err != nil {
		_, _ = w.Write([]byte("Error during converting request"))
		w.WriteHeader(500)
		return nil
	}
	return valuteResponse
}

// updateUrl updates url by current date
func updateUrl() {
	date := time.Now().Format(layout)
	url = fmt.Sprintf(url, date)
}

// getParams validates request.url and returns coordinate
func getParams(w http.ResponseWriter, r *http.Request) ([]int, bool) {
	requestParam := r.RequestURI
	if requestParam == "" || requestParam[1:] == "" {
		_, _ = w.Write([]byte("There are no required request params"))
		w.WriteHeader(404)
		return []int{0, 0}, true
	}
	requestParam = requestParam[2:]

	requestParams := strings.Split(requestParam, "&")
	if len(requestParams) < 2 {
		_, _ = w.Write([]byte("Incorrect count of required params"))
		w.WriteHeader(404)
		return []int{0, 0}, true
	}
	x, err := strconv.Atoi(requestParams[0][2:])
	if err != nil {
		_, _ = w.Write([]byte("First param contains incorrect value. Should be integer"))
		w.WriteHeader(404)
		return []int{0, 0}, true
	}
	y, err := strconv.Atoi(requestParams[1][2:])
	if err != nil {
		_, _ = w.Write([]byte("Second param contains incorrect value. Should be integer"))
		w.WriteHeader(404)
		return []int{0, 0}, true
	}
	return []int{x, y}, false
}
