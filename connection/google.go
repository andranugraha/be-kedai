package connection

import (
	"kedai/backend/be-kedai/config"
	"strconv"

	"googlemaps.github.io/maps"
)

var (
	googleMaps *maps.Client
)

func ConnectGoogleMaps() (err error) {
	rateLimit, _ := strconv.Atoi(config.GetEnv("GOOGLE_RATE_LIMIT", "10"))
	googleMaps, err = maps.NewClient(maps.WithAPIKey(config.GetEnv("GOOGLE_API_KEY", "")), maps.WithRateLimit(rateLimit))

	return
}

func GetGoogleMaps() *maps.Client {
	return googleMaps
}
