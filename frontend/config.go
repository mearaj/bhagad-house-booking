//go:build !js

package frontend

import (
	"os"
)

// LoadConfig looks for env API_URL. (if empty, then it defaults to DefaultApiURL)
func LoadConfig() (config Config) {
	url := os.Getenv(ApiURLKey)
	staticFolder := os.Getenv(StaticFolderKey)
	innerPort := os.Getenv(InnerPortKey)
	if url == "" {
		url = DefaultApiURL
	}
	if staticFolder == "" {
		staticFolder = DefaultStaticFolder
	}
	if innerPort == "" {
		innerPort = "8080"
	}
	config.ApiURL = url
	config.StaticFolder = staticFolder
	config.InnerPort = innerPort
	return config
}
