// utils/utils.go
package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

type IPInfo struct {
    Country     string
    CountryCode string
}

func LocateIP(ip string) (IPInfo, error) {
    // Implement IP geolocation here
    // For now, we'll return dummy data
    return IPInfo{
        Country:     "Unknown",
        CountryCode: "UN",
    }, nil
}

func GenerateAnonymousID(ip, userAgent string) string {
    day := time.Now().Format("2006-01-02")
    data := fmt.Sprintf("%s|%s|%s", ip, userAgent, day)
    hash := sha256.Sum256([]byte(data))
    return hex.EncodeToString(hash[:])
}