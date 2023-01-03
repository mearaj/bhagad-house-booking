//go:build !js

package frontend

import (
	"os"
)

// LoadConfig looks for env API_URL. (if empty, then it defaults to DefaultApiURL)
func LoadConfig() (config Config) {
	url := os.Getenv(ApiURLKey)
	if url == "" {
		url = DefaultApiURL
	}
	config.ApiURL = url
	return config
}
