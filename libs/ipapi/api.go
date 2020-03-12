package ipapi

import (
	"encoding/json"
	"fmt"
	"github.com/ihaiker/gokit/errors"
	"net/http"
	"time"
)

/*
{
  "status": "success",
  "country": "China",
  "countryCode": "CN",
  "region": "BJ",
  "regionName": "Beijing",
  "city": "Beijing",
  "zip": "",
  "lat": 39.9949,
  "lon": 116.316,
  "timezone": "Asia/Shanghai",
  "isp": "IDC, China Telecommunications Corporation",
  "org": "",
  "as": "AS23724 IDC, China Telecommunications Corporation",
  "query": "220.181.38.148"
}
*/

type IPInfo struct {
	Status      string  `json:"status"`
	Country     string  `json:"country"`
	CountryCode string  `json:"countryCode"`
	Region      string  `json:"region"`
	RegionName  string  `json:"regionName"`
	City        string  `json:"city"`
	Zip         string  `json:"zip"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
	Timezone    string  `json:"timezone"`
	Isp         string  `json:"isp"`
	Org         string  `json:"org"`
	As          string  `json:"as"`
	Query       string  `json:"query"`
}

func (ip IPInfo) String() string {
	if ip.Status != "success" {
		return ""
	}
	return fmt.Sprintf("%s %s %s", ip.Country, ip.RegionName, ip.City)
}

func Get(ip string) IPInfo {
	info := IPInfo{
		Status: "fail",
	}
	defer errors.Catch(func(re error) {
		info.Country = re.Error()
	})
	request, err := http.NewRequest(http.MethodGet, "http://ip-api.com/json/"+ip, nil)
	errors.Assert(err, "make request")

	client := http.Client{Timeout: time.Second * 3}
	resp, err := client.Do(request)
	errors.Assert(err, "request")

	defer func() { _ = resp.Body.Close() }()
	errors.Assert(json.NewDecoder(resp.Body).Decode(&info), "request decoder")
	return info
}
