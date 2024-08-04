// utils/utils.go
package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type IPInfo struct {
	Query       string  `json:"query"`
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
}

func LocateIP(ip string) (IPInfo, error) {
	// Construct the URL for the API request
	url := fmt.Sprintf("http://ip-api.com/json/%s", ip)

	// Make the HTTP GET request
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error making HTTP request:", err)
		return IPInfo{}, err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return IPInfo{}, err
	}

	// Parse the JSON response
	var ipInfo IPInfo
	err = json.Unmarshal(body, &ipInfo)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return IPInfo{}, err
	}

	// Return the country from the parsed response
	return ipInfo, nil
}

func GenerateAnonymousID(ip, userAgent string) string {
    day := time.Now().Format("2006-01-02")
    data := fmt.Sprintf("%s|%s|%s", ip, userAgent, day)
    hash := sha256.Sum256([]byte(data))
    return hex.EncodeToString(hash[:])
}