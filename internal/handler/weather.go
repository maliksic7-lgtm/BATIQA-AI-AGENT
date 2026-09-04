package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// weatherCache caches the last successful Pekanbaru weather lookup so we don't
// hammer the free API on every chat message.
var (
	weatherCache     string
	weatherCacheTime time.Time
)

// currentWeather returns a short human-readable weather line for Pekanbaru via
// Open-Meteo (free, no API key). Returns "" on any failure so callers degrade
// gracefully. Cached for 30 minutes.
func currentWeather(ctx context.Context) string {
	if weatherCache != "" && time.Since(weatherCacheTime) < 30*time.Minute {
		return weatherCache
	}
	// Open-Meteo supports lat/lon coordinates (no key needed).
	url := "https://api.open-meteo.com/v1/forecast?latitude=-0.5071&longitude=101.4478&current=temperature_2m,relative_humidity_2m,weather_code,wind_speed_10m"

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return ""
	}
	cli := &http.Client{Timeout: 6 * time.Second}
	resp, err := cli.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return ""
	}

	var out struct {
		Current struct {
			Temperature2M   float64 `json:"temperature_2m"`
			RelativeHumidity float64 `json:"relative_humidity_2m"`
			WeatherCode      int     `json:"weather_code"`
			WindSpeed10M    float64 `json:"wind_speed_10m"`
		} `json:"current"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return ""
	}

	desc := wmoWeatherDesc(out.Current.WeatherCode)
	line := fmt.Sprintf("%s, temp ±%.1f°C, humidity %d%%, angin ±%.0f km/j",
		desc, out.Current.Temperature2M, int(out.Current.RelativeHumidity), out.Current.WindSpeed10M)

	weatherCache = line
	weatherCacheTime = time.Now()
	return line
}

// wmoWeatherDesc maps a WMO weather code to a short Indonesian/English label.
func wmoWeatherDesc(code int) string {
	switch code {
	case 0:
		return "Cerah / cerah berawan"
	case 1, 2:
		return "Berawan sebagian"
	case 3:
		return "Berawan"
	case 45, 48:
		return "Berkabut"
	case 51, 53, 55, 56, 57:
		return "Gerbis hujan ringan"
	case 61, 63, 65, 66, 67, 80, 81, 82:
		return "Hujan"
	case 71, 73, 75, 77, 85, 86:
		return "Hujan es / salju ringan"
	case 95, 96, 99:
		return "Badai petir"
	default:
		return "Berawan"
	}
}
